// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	types "go.skia.org/infra/skcq/go/types"
)

// CurrentChangesCache is an autogenerated mock type for the CurrentChangesCache type
type CurrentChangesCache struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, changeEquivalentPatchset, changeSubject, changeOwner, repo, branch, dryRun, internal, changeID, latestPatchsetID
func (_m *CurrentChangesCache) Add(ctx context.Context, changeEquivalentPatchset string, changeSubject string, changeOwner string, repo string, branch string, dryRun bool, internal bool, changeID int64, latestPatchsetID int64) (int64, bool, error) {
	ret := _m.Called(ctx, changeEquivalentPatchset, changeSubject, changeOwner, repo, branch, dryRun, internal, changeID, latestPatchsetID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, bool, bool, int64, int64) int64); ok {
		r0 = rf(ctx, changeEquivalentPatchset, changeSubject, changeOwner, repo, branch, dryRun, internal, changeID, latestPatchsetID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string, bool, bool, int64, int64) bool); ok {
		r1 = rf(ctx, changeEquivalentPatchset, changeSubject, changeOwner, repo, branch, dryRun, internal, changeID, latestPatchsetID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, string, string, string, bool, bool, int64, int64) error); ok {
		r2 = rf(ctx, changeEquivalentPatchset, changeSubject, changeOwner, repo, branch, dryRun, internal, changeID, latestPatchsetID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Get provides a mock function with given fields:
func (_m *CurrentChangesCache) Get() map[string]*types.CurrentlyProcessingChange {
	ret := _m.Called()

	var r0 map[string]*types.CurrentlyProcessingChange
	if rf, ok := ret.Get(0).(func() map[string]*types.CurrentlyProcessingChange); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*types.CurrentlyProcessingChange)
		}
	}

	return r0
}

// Remove provides a mock function with given fields: ctx, changeEquivalentPatchset
func (_m *CurrentChangesCache) Remove(ctx context.Context, changeEquivalentPatchset string) error {
	ret := _m.Called(ctx, changeEquivalentPatchset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, changeEquivalentPatchset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
