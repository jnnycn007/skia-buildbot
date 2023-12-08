// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	config "go.skia.org/infra/perf/go/config"

	tracestore "go.skia.org/infra/perf/go/tracestore"

	types "go.skia.org/infra/perf/go/types"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

type Application_Expecter struct {
	mock *mock.Mock
}

func (_m *Application) EXPECT() *Application_Expecter {
	return &Application_Expecter{mock: &_m.Mock}
}

// ConfigCreatePubSubTopicsAndSubscriptions provides a mock function with given fields: instanceConfig
func (_m *Application) ConfigCreatePubSubTopicsAndSubscriptions(instanceConfig *config.InstanceConfig) error {
	ret := _m.Called(instanceConfig)

	if len(ret) == 0 {
		panic("no return value specified for ConfigCreatePubSubTopicsAndSubscriptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*config.InstanceConfig) error); ok {
		r0 = rf(instanceConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_ConfigCreatePubSubTopicsAndSubscriptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigCreatePubSubTopicsAndSubscriptions'
type Application_ConfigCreatePubSubTopicsAndSubscriptions_Call struct {
	*mock.Call
}

// ConfigCreatePubSubTopicsAndSubscriptions is a helper method to define mock.On call
//   - instanceConfig *config.InstanceConfig
func (_e *Application_Expecter) ConfigCreatePubSubTopicsAndSubscriptions(instanceConfig interface{}) *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call {
	return &Application_ConfigCreatePubSubTopicsAndSubscriptions_Call{Call: _e.mock.On("ConfigCreatePubSubTopicsAndSubscriptions", instanceConfig)}
}

func (_c *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call) Run(run func(instanceConfig *config.InstanceConfig)) *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*config.InstanceConfig))
	})
	return _c
}

func (_c *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call) Return(_a0 error) *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call) RunAndReturn(run func(*config.InstanceConfig) error) *Application_ConfigCreatePubSubTopicsAndSubscriptions_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseBackupAlerts provides a mock function with given fields: local, instanceConfig, outputFile
func (_m *Application) DatabaseBackupAlerts(local bool, instanceConfig *config.InstanceConfig, outputFile string) error {
	ret := _m.Called(local, instanceConfig, outputFile)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseBackupAlerts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string) error); ok {
		r0 = rf(local, instanceConfig, outputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseBackupAlerts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseBackupAlerts'
type Application_DatabaseBackupAlerts_Call struct {
	*mock.Call
}

// DatabaseBackupAlerts is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - outputFile string
func (_e *Application_Expecter) DatabaseBackupAlerts(local interface{}, instanceConfig interface{}, outputFile interface{}) *Application_DatabaseBackupAlerts_Call {
	return &Application_DatabaseBackupAlerts_Call{Call: _e.mock.On("DatabaseBackupAlerts", local, instanceConfig, outputFile)}
}

func (_c *Application_DatabaseBackupAlerts_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, outputFile string)) *Application_DatabaseBackupAlerts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string))
	})
	return _c
}

