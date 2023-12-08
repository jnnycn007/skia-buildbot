// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	autoroll "go.skia.org/infra/go/autoroll"

	mock "github.com/stretchr/testify/mock"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

type DB_Expecter struct {
	mock *mock.Mock
}

func (_m *DB) EXPECT() *DB_Expecter {
	return &DB_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, roller, issue
func (_m *DB) Get(ctx context.Context, roller string, issue int64) (*autoroll.AutoRollIssue, error) {
	ret := _m.Called(ctx, roller, issue)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *autoroll.AutoRollIssue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (*autoroll.AutoRollIssue, error)); ok {
		return rf(ctx, roller, issue)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) *autoroll.AutoRollIssue); ok {
		r0 = rf(ctx, roller, issue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*autoroll.AutoRollIssue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, roller, issue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DB_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type DB_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - roller string
//   - issue int64
func (_e *DB_Expecter) Get(ctx interface{}, roller interface{}, issue interface{}) *DB_Get_Call {
	return &DB_Get_Call{Call: _e.mock.On("Get", ctx, roller, issue)}
}

func (_c *DB_Get_Call) Run(run func(ctx context.Context, roller string, issue int64)) *DB_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int64))
	})
	return _c
}

func (_c *DB_Get_Call) Return(_a0 *autoroll.AutoRollIssue, _a1 error) *DB_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DB_Get_Call) RunAndReturn(run func(context.Context, string, int64) (*autoroll.AutoRollIssue, error)) *DB_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetRolls provides a mock function with given fields: ctx, roller, cursor
func (_m *DB) GetRolls(ctx context.Context, roller string, cursor string) ([]*autoroll.AutoRollIssue, string, error) {
	ret := _m.Called(ctx, roller, cursor)

	if len(ret) == 0 {
		panic("no return value specified for GetRolls")
	}

	var r0 []*autoroll.AutoRollIssue
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]*autoroll.AutoRollIssue, string, error)); ok {
		return rf(ctx, roller, cursor)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*autoroll.AutoRollIssue); ok {
		r0 = rf(ctx, roller, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*autoroll.AutoRollIssue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, roller, cursor)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, roller, cursor)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DB_GetRolls_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRolls'
type DB_GetRolls_Call struct {
	*mock.Call
}

// GetRolls is a helper method to define mock.On call
//   - ctx context.Context
//   - roller string
//   - cursor string
func (_e *DB_Expecter) GetRolls(ctx interface{}, roller interface{}, cursor interface{}) *DB_GetRolls_Call {
	return &DB_GetRolls_Call{Call: _e.mock.On("GetRolls", ctx, roller, cursor)}
}

func (_c *DB_GetRolls_Call) Run(run func(ctx context.Context, roller string, cursor string)) *DB_GetRolls_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *DB_GetRolls_Call) Return(_a0 []*autoroll.AutoRollIssue, _a1 string, _a2 error) *DB_GetRolls_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *DB_GetRolls_Call) RunAndReturn(run func(context.Context, string, string) ([]*autoroll.AutoRollIssue, string, error)) *DB_GetRolls_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with given fields: ctx, roller, roll
func (_m *DB) Put(ctx context.Context, roller string, roll *autoroll.AutoRollIssue) error {
	ret := _m.Called(ctx, roller, roll)

	if len(ret) == 0 {
		panic("no return value specified for Put")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *autoroll.AutoRollIssue) error); ok {
		r0 = rf(ctx, roller, roll)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type DB_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
//   - ctx context.Context
//   - roller string
//   - roll *autoroll.AutoRollIssue
func (_e *DB_Expecter) Put(ctx interface{}, roller interface{}, roll interface{}) *DB_Put_Call {
	return &DB_Put_Call{Call: _e.mock.On("Put", ctx, roller, roll)}
}

func (_c *DB_Put_Call) Run(run func(ctx context.Context, roller string, roll *autoroll.AutoRollIssue)) *DB_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*autoroll.AutoRollIssue))
	})
	return _c
}

func (_c *DB_Put_Call) Return(_a0 error) *DB_Put_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DB_Put_Call) RunAndReturn(run func(context.Context, string, *autoroll.AutoRollIssue) error) *DB_Put_Call {
	_c.Call.Return(run)
	return _c
}

// NewDB creates a new instance of DB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *DB {
	mock := &DB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
