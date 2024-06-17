// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// L1Oracle is an autogenerated mock type for the L1Oracle type
type L1Oracle struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *L1Oracle) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GasPrice provides a mock function with given fields: ctx
func (_m *L1Oracle) GasPrice(ctx context.Context) (*assets.Wei, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GasPrice")
	}

	var r0 *assets.Wei
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*assets.Wei, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *assets.Wei); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGasCost provides a mock function with given fields: ctx, tx, blockNum
func (_m *L1Oracle) GetGasCost(ctx context.Context, tx *types.Transaction, blockNum *big.Int) (*assets.Wei, error) {
	ret := _m.Called(ctx, tx, blockNum)

	if len(ret) == 0 {
		panic("no return value specified for GetGasCost")
	}

	var r0 *assets.Wei
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Transaction, *big.Int) (*assets.Wei, error)); ok {
		return rf(ctx, tx, blockNum)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.Transaction, *big.Int) *assets.Wei); ok {
		r0 = rf(ctx, tx, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.Transaction, *big.Int) error); ok {
		r1 = rf(ctx, tx, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthReport provides a mock function with given fields:
func (_m *L1Oracle) HealthReport() map[string]error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HealthReport")
	}

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *L1Oracle) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *L1Oracle) Ready() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *L1Oracle) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewL1Oracle creates a new instance of L1Oracle. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewL1Oracle(t interface {
	mock.TestingT
	Cleanup(func())
}) *L1Oracle {
	mock := &L1Oracle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
