// Code generated by mockery v2.38.0. DO NOT EDIT.

package client

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink-common/pkg/assets"

	context "context"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/smartcontractkit/chainlink/v2/common/types"
)

// mockRPC is an autogenerated mock type for the RPC type
type mockRPC[CHAIN_ID types.ID, SEQ types.Sequence, ADDR types.Hashable, BLOCK_HASH types.Hashable, TX interface{}, TX_HASH types.Hashable, EVENT interface{}, EVENT_OPS interface{}, TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH], FEE feetypes.Fee, HEAD types.Head[BLOCK_HASH]] struct {
	mock.Mock
}

// BalanceAt provides a mock function with given fields: ctx, accountAddress, blockNumber
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, accountAddress, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for BalanceAt")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, accountAddress, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) *big.Int); ok {
		r0 = rf(ctx, accountAddress, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR, *big.Int) error); ok {
		r1 = rf(ctx, accountAddress, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchCallContext provides a mock function with given fields: ctx, b
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) BatchCallContext(ctx context.Context, b []interface{}) error {
	ret := _m.Called(ctx, b)

	if len(ret) == 0 {
		panic("no return value specified for BatchCallContext")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}) error); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BlockByHash provides a mock function with given fields: ctx, hash
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for BlockByHash")
	}

	var r0 HEAD
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, BLOCK_HASH) (HEAD, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, BLOCK_HASH) HEAD); ok {
		r0 = rf(ctx, hash)
	} else {
		r0 = ret.Get(0).(HEAD)
	}

	if rf, ok := ret.Get(1).(func(context.Context, BLOCK_HASH) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockByNumber provides a mock function with given fields: ctx, number
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for BlockByNumber")
	}

	var r0 HEAD
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (HEAD, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) HEAD); ok {
		r0 = rf(ctx, number)
	} else {
		r0 = ret.Get(0).(HEAD)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallContext provides a mock function with given fields: ctx, result, method, args
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, result, method)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CallContext")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, result, method, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CallContract provides a mock function with given fields: ctx, msg, blockNumber
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) CallContract(ctx context.Context, msg interface{}, blockNumber *big.Int) ([]byte, error) {
	ret := _m.Called(ctx, msg, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for CallContract")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *big.Int) ([]byte, error)); ok {
		return rf(ctx, msg, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *big.Int) []byte); ok {
		r0 = rf(ctx, msg, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, *big.Int) error); ok {
		r1 = rf(ctx, msg, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChainID provides a mock function with given fields: ctx
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) ChainID(ctx context.Context) (CHAIN_ID, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ChainID")
	}

	var r0 CHAIN_ID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (CHAIN_ID, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) CHAIN_ID); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(CHAIN_ID)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientVersion provides a mock function with given fields: _a0
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) ClientVersion(_a0 context.Context) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ClientVersion")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) Close() {
	_m.Called()
}

// CodeAt provides a mock function with given fields: ctx, account, blockNumber
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	ret := _m.Called(ctx, account, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for CodeAt")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) ([]byte, error)); ok {
		return rf(ctx, account, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) []byte); ok {
		r0 = rf(ctx, account, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR, *big.Int) error); ok {
		r1 = rf(ctx, account, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Dial provides a mock function with given fields: ctx
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) Dial(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Dial")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DialHTTP provides a mock function with given fields:
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) DialHTTP() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DialHTTP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DisconnectAll provides a mock function with given fields:
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) DisconnectAll() {
	_m.Called()
}

