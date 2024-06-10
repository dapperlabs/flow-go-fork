// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	hotstuff "github.com/onflow/flow-go/consensus/hotstuff"
	mock "github.com/stretchr/testify/mock"
)

// TimeoutCollectorFactory is an autogenerated mock type for the TimeoutCollectorFactory type
type TimeoutCollectorFactory struct {
	mock.Mock
}

// Create provides a mock function with given fields: view
func (_m *TimeoutCollectorFactory) Create(view uint64) (hotstuff.TimeoutCollector, error) {
	ret := _m.Called(view)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 hotstuff.TimeoutCollector
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (hotstuff.TimeoutCollector, error)); ok {
		return rf(view)
	}
	if rf, ok := ret.Get(0).(func(uint64) hotstuff.TimeoutCollector); ok {
		r0 = rf(view)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(hotstuff.TimeoutCollector)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(view)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTimeoutCollectorFactory creates a new instance of TimeoutCollectorFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimeoutCollectorFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimeoutCollectorFactory {
	mock := &TimeoutCollectorFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
