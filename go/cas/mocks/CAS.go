// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CAS is an autogenerated mock type for the CAS type
type CAS struct {
	mock.Mock
}

type CAS_Expecter struct {
	mock *mock.Mock
}

func (_m *CAS) EXPECT() *CAS_Expecter {
	return &CAS_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *CAS) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CAS_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type CAS_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *CAS_Expecter) Close() *CAS_Close_Call {
	return &CAS_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *CAS_Close_Call) Run(run func()) *CAS_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CAS_Close_Call) Return(_a0 error) *CAS_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CAS_Close_Call) RunAndReturn(run func() error) *CAS_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Download provides a mock function with given fields: ctx, root, digest
func (_m *CAS) Download(ctx context.Context, root string, digest string) error {
	ret := _m.Called(ctx, root, digest)

	if len(ret) == 0 {
		panic("no return value specified for Download")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, root, digest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CAS_Download_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Download'
type CAS_Download_Call struct {
	*mock.Call
}

// Download is a helper method to define mock.On call
//   - ctx context.Context
//   - root string
//   - digest string
func (_e *CAS_Expecter) Download(ctx interface{}, root interface{}, digest interface{}) *CAS_Download_Call {
	return &CAS_Download_Call{Call: _e.mock.On("Download", ctx, root, digest)}
}

func (_c *CAS_Download_Call) Run(run func(ctx context.Context, root string, digest string)) *CAS_Download_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CAS_Download_Call) Return(_a0 error) *CAS_Download_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CAS_Download_Call) RunAndReturn(run func(context.Context, string, string) error) *CAS_Download_Call {
	_c.Call.Return(run)
	return _c
}

// Merge provides a mock function with given fields: ctx, digests
func (_m *CAS) Merge(ctx context.Context, digests []string) (string, error) {
	ret := _m.Called(ctx, digests)

	if len(ret) == 0 {
		panic("no return value specified for Merge")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) (string, error)); ok {
		return rf(ctx, digests)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) string); ok {
		r0 = rf(ctx, digests)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, digests)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CAS_Merge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Merge'
type CAS_Merge_Call struct {
	*mock.Call
}

// Merge is a helper method to define mock.On call
//   - ctx context.Context
//   - digests []string
func (_e *CAS_Expecter) Merge(ctx interface{}, digests interface{}) *CAS_Merge_Call {
	return &CAS_Merge_Call{Call: _e.mock.On("Merge", ctx, digests)}
}

func (_c *CAS_Merge_Call) Run(run func(ctx context.Context, digests []string)) *CAS_Merge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *CAS_Merge_Call) Return(_a0 string, _a1 error) *CAS_Merge_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CAS_Merge_Call) RunAndReturn(run func(context.Context, []string) (string, error)) *CAS_Merge_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: ctx, root, paths, excludes
func (_m *CAS) Upload(ctx context.Context, root string, paths []string, excludes []string) (string, error) {
	ret := _m.Called(ctx, root, paths, excludes)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, []string) (string, error)); ok {
		return rf(ctx, root, paths, excludes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, []string) string); ok {
		r0 = rf(ctx, root, paths, excludes)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []string, []string) error); ok {
		r1 = rf(ctx, root, paths, excludes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CAS_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type CAS_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - ctx context.Context
//   - root string
//   - paths []string
//   - excludes []string
func (_e *CAS_Expecter) Upload(ctx interface{}, root interface{}, paths interface{}, excludes interface{}) *CAS_Upload_Call {
	return &CAS_Upload_Call{Call: _e.mock.On("Upload", ctx, root, paths, excludes)}
}

func (_c *CAS_Upload_Call) Run(run func(ctx context.Context, root string, paths []string, excludes []string)) *CAS_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]string), args[3].([]string))
	})
	return _c
}

func (_c *CAS_Upload_Call) Return(_a0 string, _a1 error) *CAS_Upload_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CAS_Upload_Call) RunAndReturn(run func(context.Context, string, []string, []string) (string, error)) *CAS_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewCAS creates a new instance of CAS. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCAS(t interface {
	mock.TestingT
	Cleanup(func())
}) *CAS {
	mock := &CAS{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
