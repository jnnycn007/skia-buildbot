// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	chromeperf "go.skia.org/infra/perf/go/chromeperf"

	mock "github.com/stretchr/testify/mock"
)

// ChromePerfClient is an autogenerated mock type for the ChromePerfClient type
type ChromePerfClient struct {
	mock.Mock
}

type ChromePerfClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ChromePerfClient) EXPECT() *ChromePerfClient_Expecter {
	return &ChromePerfClient_Expecter{mock: &_m.Mock}
}

// SendRegression provides a mock function with given fields: ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter
func (_m *ChromePerfClient) SendRegression(ctx context.Context, testPath string, startCommitPosition int32, endCommitPosition int32, projectId string, isImprovement bool, botName string, internal bool, medianBefore float32, medianAfter float32) (*chromeperf.ChromePerfResponse, error) {
	ret := _m.Called(ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter)

	if len(ret) == 0 {
		panic("no return value specified for SendRegression")
	}

	var r0 *chromeperf.ChromePerfResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int32, int32, string, bool, string, bool, float32, float32) (*chromeperf.ChromePerfResponse, error)); ok {
		return rf(ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int32, int32, string, bool, string, bool, float32, float32) *chromeperf.ChromePerfResponse); ok {
		r0 = rf(ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chromeperf.ChromePerfResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int32, int32, string, bool, string, bool, float32, float32) error); ok {
		r1 = rf(ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChromePerfClient_SendRegression_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendRegression'
type ChromePerfClient_SendRegression_Call struct {
	*mock.Call
}

// SendRegression is a helper method to define mock.On call
//   - ctx context.Context
//   - testPath string
//   - startCommitPosition int32
//   - endCommitPosition int32
//   - projectId string
//   - isImprovement bool
//   - botName string
//   - internal bool
//   - medianBefore float32
//   - medianAfter float32
func (_e *ChromePerfClient_Expecter) SendRegression(ctx interface{}, testPath interface{}, startCommitPosition interface{}, endCommitPosition interface{}, projectId interface{}, isImprovement interface{}, botName interface{}, internal interface{}, medianBefore interface{}, medianAfter interface{}) *ChromePerfClient_SendRegression_Call {
	return &ChromePerfClient_SendRegression_Call{Call: _e.mock.On("SendRegression", ctx, testPath, startCommitPosition, endCommitPosition, projectId, isImprovement, botName, internal, medianBefore, medianAfter)}
}

func (_c *ChromePerfClient_SendRegression_Call) Run(run func(ctx context.Context, testPath string, startCommitPosition int32, endCommitPosition int32, projectId string, isImprovement bool, botName string, internal bool, medianBefore float32, medianAfter float32)) *ChromePerfClient_SendRegression_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int32), args[3].(int32), args[4].(string), args[5].(bool), args[6].(string), args[7].(bool), args[8].(float32), args[9].(float32))
	})
	return _c
}

func (_c *ChromePerfClient_SendRegression_Call) Return(_a0 *chromeperf.ChromePerfResponse, _a1 error) *ChromePerfClient_SendRegression_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChromePerfClient_SendRegression_Call) RunAndReturn(run func(context.Context, string, int32, int32, string, bool, string, bool, float32, float32) (*chromeperf.ChromePerfResponse, error)) *ChromePerfClient_SendRegression_Call {
	_c.Call.Return(run)
	return _c
}

// NewChromePerfClient creates a new instance of ChromePerfClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChromePerfClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChromePerfClient {
	mock := &ChromePerfClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
