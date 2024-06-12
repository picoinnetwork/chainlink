package client

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type totalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC RPCClient[CHAIN_ID, HEAD],
] []Node[CHAIN_ID, RPC]

func NewTotalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC RPCClient[CHAIN_ID, HEAD],
](nodes []Node[CHAIN_ID, RPC]) NodeSelector[CHAIN_ID, HEAD, RPC] {
	return totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC](nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC]) Select() Node[CHAIN_ID, RPC] {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *big.Int
	var nodes []Node[CHAIN_ID, RPC]
	var aliveNodes []Node[CHAIN_ID, RPC]

	for _, n := range s {
		state, chainInfo := n.StateAndLatest()
		if state != NodeStateAlive {
			continue
		}

		currentTD := chainInfo.BlockDifficulty

		aliveNodes = append(aliveNodes, n)
		if currentTD != nil && (highestTD == nil || currentTD.Cmp(highestTD) >= 0) {
			if highestTD == nil || currentTD.Cmp(highestTD) > 0 {
				highestTD = currentTD
				nodes = nil
			}
			nodes = append(nodes, n)
		}
	}

	//If all nodes have td == nil pick one from the nodes that are alive
	if len(nodes) == 0 {
		return firstOrHighestPriority(aliveNodes)
	}
	return firstOrHighestPriority(nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC]) Name() string {
	return NodeSelectionModeTotalDifficulty
}
