package actions_seth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var ContractDeploymentInterval = 200

// FundChainlinkNodesFromRootAddress sends native token amount (expressed in human-scale) to each Chainlink Node
// from root private key. It returns an error if any of the transactions failed.
func FundChainlinkNodesFromRootAddress(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	amount *big.Float,
) error {
	if len(client.PrivateKeys) == 0 {
		return errors.Wrap(errors.New(seth.ErrNoKeyLoaded), fmt.Sprintf("requested key: %d", 0))
	}

	return FundChainlinkNodes(logger, client, nodes, client.PrivateKeys[0], amount)
}

// FundChainlinkNodes sends native token amount (expressed in human-scale) to each Chainlink Node
// from private key's address. It returns an error if any of the transactions failed.
func FundChainlinkNodes(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
) error {
	keyAddressFn := func(cl contracts.ChainlinkNodeWithKeysAndAddress) (string, error) {
		return cl.PrimaryEthAddress()
	}
	return fundChainlinkNodesAtAnyKey(logger, client, nodes, privateKey, amount, keyAddressFn)
}

// FundChainlinkNodesAtKeyIndexFromRootAddress sends native token amount (expressed in human-scale) to each Chainlink Node
// from root private key.It returns an error if any of the transactions failed. It sends the funds to
// node address at keyIndex.
func FundChainlinkNodesAtKeyIndexFromRootAddress(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	amount *big.Float,
	keyIndex int,
) error {
	if len(client.PrivateKeys) == 0 {
		return errors.Wrap(errors.New(seth.ErrNoKeyLoaded), fmt.Sprintf("requested key: %d", 0))
	}

	return FundChainlinkNodesAtKeyIndex(logger, client, nodes, client.PrivateKeys[0], amount, keyIndex)
}

// FundChainlinkNodesAtKeyIndex sends native token amount (expressed in human-scale) to each Chainlink Node
// from private key's address. It returns an error if any of the transactions failed. It sends the funds to
// node address at keyIndex.
func FundChainlinkNodesAtKeyIndex(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
	keyIndex int,
) error {
	keyAddressFn := func(cl contracts.ChainlinkNodeWithKeysAndAddress) (string, error) {
		toAddress, err := cl.EthAddresses()
		if err != nil {
			return "", err
		}
		return toAddress[keyIndex], nil
	}
	return fundChainlinkNodesAtAnyKey(logger, client, nodes, privateKey, amount, keyAddressFn)
}

func fundChainlinkNodesAtAnyKey(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
	keyAddressFn func(contracts.ChainlinkNodeWithKeysAndAddress) (string, error),
) error {
	for _, cl := range nodes {
		toAddress, err := keyAddressFn(cl)
		// toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}

		fromAddress, err := privateKeyToAddress(privateKey)
		if err != nil {
			return err
		}

		receipt, err := SendFunds(logger, client, FundsToSendPayload{
			ToAddress:  common.HexToAddress(toAddress),
			Amount:     conversions.EtherToWei(amount),
			PrivateKey: privateKey,
		})
		if err != nil {
			logger.Err(err).
				Str("From", fromAddress.Hex()).
				Str("To", toAddress).
				Msg("Failed to fund Chainlink node")

			return err
		}

		txHash := "(none)"
		if receipt != nil {
			txHash = receipt.TxHash.String()
		}

		logger.Info().
			Str("From", fromAddress.Hex()).
			Str("To", toAddress).
			Str("TxHash", txHash).
			Str("Amount", amount.String()).
			Msg("Funded Chainlink node")
	}

	return nil
}

type FundsToSendPayload struct {
	ToAddress  common.Address
	Amount     *big.Int
	PrivateKey *ecdsa.PrivateKey
	GasLimit   *int64
	GasPrice   *big.Int
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	TxTimeout  *time.Duration
}

