// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	protocol "github.com/onflow/flow-go/state/protocol"
)

// ClusterRootQCVoter is an autogenerated mock type for the ClusterRootQCVoter type
type ClusterRootQCVoter struct {
	mock.Mock
}

// Vote provides a mock function with given fields: _a0, _a1
func (_m *ClusterRootQCVoter) Vote(_a0 context.Context, _a1 protocol.Epoch) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, protocol.Epoch) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClusterRootQCVoter interface {
	mock.TestingT
	Cleanup(func())
}

// NewClusterRootQCVoter creates a new instance of ClusterRootQCVoter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClusterRootQCVoter(t mockConstructorTestingTNewClusterRootQCVoter) *ClusterRootQCVoter {
	mock := &ClusterRootQCVoter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
