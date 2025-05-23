// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	gopubsub "cloud.google.com/go/pubsub"
	mock "github.com/stretchr/testify/mock"
)

// Snapshot is an autogenerated mock type for the Snapshot type
type Snapshot struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx
func (_m *Snapshot) Delete(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ID provides a mock function with no fields
func (_m *Snapshot) ID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetLabels provides a mock function with given fields: ctx, label
func (_m *Snapshot) SetLabels(ctx context.Context, label map[string]string) (*gopubsub.SnapshotConfig, error) {
	ret := _m.Called(ctx, label)

	if len(ret) == 0 {
		panic("no return value specified for SetLabels")
	}

	var r0 *gopubsub.SnapshotConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) (*gopubsub.SnapshotConfig, error)); ok {
		return rf(ctx, label)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) *gopubsub.SnapshotConfig); ok {
		r0 = rf(ctx, label)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gopubsub.SnapshotConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]string) error); ok {
		r1 = rf(ctx, label)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSnapshot creates a new instance of Snapshot. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSnapshot(t interface {
	mock.TestingT
	Cleanup(func())
}) *Snapshot {
	mock := &Snapshot{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