// EstimateGas provides a mock function with given fields: ctx, call
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) EstimateGas(ctx context.Context, call interface{}) (uint64, error) {
	ret := _m.Called(ctx, call)

	if len(ret) == 0 {
		panic("no return value specified for EstimateGas")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) (uint64, error)); ok {
		return rf(ctx, call)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) uint64); ok {
		r0 = rf(ctx, call)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, call)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterEvents provides a mock function with given fields: ctx, query
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for FilterEvents")
	}

	var r0 []EVENT
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, EVENT_OPS) ([]EVENT, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, EVENT_OPS) []EVENT); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EVENT)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, EVENT_OPS) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LINKBalance provides a mock function with given fields: ctx, accountAddress, linkAddress
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	ret := _m.Called(ctx, accountAddress, linkAddress)

	if len(ret) == 0 {
		panic("no return value specified for LINKBalance")
	}

	var r0 *assets.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, ADDR) (*assets.Link, error)); ok {
		return rf(ctx, accountAddress, linkAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, ADDR) *assets.Link); ok {
		r0 = rf(ctx, accountAddress, linkAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR, ADDR) error); ok {
		r1 = rf(ctx, accountAddress, linkAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestBlockHeight provides a mock function with given fields: _a0
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) LatestBlockHeight(_a0 context.Context) (*big.Int, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for LatestBlockHeight")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestFinalizedBlock provides a mock function with given fields: ctx
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) LatestFinalizedBlock(ctx context.Context) (HEAD, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for LatestFinalizedBlock")
	}

	var r0 HEAD
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (HEAD, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) HEAD); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(HEAD)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PendingCallContract provides a mock function with given fields: ctx, msg
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) PendingCallContract(ctx context.Context, msg interface{}) ([]byte, error) {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for PendingCallContract")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) ([]byte, error)); ok {
		return rf(ctx, msg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) []byte); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PendingSequenceAt provides a mock function with given fields: ctx, addr
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for PendingSequenceAt")
	}

	var r0 SEQ
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR) (SEQ, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR) SEQ); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(SEQ)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendEmptyTransaction provides a mock function with given fields: ctx, newTxAttempt, seq, gasLimit, fee, fromAddress
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SendEmptyTransaction(ctx context.Context, newTxAttempt func(SEQ, uint32, FEE, ADDR) (interface{}, error), seq SEQ, gasLimit uint32, fee FEE, fromAddress ADDR) (string, error) {
	ret := _m.Called(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)

	if len(ret) == 0 {
		panic("no return value specified for SendEmptyTransaction")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, func(SEQ, uint32, FEE, ADDR) (interface{}, error), SEQ, uint32, FEE, ADDR) (string, error)); ok {
		return rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, func(SEQ, uint32, FEE, ADDR) (interface{}, error), SEQ, uint32, FEE, ADDR) string); ok {
		r0 = rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, func(SEQ, uint32, FEE, ADDR) (interface{}, error), SEQ, uint32, FEE, ADDR) error); ok {
		r1 = rf(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendTransaction provides a mock function with given fields: ctx, tx
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SendTransaction(ctx context.Context, tx TX) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for SendTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, TX) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SequenceAt provides a mock function with given fields: ctx, accountAddress, blockNumber
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (SEQ, error) {
	ret := _m.Called(ctx, accountAddress, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for SequenceAt")
	}

	var r0 SEQ
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) (SEQ, error)); ok {
		return rf(ctx, accountAddress, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, *big.Int) SEQ); ok {
		r0 = rf(ctx, accountAddress, blockNumber)
	} else {
		r0 = ret.Get(0).(SEQ)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR, *big.Int) error); ok {
		r1 = rf(ctx, accountAddress, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAliveLoopSub provides a mock function with given fields: _a0
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SetAliveLoopSub(_a0 types.Subscription) {
	_m.Called(_a0)
}

// SimulateTransaction provides a mock function with given fields: ctx, tx
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SimulateTransaction(ctx context.Context, tx TX) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for SimulateTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, TX) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: ctx, channel, args
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, channel)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Subscribe")
	}

	var r0 types.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, chan<- HEAD, ...interface{}) (types.Subscription, error)); ok {
		return rf(ctx, channel, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, chan<- HEAD, ...interface{}) types.Subscription); ok {
		r0 = rf(ctx, channel, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, chan<- HEAD, ...interface{}) error); ok {
		r1 = rf(ctx, channel, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribersCount provides a mock function with given fields:
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) SubscribersCount() int32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SubscribersCount")
	}

	var r0 int32
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	return r0
}

// TokenBalance provides a mock function with given fields: ctx, accountAddress, tokenAddress
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error) {
	ret := _m.Called(ctx, accountAddress, tokenAddress)

	if len(ret) == 0 {
		panic("no return value specified for TokenBalance")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, ADDR) (*big.Int, error)); ok {
		return rf(ctx, accountAddress, tokenAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ADDR, ADDR) *big.Int); ok {
		r0 = rf(ctx, accountAddress, tokenAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ADDR, ADDR) error); ok {
		r1 = rf(ctx, accountAddress, tokenAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionByHash provides a mock function with given fields: ctx, txHash
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) TransactionByHash(ctx context.Context, txHash TX_HASH) (TX, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionByHash")
	}

	var r0 TX
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, TX_HASH) (TX, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, TX_HASH) TX); ok {
		r0 = rf(ctx, txHash)
	} else {
		r0 = ret.Get(0).(TX)
	}

	if rf, ok := ret.Get(1).(func(context.Context, TX_HASH) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionReceipt provides a mock function with given fields: ctx, txHash
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (TX_RECEIPT, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionReceipt")
	}

	var r0 TX_RECEIPT
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, TX_HASH) (TX_RECEIPT, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, TX_HASH) TX_RECEIPT); ok {
		r0 = rf(ctx, txHash)
	} else {
		r0 = ret.Get(0).(TX_RECEIPT)
	}

	if rf, ok := ret.Get(1).(func(context.Context, TX_HASH) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeAllExceptAliveLoop provides a mock function with given fields:
func (_m *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]) UnsubscribeAllExceptAliveLoop() {
	_m.Called()
}

// newMockRPC creates a new instance of mockRPC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRPC[CHAIN_ID types.ID, SEQ types.Sequence, ADDR types.Hashable, BLOCK_HASH types.Hashable, TX interface{}, TX_HASH types.Hashable, EVENT interface{}, EVENT_OPS interface{}, TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH], FEE feetypes.Fee, HEAD types.Head[BLOCK_HASH]](t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD] {
	mock := &mockRPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