func (_c *Application_DatabaseBackupAlerts_Call) Return(_a0 error) *Application_DatabaseBackupAlerts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseBackupAlerts_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string) error) *Application_DatabaseBackupAlerts_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseBackupRegressions provides a mock function with given fields: local, instanceConfig, outputFile, backupTo
func (_m *Application) DatabaseBackupRegressions(local bool, instanceConfig *config.InstanceConfig, outputFile string, backupTo string) error {
	ret := _m.Called(local, instanceConfig, outputFile, backupTo)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseBackupRegressions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string, string) error); ok {
		r0 = rf(local, instanceConfig, outputFile, backupTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseBackupRegressions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseBackupRegressions'
type Application_DatabaseBackupRegressions_Call struct {
	*mock.Call
}

// DatabaseBackupRegressions is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - outputFile string
//   - backupTo string
func (_e *Application_Expecter) DatabaseBackupRegressions(local interface{}, instanceConfig interface{}, outputFile interface{}, backupTo interface{}) *Application_DatabaseBackupRegressions_Call {
	return &Application_DatabaseBackupRegressions_Call{Call: _e.mock.On("DatabaseBackupRegressions", local, instanceConfig, outputFile, backupTo)}
}

func (_c *Application_DatabaseBackupRegressions_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, outputFile string, backupTo string)) *Application_DatabaseBackupRegressions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *Application_DatabaseBackupRegressions_Call) Return(_a0 error) *Application_DatabaseBackupRegressions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseBackupRegressions_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string, string) error) *Application_DatabaseBackupRegressions_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseBackupShortcuts provides a mock function with given fields: local, instanceConfig, outputFile
func (_m *Application) DatabaseBackupShortcuts(local bool, instanceConfig *config.InstanceConfig, outputFile string) error {
	ret := _m.Called(local, instanceConfig, outputFile)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseBackupShortcuts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string) error); ok {
		r0 = rf(local, instanceConfig, outputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseBackupShortcuts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseBackupShortcuts'
type Application_DatabaseBackupShortcuts_Call struct {
	*mock.Call
}

// DatabaseBackupShortcuts is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - outputFile string
func (_e *Application_Expecter) DatabaseBackupShortcuts(local interface{}, instanceConfig interface{}, outputFile interface{}) *Application_DatabaseBackupShortcuts_Call {
	return &Application_DatabaseBackupShortcuts_Call{Call: _e.mock.On("DatabaseBackupShortcuts", local, instanceConfig, outputFile)}
}

func (_c *Application_DatabaseBackupShortcuts_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, outputFile string)) *Application_DatabaseBackupShortcuts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string))
	})
	return _c
}

func (_c *Application_DatabaseBackupShortcuts_Call) Return(_a0 error) *Application_DatabaseBackupShortcuts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseBackupShortcuts_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string) error) *Application_DatabaseBackupShortcuts_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseRestoreAlerts provides a mock function with given fields: local, instanceConfig, inputFile
func (_m *Application) DatabaseRestoreAlerts(local bool, instanceConfig *config.InstanceConfig, inputFile string) error {
	ret := _m.Called(local, instanceConfig, inputFile)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseRestoreAlerts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string) error); ok {
		r0 = rf(local, instanceConfig, inputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseRestoreAlerts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseRestoreAlerts'
type Application_DatabaseRestoreAlerts_Call struct {
	*mock.Call
}

// DatabaseRestoreAlerts is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - inputFile string
func (_e *Application_Expecter) DatabaseRestoreAlerts(local interface{}, instanceConfig interface{}, inputFile interface{}) *Application_DatabaseRestoreAlerts_Call {
	return &Application_DatabaseRestoreAlerts_Call{Call: _e.mock.On("DatabaseRestoreAlerts", local, instanceConfig, inputFile)}
}

func (_c *Application_DatabaseRestoreAlerts_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, inputFile string)) *Application_DatabaseRestoreAlerts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string))
	})
	return _c
}

func (_c *Application_DatabaseRestoreAlerts_Call) Return(_a0 error) *Application_DatabaseRestoreAlerts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseRestoreAlerts_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string) error) *Application_DatabaseRestoreAlerts_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseRestoreRegressions provides a mock function with given fields: local, instanceConfig, inputFile
func (_m *Application) DatabaseRestoreRegressions(local bool, instanceConfig *config.InstanceConfig, inputFile string) error {
	ret := _m.Called(local, instanceConfig, inputFile)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseRestoreRegressions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string) error); ok {
		r0 = rf(local, instanceConfig, inputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseRestoreRegressions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseRestoreRegressions'
type Application_DatabaseRestoreRegressions_Call struct {
	*mock.Call
}

// DatabaseRestoreRegressions is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - inputFile string
func (_e *Application_Expecter) DatabaseRestoreRegressions(local interface{}, instanceConfig interface{}, inputFile interface{}) *Application_DatabaseRestoreRegressions_Call {
	return &Application_DatabaseRestoreRegressions_Call{Call: _e.mock.On("DatabaseRestoreRegressions", local, instanceConfig, inputFile)}
}

func (_c *Application_DatabaseRestoreRegressions_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, inputFile string)) *Application_DatabaseRestoreRegressions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string))
	})
	return _c
}

