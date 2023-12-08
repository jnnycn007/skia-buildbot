// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// PeerFinder is an autogenerated mock type for the PeerFinder type
type PeerFinder struct {
	mock.Mock
}

type PeerFinder_Expecter struct {
	mock *mock.Mock
}

func (_m *PeerFinder) EXPECT() *PeerFinder_Expecter {
	return &PeerFinder_Expecter{mock: &_m.Mock}
}

// Start provides a mock function with given fields: ctx
func (_m *PeerFinder) Start(ctx context.Context) ([]string, <-chan []string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 []string
	var r1 <-chan []string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, <-chan []string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) <-chan []string); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan []string)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PeerFinder_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type PeerFinder_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PeerFinder_Expecter) Start(ctx interface{}) *PeerFinder_Start_Call {
	return &PeerFinder_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *PeerFinder_Start_Call) Run(run func(ctx context.Context)) *PeerFinder_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PeerFinder_Start_Call) Return(_a0 []string, _a1 <-chan []string, _a2 error) *PeerFinder_Start_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *PeerFinder_Start_Call) RunAndReturn(run func(context.Context) ([]string, <-chan []string, error)) *PeerFinder_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewPeerFinder creates a new instance of PeerFinder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPeerFinder(t interface {
	mock.TestingT
	Cleanup(func())
}) *PeerFinder {
	mock := &PeerFinder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