// TODO: move to CTF?
// SendFunds sends native token amount (expressed in human-scale) from address controlled by private key
// to given address. You can override any or none of the following: gas limit, gas price, gas fee cap, gas tip cap.
// Values that are not set will be estimated or taken from config.
func SendFunds(logger zerolog.Logger, client *seth.Client, payload FundsToSendPayload) (*types.Receipt, error) {
	fromAddress, err := privateKeyToAddress(payload.PrivateKey)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	nonce, err := client.Client.PendingNonceAt(ctx, fromAddress)
	defer cancel()
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(client.Cfg.Network.TransferGasFee)
	gasPrice := big.NewInt(0)
	gasFeeCap := big.NewInt(0)
	gasTipCap := big.NewInt(0)

	if payload.GasLimit != nil {
		gasLimit = uint64(*payload.GasLimit)
	}

	if client.Cfg.Network.EIP1559DynamicFees {
		// if any of the dynamic fees are not set, we need to either estimate them or read them from config
		if payload.GasFeeCap == nil || payload.GasTipCap == nil {
			// estimatior or config reading happens here
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasFeeCap = txOptions.GasFeeCap
			gasTipCap = txOptions.GasTipCap
		}

		// override with payload values if they are set
		if payload.GasFeeCap != nil {
			gasFeeCap = payload.GasFeeCap
		}

		if payload.GasTipCap != nil {
			gasTipCap = payload.GasTipCap
		}
	}

	if !client.Cfg.Network.EIP1559DynamicFees {
		if payload.GasPrice == nil {
			txOptions := client.NewTXOpts((seth.WithGasLimit(gasLimit)))
			gasPrice = txOptions.GasPrice
		} else {
			gasPrice = payload.GasPrice
		}
	}

	var rawTx types.TxData

	if client.Cfg.Network.EIP1559DynamicFees {
		rawTx = &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &payload.ToAddress,
			Value:     payload.Amount,
			Gas:       gasLimit,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
		}
	} else {
		rawTx = &types.LegacyTx{
			Nonce:    nonce,
			To:       &payload.ToAddress,
			Value:    payload.Amount,
			Gas:      gasLimit,
			GasPrice: gasPrice,
		}
	}

	signedTx, err := types.SignNewTx(payload.PrivateKey, types.LatestSignerForChainID(big.NewInt(client.ChainID)), rawTx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign tx")
	}

	txTimeout := client.Cfg.Network.TxnTimeout.Duration()
	if payload.TxTimeout != nil {
		txTimeout = *payload.TxTimeout
	}

	ctx, cancel = context.WithTimeout(ctx, txTimeout)
	defer cancel()
	err = client.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send transaction")
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("TxHash", signedTx.Hash().String()).
		Str("Amount", conversions.WeiToEther(payload.Amount).String()).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Str("Gas Price", gasPrice.String()).
		Str("Gas Fee Cap", gasFeeCap.String()).
		Str("Gas Tip Cap", gasTipCap.String()).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("Sent funds")

	return client.WaitMined(ctx, logger, client.Client, signedTx)
}

// DeployForwarderContracts first deploys Operator Factory and then uses it to deploy given number of
// operator and forwarder pairs. It waits for each transaction to be mined and then extracts operator and
// forwarder addresses from emitted events.
func DeployForwarderContracts(
	t *testing.T,
	seth *seth.Client,
	linkTokenAddress common.Address,
	numberOfOperatorForwarderPairs int,
) (operators []common.Address, authorizedForwarders []common.Address, operatorFactoryInstance contracts.OperatorFactory) {
	instance, err := contracts.DeployEthereumOperatorFactory(seth, linkTokenAddress)
	require.NoError(t, err, "failed to create new instance of operator factory")
	operatorFactoryInstance = &instance

	for i := 0; i < numberOfOperatorForwarderPairs; i++ {
		decodedTx, err := seth.Decode(operatorFactoryInstance.DeployNewOperatorAndForwarder())
		require.NoError(t, err, "Deploying new operator with proposed ownership with forwarder shouldn't fail")

		for i, event := range decodedTx.Events {
			require.True(t, len(event.Topics) > 0, fmt.Sprintf("Event %d should have topics", i))
			switch event.Topics[0] {
			case operator_factory.OperatorFactoryOperatorCreated{}.Topic().String():
				if address, ok := event.EventData["operator"]; ok {
					operators = append(operators, address.(common.Address))
				} else {
					require.Fail(t, "Operator address not found in event", event)
				}
			case operator_factory.OperatorFactoryAuthorizedForwarderCreated{}.Topic().String():
				if address, ok := event.EventData["forwarder"]; ok {
					authorizedForwarders = append(authorizedForwarders, address.(common.Address))
				} else {
					require.Fail(t, "Forwarder address not found in event", event)
				}
			}
		}
	}
	return operators, authorizedForwarders, operatorFactoryInstance
}

