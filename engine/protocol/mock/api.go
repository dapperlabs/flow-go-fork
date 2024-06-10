// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	access "github.com/onflow/flow-go/access"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// GetBlockByHeight provides a mock function with given fields: ctx, height
func (_m *API) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockByHeight")
	}

	var r0 *flow.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*flow.Block, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *flow.Block); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByID provides a mock function with given fields: ctx, id
func (_m *API) GetBlockByID(ctx context.Context, id flow.Identifier) (*flow.Block, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockByID")
	}

	var r0 *flow.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.Block, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.Block); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeaderByHeight provides a mock function with given fields: ctx, height
func (_m *API) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.Header, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockHeaderByHeight")
	}

	var r0 *flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*flow.Header, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *flow.Header); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeaderByID provides a mock function with given fields: ctx, id
func (_m *API) GetBlockHeaderByID(ctx context.Context, id flow.Identifier) (*flow.Header, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockHeaderByID")
	}

	var r0 *flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) (*flow.Header, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *flow.Header); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlock provides a mock function with given fields: ctx, isSealed
func (_m *API) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	ret := _m.Called(ctx, isSealed)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlock")
	}

	var r0 *flow.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool) (*flow.Block, error)); ok {
		return rf(ctx, isSealed)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool) *flow.Block); ok {
		r0 = rf(ctx, isSealed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, isSealed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlockHeader provides a mock function with given fields: ctx, isSealed
func (_m *API) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.Header, error) {
	ret := _m.Called(ctx, isSealed)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockHeader")
	}

	var r0 *flow.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool) (*flow.Header, error)); ok {
		return rf(ctx, isSealed)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool) *flow.Header); ok {
		r0 = rf(ctx, isSealed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, isSealed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestProtocolStateSnapshot provides a mock function with given fields: ctx
func (_m *API) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestProtocolStateSnapshot")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]byte, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetworkParameters provides a mock function with given fields: ctx
func (_m *API) GetNetworkParameters(ctx context.Context) access.NetworkParameters {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetNetworkParameters")
	}

	var r0 access.NetworkParameters
	if rf, ok := ret.Get(0).(func(context.Context) access.NetworkParameters); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(access.NetworkParameters)
	}

	return r0
}

// GetNodeVersionInfo provides a mock function with given fields: ctx
func (_m *API) GetNodeVersionInfo(ctx context.Context) (*access.NodeVersionInfo, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetNodeVersionInfo")
	}

	var r0 *access.NodeVersionInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*access.NodeVersionInfo, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *access.NodeVersionInfo); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.NodeVersionInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProtocolStateSnapshotByBlockID provides a mock function with given fields: ctx, blockID
func (_m *API) GetProtocolStateSnapshotByBlockID(ctx context.Context, blockID flow.Identifier) ([]byte, error) {
	ret := _m.Called(ctx, blockID)

	if len(ret) == 0 {
		panic("no return value specified for GetProtocolStateSnapshotByBlockID")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) ([]byte, error)); ok {
		return rf(ctx, blockID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) []byte); ok {
		r0 = rf(ctx, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProtocolStateSnapshotByHeight provides a mock function with given fields: ctx, blockHeight
func (_m *API) GetProtocolStateSnapshotByHeight(ctx context.Context, blockHeight uint64) ([]byte, error) {
	ret := _m.Called(ctx, blockHeight)

	if len(ret) == 0 {
		panic("no return value specified for GetProtocolStateSnapshotByHeight")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) ([]byte, error)); ok {
		return rf(ctx, blockHeight)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) []byte); ok {
		r0 = rf(ctx, blockHeight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, blockHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
