// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// TimeTicker is an autogenerated mock type for the TimeTicker type
type TimeTicker struct {
	mock.Mock
}

type TimeTicker_Expecter struct {
	mock *mock.Mock
}

func (_m *TimeTicker) EXPECT() *TimeTicker_Expecter {
	return &TimeTicker_Expecter{mock: &_m.Mock}
}

// C provides a mock function with given fields:
func (_m *TimeTicker) C() <-chan time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for C")
	}

	var r0 <-chan time.Time
	if rf, ok := ret.Get(0).(func() <-chan time.Time); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan time.Time)
		}
	}

	return r0
}

// TimeTicker_C_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'C'
type TimeTicker_C_Call struct {
	*mock.Call
}

// C is a helper method to define mock.On call
func (_e *TimeTicker_Expecter) C() *TimeTicker_C_Call {
	return &TimeTicker_C_Call{Call: _e.mock.On("C")}
}

func (_c *TimeTicker_C_Call) Run(run func()) *TimeTicker_C_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TimeTicker_C_Call) Return(_a0 <-chan time.Time) *TimeTicker_C_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TimeTicker_C_Call) RunAndReturn(run func() <-chan time.Time) *TimeTicker_C_Call {
	_c.Call.Return(run)
	return _c
}

// Reset provides a mock function with given fields: d
func (_m *TimeTicker) Reset(d time.Duration) {
	_m.Called(d)
}

// TimeTicker_Reset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reset'
type TimeTicker_Reset_Call struct {
	*mock.Call
}

// Reset is a helper method to define mock.On call
//   - d time.Duration
func (_e *TimeTicker_Expecter) Reset(d interface{}) *TimeTicker_Reset_Call {
	return &TimeTicker_Reset_Call{Call: _e.mock.On("Reset", d)}
}

func (_c *TimeTicker_Reset_Call) Run(run func(d time.Duration)) *TimeTicker_Reset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *TimeTicker_Reset_Call) Return() *TimeTicker_Reset_Call {
	_c.Call.Return()
	return _c
}

func (_c *TimeTicker_Reset_Call) RunAndReturn(run func(time.Duration)) *TimeTicker_Reset_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *TimeTicker) Stop() {
	_m.Called()
}

// TimeTicker_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type TimeTicker_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *TimeTicker_Expecter) Stop() *TimeTicker_Stop_Call {
	return &TimeTicker_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *TimeTicker_Stop_Call) Run(run func()) *TimeTicker_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TimeTicker_Stop_Call) Return() *TimeTicker_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *TimeTicker_Stop_Call) RunAndReturn(run func()) *TimeTicker_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// NewTimeTicker creates a new instance of TimeTicker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimeTicker(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimeTicker {
	mock := &TimeTicker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
