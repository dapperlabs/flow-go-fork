// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// ScriptExecutor is an autogenerated mock type for the ScriptExecutor type
type ScriptExecutor struct {
	mock.Mock
}

// ExecuteScriptAtBlockID provides a mock function with given fields: ctx, script, arguments, blockID
func (_m *ScriptExecutor) ExecuteScriptAtBlockID(ctx context.Context, script []byte, arguments [][]byte, blockID flow.Identifier) ([]byte, uint64, error) {
	ret := _m.Called(ctx, script, arguments, blockID)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteScriptAtBlockID")
	}

	var r0 []byte
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, flow.Identifier) ([]byte, uint64, error)); ok {
		return rf(ctx, script, arguments, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, [][]byte, flow.Identifier) []byte); ok {
		r0 = rf(ctx, script, arguments, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, [][]byte, flow.Identifier) uint64); ok {
		r1 = rf(ctx, script, arguments, blockID)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, []byte, [][]byte, flow.Identifier) error); ok {
		r2 = rf(ctx, script, arguments, blockID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAccount provides a mock function with given fields: ctx, address, blockID
func (_m *ScriptExecutor) GetAccount(ctx context.Context, address flow.Address, blockID flow.Identifier) (*flow.Account, error) {
	ret := _m.Called(ctx, address, blockID)

	if len(ret) == 0 {
		panic("no return value specified for GetAccount")
	}

	var r0 *flow.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, flow.Identifier) (*flow.Account, error)); ok {
		return rf(ctx, address, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Address, flow.Identifier) *flow.Account); ok {
		r0 = rf(ctx, address, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Address, flow.Identifier) error); ok {
		r1 = rf(ctx, address, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterAtBlockID provides a mock function with given fields: ctx, owner, key, blockID
func (_m *ScriptExecutor) GetRegisterAtBlockID(ctx context.Context, owner []byte, key []byte, blockID flow.Identifier) ([]byte, error) {
	ret := _m.Called(ctx, owner, key, blockID)

	if len(ret) == 0 {
		panic("no return value specified for GetRegisterAtBlockID")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte, flow.Identifier) ([]byte, error)); ok {
		return rf(ctx, owner, key, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte, flow.Identifier) []byte); ok {
		r0 = rf(ctx, owner, key, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, []byte, flow.Identifier) error); ok {
		r1 = rf(ctx, owner, key, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScriptExecutor creates a new instance of ScriptExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScriptExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScriptExecutor {
	mock := &ScriptExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