func (_c *Application_DatabaseRestoreRegressions_Call) Return(_a0 error) *Application_DatabaseRestoreRegressions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseRestoreRegressions_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string) error) *Application_DatabaseRestoreRegressions_Call {
	_c.Call.Return(run)
	return _c
}

// DatabaseRestoreShortcuts provides a mock function with given fields: local, instanceConfig, inputFile
func (_m *Application) DatabaseRestoreShortcuts(local bool, instanceConfig *config.InstanceConfig, inputFile string) error {
	ret := _m.Called(local, instanceConfig, inputFile)

	if len(ret) == 0 {
		panic("no return value specified for DatabaseRestoreShortcuts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string) error); ok {
		r0 = rf(local, instanceConfig, inputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_DatabaseRestoreShortcuts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DatabaseRestoreShortcuts'
type Application_DatabaseRestoreShortcuts_Call struct {
	*mock.Call
}

// DatabaseRestoreShortcuts is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - inputFile string
func (_e *Application_Expecter) DatabaseRestoreShortcuts(local interface{}, instanceConfig interface{}, inputFile interface{}) *Application_DatabaseRestoreShortcuts_Call {
	return &Application_DatabaseRestoreShortcuts_Call{Call: _e.mock.On("DatabaseRestoreShortcuts", local, instanceConfig, inputFile)}
}

func (_c *Application_DatabaseRestoreShortcuts_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, inputFile string)) *Application_DatabaseRestoreShortcuts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string))
	})
	return _c
}

func (_c *Application_DatabaseRestoreShortcuts_Call) Return(_a0 error) *Application_DatabaseRestoreShortcuts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_DatabaseRestoreShortcuts_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string) error) *Application_DatabaseRestoreShortcuts_Call {
	_c.Call.Return(run)
	return _c
}

// IngestForceReingest provides a mock function with given fields: local, instanceConfig, start, stop, dryrun
func (_m *Application) IngestForceReingest(local bool, instanceConfig *config.InstanceConfig, start string, stop string, dryrun bool) error {
	ret := _m.Called(local, instanceConfig, start, stop, dryrun)

	if len(ret) == 0 {
		panic("no return value specified for IngestForceReingest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, *config.InstanceConfig, string, string, bool) error); ok {
		r0 = rf(local, instanceConfig, start, stop, dryrun)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_IngestForceReingest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IngestForceReingest'
type Application_IngestForceReingest_Call struct {
	*mock.Call
}

// IngestForceReingest is a helper method to define mock.On call
//   - local bool
//   - instanceConfig *config.InstanceConfig
//   - start string
//   - stop string
//   - dryrun bool
func (_e *Application_Expecter) IngestForceReingest(local interface{}, instanceConfig interface{}, start interface{}, stop interface{}, dryrun interface{}) *Application_IngestForceReingest_Call {
	return &Application_IngestForceReingest_Call{Call: _e.mock.On("IngestForceReingest", local, instanceConfig, start, stop, dryrun)}
}

func (_c *Application_IngestForceReingest_Call) Run(run func(local bool, instanceConfig *config.InstanceConfig, start string, stop string, dryrun bool)) *Application_IngestForceReingest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(*config.InstanceConfig), args[2].(string), args[3].(string), args[4].(bool))
	})
	return _c
}

func (_c *Application_IngestForceReingest_Call) Return(_a0 error) *Application_IngestForceReingest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_IngestForceReingest_Call) RunAndReturn(run func(bool, *config.InstanceConfig, string, string, bool) error) *Application_IngestForceReingest_Call {
	_c.Call.Return(run)
	return _c
}

// IngestValidate provides a mock function with given fields: inputFile, verbose
func (_m *Application) IngestValidate(inputFile string, verbose bool) error {
	ret := _m.Called(inputFile, verbose)

	if len(ret) == 0 {
		panic("no return value specified for IngestValidate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool) error); ok {
		r0 = rf(inputFile, verbose)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_IngestValidate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IngestValidate'
type Application_IngestValidate_Call struct {
	*mock.Call
}

// IngestValidate is a helper method to define mock.On call
//   - inputFile string
//   - verbose bool
func (_e *Application_Expecter) IngestValidate(inputFile interface{}, verbose interface{}) *Application_IngestValidate_Call {
	return &Application_IngestValidate_Call{Call: _e.mock.On("IngestValidate", inputFile, verbose)}
}

func (_c *Application_IngestValidate_Call) Run(run func(inputFile string, verbose bool)) *Application_IngestValidate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool))
	})
	return _c
}