// WatchNewOCRRound watches for a new OCR round, similarly to StartNewRound, but it does not explicitly request a new
// round from the contract, as this can cause some odd behavior in some cases. It announces success if latest round
// is >= roundNumber.
func WatchNewOCRRound(
	l zerolog.Logger,
	seth *seth.Client,
	roundNumber int64,
	ocrInstances []contracts.OffChainAggregatorWithRounds,
	timeout time.Duration,
) error {
	confirmed := make(map[string]bool)
	timeoutC := time.After(timeout)
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	l.Info().Msgf("Waiting for round %d to be confirmed by all nodes", roundNumber)

	for {
		select {
		case <-timeoutC:
			return fmt.Errorf("timeout waiting for round %d to be confirmed. %d/%d nodes confirmed it", roundNumber, len(confirmed), len(ocrInstances))
		case <-ticker.C:
			for i := 0; i < len(ocrInstances); i++ {
				if confirmed[ocrInstances[i].Address()] {
					continue
				}
				ctx, cancel := context.WithTimeout(context.Background(), seth.Cfg.Network.TxnTimeout.Duration())
				roundData, err := ocrInstances[i].GetLatestRound(ctx)
				if err != nil {
					cancel()
					return fmt.Errorf("getting latest round from OCR instance %d have failed: %w", i+1, err)
				}
				cancel()
				if roundData.RoundId.Cmp(big.NewInt(roundNumber)) >= 0 {
					l.Debug().Msgf("OCR instance %d/%d confirmed round %d", i+1, len(ocrInstances), roundNumber)
					confirmed[ocrInstances[i].Address()] = true
				}
			}
			if len(confirmed) == len(ocrInstances) {
				return nil
			}
		}
	}
}

// AcceptAuthorizedReceiversOperator sets authorized receivers for each operator contract to
// authorizedForwarder and authorized EA to nodeAddresses. Once done, it confirms that authorizations
// were set correctly.
func AcceptAuthorizedReceiversOperator(
	t *testing.T,
	logger zerolog.Logger,
	seth *seth.Client,
	operator common.Address,
	authorizedForwarder common.Address,
	nodeAddresses []common.Address,
) {
	operatorInstance, err := contracts.LoadEthereumOperator(logger, seth, operator)
	require.NoError(t, err, "Loading operator contract shouldn't fail")
	forwarderInstance, err := contracts.LoadEthereumAuthorizedForwarder(seth, authorizedForwarder)
	require.NoError(t, err, "Loading authorized forwarder contract shouldn't fail")

	err = operatorInstance.AcceptAuthorizedReceivers([]common.Address{authorizedForwarder}, nodeAddresses)
	require.NoError(t, err, "Accepting authorized forwarder shouldn't fail")

	senders, err := forwarderInstance.GetAuthorizedSenders(testcontext.Get(t))
	require.NoError(t, err, "Getting authorized senders shouldn't fail")
	var nodesAddrs []string
	for _, o := range nodeAddresses {
		nodesAddrs = append(nodesAddrs, o.Hex())
	}
	require.Equal(t, nodesAddrs, senders, "Senders addresses should match node addresses")

	owner, err := forwarderInstance.Owner(testcontext.Get(t))
	require.NoError(t, err, "Getting authorized forwarder owner shouldn't fail")
	require.Equal(t, operator.Hex(), owner, "Forwarder owner should match operator")
}

