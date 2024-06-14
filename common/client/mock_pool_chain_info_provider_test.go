// Code generated by mockery v2.42.2. DO NOT EDIT.

package client

import mock "github.com/stretchr/testify/mock"

// mockPoolChainInfoProvider is an autogenerated mock type for the PoolChainInfoProvider type
type mockPoolChainInfoProvider struct {
	mock.Mock
}

// HighestUserObservations provides a mock function with given fields:
func (_m *mockPoolChainInfoProvider) HighestUserObservations() ChainInfo {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HighestUserObservations")
	}

	var r0 ChainInfo
	if rf, ok := ret.Get(0).(func() ChainInfo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ChainInfo)
	}

	return r0
}

// LatestChainInfo provides a mock function with given fields:
func (_m *mockPoolChainInfoProvider) LatestChainInfo() (int, ChainInfo) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LatestChainInfo")
	}

	var r0 int
	var r1 ChainInfo
	if rf, ok := ret.Get(0).(func() (int, ChainInfo)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() ChainInfo); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(ChainInfo)
	}

	return r0, r1
}

// newMockPoolChainInfoProvider creates a new instance of mockPoolChainInfoProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockPoolChainInfoProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockPoolChainInfoProvider {
	mock := &mockPoolChainInfoProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}