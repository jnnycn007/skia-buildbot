// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	alerts "go.skia.org/infra/perf/go/alerts"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Store) Delete(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, includeDeleted
func (_m *Store) List(ctx context.Context, includeDeleted bool) ([]*alerts.Alert, error) {
	ret := _m.Called(ctx, includeDeleted)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*alerts.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool) ([]*alerts.Alert, error)); ok {
		return rf(ctx, includeDeleted)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool) []*alerts.Alert); ok {
		r0 = rf(ctx, includeDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*alerts.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, includeDeleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForSubscription provides a mock function with given fields: ctx, subName
func (_m *Store) ListForSubscription(ctx context.Context, subName string) ([]*alerts.Alert, error) {
	ret := _m.Called(ctx, subName)

	if len(ret) == 0 {
		panic("no return value specified for ListForSubscription")
	}

	var r0 []*alerts.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*alerts.Alert, error)); ok {
		return rf(ctx, subName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*alerts.Alert); ok {
		r0 = rf(ctx, subName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*alerts.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, subName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplaceAll provides a mock function with given fields: ctx, reqs, tx
func (_m *Store) ReplaceAll(ctx context.Context, reqs []*alerts.SaveRequest, tx pgx.Tx) error {
	ret := _m.Called(ctx, reqs, tx)

	if len(ret) == 0 {
		panic("no return value specified for ReplaceAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*alerts.SaveRequest, pgx.Tx) error); ok {
		r0 = rf(ctx, reqs, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, req
func (_m *Store) Save(ctx context.Context, req *alerts.SaveRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *alerts.SaveRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
