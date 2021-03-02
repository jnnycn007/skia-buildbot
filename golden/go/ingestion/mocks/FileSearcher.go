// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// FileSearcher is an autogenerated mock type for the FileSearcher type
type FileSearcher struct {
	mock.Mock
}

// SearchForFiles provides a mock function with given fields: ctx, start, end
func (_m *FileSearcher) SearchForFiles(ctx context.Context, start time.Time, end time.Time) []string {
	ret := _m.Called(ctx, start, end)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) []string); ok {
		r0 = rf(ctx, start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