func (_c *Application_IngestValidate_Call) Return(_a0 error) *Application_IngestValidate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_IngestValidate_Call) RunAndReturn(run func(string, bool) error) *Application_IngestValidate_Call {
	_c.Call.Return(run)
	return _c
}

// TilesLast provides a mock function with given fields: store
func (_m *Application) TilesLast(store tracestore.TraceStore) error {
	ret := _m.Called(store)

	if len(ret) == 0 {
		panic("no return value specified for TilesLast")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(tracestore.TraceStore) error); ok {
		r0 = rf(store)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_TilesLast_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TilesLast'
type Application_TilesLast_Call struct {
	*mock.Call
}

// TilesLast is a helper method to define mock.On call
//   - store tracestore.TraceStore
func (_e *Application_Expecter) TilesLast(store interface{}) *Application_TilesLast_Call {
	return &Application_TilesLast_Call{Call: _e.mock.On("TilesLast", store)}
}

func (_c *Application_TilesLast_Call) Run(run func(store tracestore.TraceStore)) *Application_TilesLast_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tracestore.TraceStore))
	})
	return _c
}

func (_c *Application_TilesLast_Call) Return(_a0 error) *Application_TilesLast_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_TilesLast_Call) RunAndReturn(run func(tracestore.TraceStore) error) *Application_TilesLast_Call {
	_c.Call.Return(run)
	return _c
}

// TilesList provides a mock function with given fields: store, num
func (_m *Application) TilesList(store tracestore.TraceStore, num int) error {
	ret := _m.Called(store, num)

	if len(ret) == 0 {
		panic("no return value specified for TilesList")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(tracestore.TraceStore, int) error); ok {
		r0 = rf(store, num)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_TilesList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TilesList'
type Application_TilesList_Call struct {
	*mock.Call
}

// TilesList is a helper method to define mock.On call
//   - store tracestore.TraceStore
//   - num int
func (_e *Application_Expecter) TilesList(store interface{}, num interface{}) *Application_TilesList_Call {
	return &Application_TilesList_Call{Call: _e.mock.On("TilesList", store, num)}
}

func (_c *Application_TilesList_Call) Run(run func(store tracestore.TraceStore, num int)) *Application_TilesList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tracestore.TraceStore), args[1].(int))
	})
	return _c
}

func (_c *Application_TilesList_Call) Return(_a0 error) *Application_TilesList_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_TilesList_Call) RunAndReturn(run func(tracestore.TraceStore, int) error) *Application_TilesList_Call {
	_c.Call.Return(run)
	return _c
}

// TracesExport provides a mock function with given fields: store, queryString, begin, end, outputFile
func (_m *Application) TracesExport(store tracestore.TraceStore, queryString string, begin types.CommitNumber, end types.CommitNumber, outputFile string) error {
	ret := _m.Called(store, queryString, begin, end, outputFile)

	if len(ret) == 0 {
		panic("no return value specified for TracesExport")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(tracestore.TraceStore, string, types.CommitNumber, types.CommitNumber, string) error); ok {
		r0 = rf(store, queryString, begin, end, outputFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_TracesExport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TracesExport'
type Application_TracesExport_Call struct {
	*mock.Call
}

// TracesExport is a helper method to define mock.On call
//   - store tracestore.TraceStore
//   - queryString string
//   - begin types.CommitNumber
//   - end types.CommitNumber
//   - outputFile string
func (_e *Application_Expecter) TracesExport(store interface{}, queryString interface{}, begin interface{}, end interface{}, outputFile interface{}) *Application_TracesExport_Call {
	return &Application_TracesExport_Call{Call: _e.mock.On("TracesExport", store, queryString, begin, end, outputFile)}
}

func (_c *Application_TracesExport_Call) Run(run func(store tracestore.TraceStore, queryString string, begin types.CommitNumber, end types.CommitNumber, outputFile string)) *Application_TracesExport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tracestore.TraceStore), args[1].(string), args[2].(types.CommitNumber), args[3].(types.CommitNumber), args[4].(string))
	})
	return _c
}

