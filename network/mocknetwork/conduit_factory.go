// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocknetwork

import (
	context "context"

	network "github.com/onflow/flow-go/network"
	mock "github.com/stretchr/testify/mock"
)

// ConduitFactory is an autogenerated mock type for the ConduitFactory type
type ConduitFactory struct {
	mock.Mock
}

// NewConduit provides a mock function with given fields: _a0, _a1
func (_m *ConduitFactory) NewConduit(_a0 context.Context, _a1 network.Channel) (network.Conduit, error) {
	ret := _m.Called(_a0, _a1)

	var r0 network.Conduit
	if rf, ok := ret.Get(0).(func(context.Context, network.Channel) network.Conduit); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conduit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, network.Channel) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterAdapter provides a mock function with given fields: _a0
func (_m *ConduitFactory) RegisterAdapter(_a0 network.Adapter) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(network.Adapter) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewConduitFactoryT interface {
	mock.TestingT
	Cleanup(func())
}

// NewConduitFactory creates a new instance of ConduitFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConduitFactory(t NewConduitFactoryT) *ConduitFactory {
	mock := &ConduitFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
