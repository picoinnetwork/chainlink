package types

import "time"

type TransactionManagerChainConfig interface {
	BroadcasterChainConfig
	ConfirmerChainConfig
	ReaperChainConfig
}

type TransactionManagerFeeConfig interface {
	BroadcasterFeeConfig
	ResenderFeeConfig
}

type TransactionManagerTransactionsConfig interface {
	BroadcasterTransactionsConfig
	ConfirmerTransactionsConfig
	ResenderTransactionsConfig
	ReaperTransactionsConfig

	ForwardersEnabled() bool
	MaxQueued() uint64
}

type BroadcasterChainConfig interface {
	IsL2() bool
}

type BroadcasterFeeConfig interface {
	MaxFeePrice() string     // logging value
	FeePriceDefault() string // logging value
}

type BroadcasterTransactionsConfig interface {
	MaxInFlight() uint32
}

type BroadcasterListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type ConfirmerChainConfig interface {
	RPCDefaultBatchSize() uint32
	FinalityDepth() uint32
}

type ConfirmerTransactionsConfig interface {
	MaxInFlight() uint32
	ForwardersEnabled() bool
}

type ResenderChainConfig interface {
	RPCDefaultBatchSize() uint32
}

type ResenderDatabaseConfig interface {
	// from pg.QConfig
	DefaultQueryTimeout() time.Duration
}

type ResenderFeeConfig interface {
	BumpTxDepth() uint32
	LimitDefault() uint32

	// from gas.Config
	BumpThreshold() uint64
	MaxFeePrice() string // logging value
	BumpPercent() uint16
}

type ResenderTransactionsConfig interface {
	ResendAfterThreshold() time.Duration
	MaxInFlight() uint32
}

// ReaperConfig is the config subset used by the reaper
//
//go:generate mockery --quiet --name ReaperChainConfig --structname ReaperConfig --output ./mocks/ --case=underscore
type ReaperChainConfig interface {
	FinalityDepth() uint32
}

type ReaperTransactionsConfig interface {
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
}
