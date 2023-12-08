// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	watch "k8s.io/apimachinery/pkg/watch"
)

// Interface is an autogenerated mock type for the Interface type
type Interface struct {
	mock.Mock
}

type Interface_Expecter struct {
	mock *mock.Mock
}

func (_m *Interface) EXPECT() *Interface_Expecter {
	return &Interface_Expecter{mock: &_m.Mock}
}

// ResultChan provides a mock function with given fields:
func (_m *Interface) ResultChan() <-chan watch.Event {
	ret := _m.Called()

	var r0 <-chan watch.Event
	if rf, ok := ret.Get(0).(func() <-chan watch.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan watch.Event)
		}
	}

	return r0
}

// Interface_ResultChan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResultChan'
type Interface_ResultChan_Call struct {
	*mock.Call
}

// ResultChan is a helper method to define mock.On call
func (_e *Interface_Expecter) ResultChan() *Interface_ResultChan_Call {
	return &Interface_ResultChan_Call{Call: _e.mock.On("ResultChan")}
}

func (_c *Interface_ResultChan_Call) Run(run func()) *Interface_ResultChan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Interface_ResultChan_Call) Return(_a0 <-chan watch.Event) *Interface_ResultChan_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Interface_ResultChan_Call) RunAndReturn(run func() <-chan watch.Event) *Interface_ResultChan_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *Interface) Stop() {
	_m.Called()
}

// Interface_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type Interface_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *Interface_Expecter) Stop() *Interface_Stop_Call {
	return &Interface_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *Interface_Stop_Call) Run(run func()) *Interface_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Interface_Stop_Call) Return() *Interface_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *Interface_Stop_Call) RunAndReturn(run func()) *Interface_Stop_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewInterface creates a new instance of Interface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInterface(t mockConstructorTestingTNewInterface) *Interface {
	mock := &Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
