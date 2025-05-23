// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	types "go.skia.org/infra/task_scheduler/go/types"
)

// JobCache is an autogenerated mock type for the JobCache type
type JobCache struct {
	mock.Mock
}

// AddJobs provides a mock function with given fields: _a0
func (_m *JobCache) AddJobs(_a0 []*types.Job) {
	_m.Called(_a0)
}

// GetAllCachedJobs provides a mock function with no fields
func (_m *JobCache) GetAllCachedJobs() []*types.Job {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllCachedJobs")
	}

	var r0 []*types.Job
	if rf, ok := ret.Get(0).(func() []*types.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Job)
		}
	}

	return r0
}

// GetJob provides a mock function with given fields: _a0
func (_m *JobCache) GetJob(_a0 string) (*types.Job, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetJob")
	}

	var r0 *types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*types.Job, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *types.Job); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobMaybeExpired provides a mock function with given fields: _a0, _a1
func (_m *JobCache) GetJobMaybeExpired(_a0 context.Context, _a1 string) (*types.Job, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetJobMaybeExpired")
	}

	var r0 *types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*types.Job, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.Job); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobsByRepoState provides a mock function with given fields: _a0, _a1
func (_m *JobCache) GetJobsByRepoState(_a0 string, _a1 types.RepoState) ([]*types.Job, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetJobsByRepoState")
	}

	var r0 []*types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(string, types.RepoState) ([]*types.Job, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, types.RepoState) []*types.Job); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(string, types.RepoState) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobsFromDateRange provides a mock function with given fields: _a0, _a1
func (_m *JobCache) GetJobsFromDateRange(_a0 time.Time, _a1 time.Time) ([]*types.Job, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetJobsFromDateRange")
	}

	var r0 []*types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) ([]*types.Job, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) []*types.Job); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMatchingJobsFromDateRange provides a mock function with given fields: names, from, to
func (_m *JobCache) GetMatchingJobsFromDateRange(names []string, from time.Time, to time.Time) (map[string][]*types.Job, error) {
	ret := _m.Called(names, from, to)

	if len(ret) == 0 {
		panic("no return value specified for GetMatchingJobsFromDateRange")
	}

	var r0 map[string][]*types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func([]string, time.Time, time.Time) (map[string][]*types.Job, error)); ok {
		return rf(names, from, to)
	}
	if rf, ok := ret.Get(0).(func([]string, time.Time, time.Time) map[string][]*types.Job); ok {
		r0 = rf(names, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func([]string, time.Time, time.Time) error); ok {
		r1 = rf(names, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InProgressJobs provides a mock function with no fields
func (_m *JobCache) InProgressJobs() ([]*types.Job, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for InProgressJobs")
	}

	var r0 []*types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*types.Job, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*types.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LastUpdated provides a mock function with no fields
func (_m *JobCache) LastUpdated() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LastUpdated")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// RequestedJobs provides a mock function with no fields
func (_m *JobCache) RequestedJobs() ([]*types.Job, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RequestedJobs")
	}

	var r0 []*types.Job
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*types.Job, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*types.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Job)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx
func (_m *JobCache) Update(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewJobCache creates a new instance of JobCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJobCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *JobCache {
	mock := &JobCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
