// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ExecutionDataRequesterV2Metrics is an autogenerated mock type for the ExecutionDataRequesterV2Metrics type
type ExecutionDataRequesterV2Metrics struct {
	mock.Mock
}

// FulfilledHeight provides a mock function with given fields: blockHeight
func (_m *ExecutionDataRequesterV2Metrics) FulfilledHeight(blockHeight uint64) {
	_m.Called(blockHeight)
}

// ReceiptSkipped provides a mock function with given fields:
func (_m *ExecutionDataRequesterV2Metrics) ReceiptSkipped() {
	_m.Called()
}

// RequestCanceled provides a mock function with given fields:
func (_m *ExecutionDataRequesterV2Metrics) RequestCanceled() {
	_m.Called()
}

// RequestFailed provides a mock function with given fields: duration, retryable
func (_m *ExecutionDataRequesterV2Metrics) RequestFailed(duration time.Duration, retryable bool) {
	_m.Called(duration, retryable)
}

// RequestSucceeded provides a mock function with given fields: blockHeight, duration, totalSize, numberOfAttempts
func (_m *ExecutionDataRequesterV2Metrics) RequestSucceeded(blockHeight uint64, duration time.Duration, totalSize uint64, numberOfAttempts int) {
	_m.Called(blockHeight, duration, totalSize, numberOfAttempts)
}

// ResponseDropped provides a mock function with given fields:
func (_m *ExecutionDataRequesterV2Metrics) ResponseDropped() {
	_m.Called()
}

// NewExecutionDataRequesterV2Metrics creates a new instance of ExecutionDataRequesterV2Metrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutionDataRequesterV2Metrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExecutionDataRequesterV2Metrics {
	mock := &ExecutionDataRequesterV2Metrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
