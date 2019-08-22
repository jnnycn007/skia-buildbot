// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	continuous_integration "go.skia.org/infra/golden/go/continuous_integration"

	tjstore "go.skia.org/infra/golden/go/tjstore"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// GetResults provides a mock function with given fields: ctx, psID
func (_m *Store) GetResults(ctx context.Context, psID tjstore.CombinedPSID) ([]tjstore.TryJobResult, error) {
	ret := _m.Called(ctx, psID)

	var r0 []tjstore.TryJobResult
	if rf, ok := ret.Get(0).(func(context.Context, tjstore.CombinedPSID) []tjstore.TryJobResult); ok {
		r0 = rf(ctx, psID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tjstore.TryJobResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, tjstore.CombinedPSID) error); ok {
		r1 = rf(ctx, psID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRunningTryJobs provides a mock function with given fields: ctx
func (_m *Store) GetRunningTryJobs(ctx context.Context) ([]continuous_integration.TryJob, error) {
	ret := _m.Called(ctx)

	var r0 []continuous_integration.TryJob
	if rf, ok := ret.Get(0).(func(context.Context) []continuous_integration.TryJob); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]continuous_integration.TryJob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTryJob provides a mock function with given fields: ctx, id
func (_m *Store) GetTryJob(ctx context.Context, id string) (continuous_integration.TryJob, error) {
	ret := _m.Called(ctx, id)

	var r0 continuous_integration.TryJob
	if rf, ok := ret.Get(0).(func(context.Context, string) continuous_integration.TryJob); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(continuous_integration.TryJob)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutResults provides a mock function with given fields: ctx, psID, r
func (_m *Store) PutResults(ctx context.Context, psID tjstore.CombinedPSID, r []tjstore.TryJobResult) error {
	ret := _m.Called(ctx, psID, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, tjstore.CombinedPSID, []tjstore.TryJobResult) error); ok {
		r0 = rf(ctx, psID, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutTryJob provides a mock function with given fields: ctx, psID, tj
func (_m *Store) PutTryJob(ctx context.Context, psID tjstore.CombinedPSID, tj continuous_integration.TryJob) error {
	ret := _m.Called(ctx, psID, tj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, tjstore.CombinedPSID, continuous_integration.TryJob) error); ok {
		r0 = rf(ctx, psID, tj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
