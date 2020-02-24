// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	digesttools "go.skia.org/infra/golden/go/digesttools"
	expectations "go.skia.org/infra/golden/go/expectations"

	mock "github.com/stretchr/testify/mock"

	types "go.skia.org/infra/golden/go/types"
)

// ClosestDiffFinder is an autogenerated mock type for the ClosestDiffFinder type
type ClosestDiffFinder struct {
	mock.Mock
}

// ClosestDigest provides a mock function with given fields: ctx, test, digest, label
func (_m *ClosestDiffFinder) ClosestDigest(ctx context.Context, test types.TestName, digest types.Digest, label expectations.Label) (*digesttools.Closest, error) {
	ret := _m.Called(ctx, test, digest, label)

	var r0 *digesttools.Closest
	if rf, ok := ret.Get(0).(func(context.Context, types.TestName, types.Digest, expectations.Label) *digesttools.Closest); ok {
		r0 = rf(ctx, test, digest, label)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*digesttools.Closest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.TestName, types.Digest, expectations.Label) error); ok {
		r1 = rf(ctx, test, digest, label)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Precompute provides a mock function with given fields: ctx
func (_m *ClosestDiffFinder) Precompute(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
