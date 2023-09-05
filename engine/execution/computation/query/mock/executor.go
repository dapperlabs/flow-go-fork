// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	snapshot "github.com/onflow/flow-go/fvm/storage/snapshot"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// ExecuteScript provides a mock function with given fields: ctx, script, arguments, blockHeader, _a4
func (_m *Executor) ExecuteScript(ctx context.Context, script []byte, arguments [][]byte, blockHeader *flow.Header, _a4 snapshot.StorageSnapshot) ([]byte, error) {
	ret := _m.Called(ctx, script, arguments, blockHeader, _a4)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, *flow.Header, snapshot.StorageSnapshot) ([]byte, error)); ok {
		return rf(ctx, script, arguments, blockHeader, _a4)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, *flow.Header, snapshot.StorageSnapshot) []byte); ok {
		r0 = rf(ctx, script, arguments, blockHeader, _a4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, [][]byte, *flow.Header, snapshot.StorageSnapshot) error); ok {
		r1 = rf(ctx, script, arguments, blockHeader, _a4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccount provides a mock function with given fields: ctx, addr, header, _a3
func (_m *Executor) GetAccount(ctx context.Context, addr flow.Address, header *flow.Header, _a3 snapshot.StorageSnapshot) (*flow.Account, error) {
	ret := _m.Called(ctx, addr, header, _a3)

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, *flow.Header, snapshot.StorageSnapshot) (*flow.Account, error)); ok {
		return rf(ctx, addr, header, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, *flow.Header, snapshot.StorageSnapshot) *flow.Account); ok {
		r0 = rf(ctx, addr, header, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, *flow.Header, snapshot.StorageSnapshot) error); ok {
		r1 = rf(ctx, addr, header, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutor creates a new instance of Executor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutor(t mockConstructorTestingTNewExecutor) *Executor {
	mock := &Executor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}