// TrackForwarder creates forwarder track for a given Chainlink node
func TrackForwarder(
	t *testing.T,
	seth *seth.Client,
	authorizedForwarder common.Address,
	node contracts.ChainlinkNodeWithForwarder,
) {
	l := logging.GetTestLogger(t)
	chainID := big.NewInt(seth.ChainID)
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	require.NoError(t, err, "Forwarder track should be created")
	l.Info().Str("NodeURL", node.GetConfig().URL).
		Str("ForwarderAddress", authorizedForwarder.Hex()).
		Str("ChaindID", chainID.String()).
		Msg("Forwarder tracked")
}

// DeployOCRv2Contracts deploys a number of OCRv2 contracts and configures them with defaults
func DeployOCRv2Contracts(
	l zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenAddress common.Address,
	transmitters []string,
	ocrOptions contracts.OffchainOptions,
) ([]contracts.OffchainAggregatorV2, error) {
	var ocrInstances []contracts.OffchainAggregatorV2
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contracts.DeployOffchainAggregatorV2(
			l,
			seth,
			linkTokenAddress,
			ocrOptions,
		)
		if err != nil {
			return nil, fmt.Errorf("OCRv2 instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, &ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Gather address payees
	var payees []string
	for range transmitters {
		payees = append(payees, seth.Addresses[0].Hex())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err := ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}
	return ocrInstances, nil
}

// ConfigureOCRv2AggregatorContracts sets configuration for a number of OCRv2 contracts
func ConfigureOCRv2AggregatorContracts(
	contractConfig *contracts.OCRv2Config,
	ocrv2Contracts []contracts.OffchainAggregatorV2,
) error {
	for contractCount, ocrInstance := range ocrv2Contracts {
		// Exclude the first node, which will be used as a bootstrapper
		err := ocrInstance.SetConfig(contractConfig)
		if err != nil {
			return fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

// TeardownSuite tears down networks/clients and environment and creates a logs folder for failed tests in the
// specified path. Can also accept a testreporter (if one was used) to log further results
func TeardownSuite(
	t *testing.T,
	chainClient *seth.Client,
	env *environment.Environment,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	failingLogLevel zapcore.Level, // Examines logs after the test, and fails the test if any Chainlink logs are found at or above provided level
	grafnaUrlProvider testreporters.GrafanaURLProvider,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel, grafnaUrlProvider); err != nil {
		return fmt.Errorf("Error dumping environment logs, leaving environment running for manual retrieval, err: %w", err)
	}
	// Delete all jobs to stop depleting the funds
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if chainlinkNodes != nil && len(chainlinkNodes) > 0 {
		if err := ReturnFundsFromNodes(l, chainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
			// This printed line is required for tests that use real funds to propagate the failure
			// out to the system running the test. Do not remove
			fmt.Println(environment.FAILED_FUND_RETURN)
			l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
					"Environment is left running so you can try manually!")
		}
	} else {
		l.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}

	return env.Shutdown()
}

// TeardownRemoteSuite sends a report and returns funds from chainlink nodes to network's default wallet
func TeardownRemoteSuite(
	t *testing.T,
	client *seth.Client,
	namespace string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	grafnaUrlProvider testreporters.GrafanaURLProvider,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.SendReport(t, namespace, "./", optionalTestReporter, grafnaUrlProvider); err != nil {
		l.Warn().Err(err).Msg("Error writing test report")
	}
	// Delete all jobs to stop depleting the funds
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if err = ReturnFundsFromNodes(l, client, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
		l.Error().Err(err).Str("Namespace", namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}

	// This is a failsafe, we should never use ephemeral keys on live networks
	if !client.Cfg.IsSimulatedNetwork() {
		if err := ReturnFundFromEphemeralKeys(l, client); err != nil {
			l.Error().Err(err).Str("Namespace", namespace).
				Msg("Error attempting to return funds from ephemeral keys to root key")

			pkStrings := []string{}
			for _, pk := range client.PrivateKeys {
				// Convert the D field (private key) to a byte slice of length 32
				privateKeyBytes := pk.D.Bytes()
				// Ensure the byte slice is exactly 32 bytes long, as required for Ethereum
				privateKeyBytes32 := make([]byte, 32)
				copy(privateKeyBytes32[32-len(privateKeyBytes):], privateKeyBytes)

				// Convert to a hexadecimal string
				privateKeyHex := hex.EncodeToString(privateKeyBytes32)

				pkStrings = append(pkStrings, "0x"+privateKeyHex)
			}

			privateKeyJson, jsonErr := json.Marshal(pkStrings)
			if jsonErr != nil {
				l.Error().
					Err(jsonErr).
					Msg("Error marshalling private keys to JSON. Funds are left in ephemeral keys")

				return err
			}

			fileName := "ephemeral_addresses_private_keys.json"
			writeErr := os.WriteFile(fileName, privateKeyJson, 0600)
			if writeErr != nil {
				l.Error().
					Err(writeErr).
					Msg("Error writing ephemeral addresses private keys to file. Funds are left in ephemeral keys")

				return err
			}
			absolutePath, pathErr := filepath.Abs(fileName)
			if pathErr != nil {
				l.Error().
					Err(pathErr).
					Str("FileName", fileName).
					Msg("Error getting absolute path of file with private keys. Try looking for it yourself.")

				return err
			}

			l.Info().
				Str("Filepath", absolutePath).
				Msg("Private keys for ephemeral addresses are saved in the file. You can use them to return funds manually.")
		}
	}

	return err
}

// DeleteAllJobs deletes all jobs from all chainlink nodes
// added here temporarily to avoid circular import
func DeleteAllJobs(chainlinkNodes []*client.ChainlinkK8sClient) error {
	for _, node := range chainlinkNodes {
		if node == nil {
			return fmt.Errorf("found a nil chainlink node in the list of chainlink nodes while tearing down: %v", chainlinkNodes)
		}
		jobs, _, err := node.ReadJobs()
		if err != nil {
			return fmt.Errorf("error reading jobs from chainlink node, err: %w", err)
		}
		for _, maps := range jobs.Data {
			if _, ok := maps["id"]; !ok {
				return fmt.Errorf("error reading job id from chainlink node's jobs %+v", jobs.Data)
			}
			id := maps["id"].(string)
			_, err := node.DeleteJob(id)
			if err != nil {
				return fmt.Errorf("error deleting job from chainlink node, err: %w", err)
			}
		}
	}
	return nil
}

// StartNewRound requests a new round from the ocr contracts and returns once transaction was mined
func StartNewRound(
	ocrInstances []contracts.OffChainAggregatorWithRounds,
) error {
	for i := 0; i < len(ocrInstances); i++ {
		err := ocrInstances[i].RequestNewRound()
		if err != nil {
			return fmt.Errorf("requesting new OCR round %d have failed: %w", i+1, err)
		}
	}
	return nil
}

// DeployOCRContractsForwarderFlow deploys and funds a certain number of offchain
// aggregator contracts with forwarders as effectiveTransmitters
func DeployOCRContractsForwarderFlow(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
	forwarderAddresses []common.Address,
) ([]contracts.OffchainAggregator, error) {
	transmitterPayeesFn := func() (transmitters []string, payees []string, err error) {
		transmitters = make([]string, 0)
		payees = make([]string, 0)
		for _, forwarderCommonAddress := range forwarderAddresses {
			forwarderAddress := forwarderCommonAddress.Hex()
			transmitters = append(transmitters, forwarderAddress)
			payees = append(payees, seth.Addresses[0].Hex())
		}

		return
	}

	transmitterAddressesFn := func() ([]common.Address, error) {
		return forwarderAddresses, nil
	}

	return deployAnyOCRv1Contracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn, transmitterAddressesFn)
}

// DeployOCRv1Contracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRv1Contracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
) ([]contracts.OffchainAggregator, error) {
	transmitterPayeesFn := func() (transmitters []string, payees []string, err error) {
		transmitters = make([]string, 0)
		payees = make([]string, 0)
		for _, node := range workerNodes {
			var addr string
			addr, err = node.PrimaryEthAddress()
			if err != nil {
				err = fmt.Errorf("error getting node's primary ETH address: %w", err)
				return
			}
			transmitters = append(transmitters, addr)
			payees = append(payees, seth.Addresses[0].Hex())
		}

		return
	}

	transmitterAddressesFn := func() ([]common.Address, error) {
		transmitterAddresses := make([]common.Address, 0)
		for _, node := range workerNodes {
			primaryAddress, err := node.PrimaryEthAddress()
			if err != nil {
				return nil, err
			}
			transmitterAddresses = append(transmitterAddresses, common.HexToAddress(primaryAddress))
		}

		return transmitterAddresses, nil
	}

	return deployAnyOCRv1Contracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn, transmitterAddressesFn)
}

func deployAnyOCRv1Contracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
	getTransmitterAndPayeesFn func() ([]string, []string, error),
	getTransmitterAddressesFn func() ([]common.Address, error),
) ([]contracts.OffchainAggregator, error) {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contracts.DeployOffchainAggregator(logger, seth, linkTokenContractAddress, contracts.DefaultOffChainAggregatorOptions())
		if err != nil {
			return nil, fmt.Errorf("OCR instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, &ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	var err error
	transmitters, payees, err = getTransmitterAndPayeesFn()
	if err != nil {
		return nil, fmt.Errorf("error getting transmitter and payees: %w", err)
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err := ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Set Config
	transmitterAddresses, err := getTransmitterAddressesFn()
	if err != nil {
		return nil, fmt.Errorf("getting transmitter addresses should not fail: %w", err)
	}

	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			workerNodes,
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			transmitterAddresses,
		)
		if err != nil {
			return nil, fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	return ocrInstances, nil
}

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func WatchNewFluxRound(
	l zerolog.Logger,
	seth *seth.Client,
	roundNumber int64,
	fluxInstance contracts.FluxAggregator,
	timeout time.Duration,
) error {
	timeoutC := time.After(timeout)
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	l.Info().Msgf("Waiting for flux round %d to be confirmed by flux aggregator", roundNumber)

	for {
		select {
		case <-timeoutC:
			return fmt.Errorf("timeout waiting for round %d to be confirmed", roundNumber)
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), seth.Cfg.Network.TxnTimeout.Duration())
			roundId, err := fluxInstance.LatestRoundID(ctx)
			if err != nil {
				cancel()
				return fmt.Errorf("getting latest round from flux instance has failed: %w", err)
			}
			cancel()
			if roundId.Cmp(big.NewInt(roundNumber)) >= 0 {
				l.Debug().Msgf("Flux instance confirmed round %d", roundNumber)
				return nil
			}
		}
	}
}

func EstimateCostForChainlinkOperations(l zerolog.Logger, client *seth.Client, network blockchain.EVMNetwork, amountOfOperations int) (*big.Float, error) {
	bigAmountOfOperations := big.NewInt(int64(amountOfOperations))
	estimations := client.CalculateGasEstimations(client.NewDefaultGasEstimationRequest())

	// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
	// total gas limit = chainlink gas limit + gas limit buffer
	gasLimit := network.GasEstimationBuffer + network.ChainlinkTransactionLimit
	// gas cost for TX = total gas limit * estimated gas price

	var gasPriceInWei *big.Int
	if client.Cfg.Network.EIP1559DynamicFees {
		gasPriceInWei = estimations.GasFeeCap
	} else {
		gasPriceInWei = estimations.GasPrice
	}

	gasCostPerOperationWei := big.NewInt(1).Mul(big.NewInt(1).SetUint64(gasLimit), gasPriceInWei)
	gasCostPerOperationETH := conversions.WeiToEther(gasCostPerOperationWei)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalEthForAllOperations := conversions.WeiToEther(totalWeiForAllOperations)

	l.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
}