func (_c *Application_TracesExport_Call) Return(_a0 error) *Application_TracesExport_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_TracesExport_Call) RunAndReturn(run func(tracestore.TraceStore, string, types.CommitNumber, types.CommitNumber, string) error) *Application_TracesExport_Call {
	_c.Call.Return(run)
	return _c
}

// TracesList provides a mock function with given fields: store, queryString, tileNumber
func (_m *Application) TracesList(store tracestore.TraceStore, queryString string, tileNumber types.TileNumber) error {
	ret := _m.Called(store, queryString, tileNumber)

	if len(ret) == 0 {
		panic("no return value specified for TracesList")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(tracestore.TraceStore, string, types.TileNumber) error); ok {
		r0 = rf(store, queryString, tileNumber)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_TracesList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TracesList'
type Application_TracesList_Call struct {
	*mock.Call
}

// TracesList is a helper method to define mock.On call
//   - store tracestore.TraceStore
//   - queryString string
//   - tileNumber types.TileNumber
func (_e *Application_Expecter) TracesList(store interface{}, queryString interface{}, tileNumber interface{}) *Application_TracesList_Call {
	return &Application_TracesList_Call{Call: _e.mock.On("TracesList", store, queryString, tileNumber)}
}

func (_c *Application_TracesList_Call) Run(run func(store tracestore.TraceStore, queryString string, tileNumber types.TileNumber)) *Application_TracesList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tracestore.TraceStore), args[1].(string), args[2].(types.TileNumber))
	})
	return _c
}

func (_c *Application_TracesList_Call) Return(_a0 error) *Application_TracesList_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_TracesList_Call) RunAndReturn(run func(tracestore.TraceStore, string, types.TileNumber) error) *Application_TracesList_Call {
	_c.Call.Return(run)
	return _c
}

// TrybotReference provides a mock function with given fields: local, store, instanceConfig, trybotFilename, outputFilename, numCommits
func (_m *Application) TrybotReference(local bool, store tracestore.TraceStore, instanceConfig *config.InstanceConfig, trybotFilename string, outputFilename string, numCommits int) error {
	ret := _m.Called(local, store, instanceConfig, trybotFilename, outputFilename, numCommits)

	if len(ret) == 0 {
		panic("no return value specified for TrybotReference")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, tracestore.TraceStore, *config.InstanceConfig, string, string, int) error); ok {
		r0 = rf(local, store, instanceConfig, trybotFilename, outputFilename, numCommits)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Application_TrybotReference_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TrybotReference'
type Application_TrybotReference_Call struct {
	*mock.Call
}

// TrybotReference is a helper method to define mock.On call
//   - local bool
//   - store tracestore.TraceStore
//   - instanceConfig *config.InstanceConfig
//   - trybotFilename string
//   - outputFilename string
//   - numCommits int
func (_e *Application_Expecter) TrybotReference(local interface{}, store interface{}, instanceConfig interface{}, trybotFilename interface{}, outputFilename interface{}, numCommits interface{}) *Application_TrybotReference_Call {
	return &Application_TrybotReference_Call{Call: _e.mock.On("TrybotReference", local, store, instanceConfig, trybotFilename, outputFilename, numCommits)}
}

func (_c *Application_TrybotReference_Call) Run(run func(local bool, store tracestore.TraceStore, instanceConfig *config.InstanceConfig, trybotFilename string, outputFilename string, numCommits int)) *Application_TrybotReference_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(tracestore.TraceStore), args[2].(*config.InstanceConfig), args[3].(string), args[4].(string), args[5].(int))
	})
	return _c
}

func (_c *Application_TrybotReference_Call) Return(_a0 error) *Application_TrybotReference_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Application_TrybotReference_Call) RunAndReturn(run func(bool, tracestore.TraceStore, *config.InstanceConfig, string, string, int) error) *Application_TrybotReference_Call {
	_c.Call.Return(run)
	return _c
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
