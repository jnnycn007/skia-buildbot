/*
	The Cluster Telemetry Frontend.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/skia-dev/glog"

	"go.skia.org/infra/ct/go/db"
	"go.skia.org/infra/ct/go/util"
	"go.skia.org/infra/go/common"
	"go.skia.org/infra/go/influxdb"
	"go.skia.org/infra/go/login"
	"go.skia.org/infra/go/metadata"
	"go.skia.org/infra/go/skiaversion"
	skutil "go.skia.org/infra/go/util"
)

const (
	// Default page size used for pagination.
	DEFAULT_PAGE_SIZE = 5

	// Maximum page size used for pagination.
	MAX_PAGE_SIZE = 100
)

var (
	taskTables = []string{
		db.TABLE_CHROMIUM_PERF_TASKS,
		db.TABLE_RECREATE_PAGE_SETS_TASKS,
		db.TABLE_RECREATE_WEBPAGE_ARCHIVES_TASKS,
	}

	chromiumPerfTemplate                       *template.Template = nil
	chromiumPerfRunsHistoryTemplate            *template.Template = nil
	adminTasksTemplate                         *template.Template = nil
	recreatePageSetsRunsHistoryTemplate        *template.Template = nil
	recreateWebpageArchivesRunsHistoryTemplate *template.Template = nil
	runsHistoryTemplate                        *template.Template = nil
	pendingTasksTemplate                       *template.Template = nil

	dbClient *influxdb.Client = nil
)

// flags
var (
	graphiteServer = flag.String("graphite_server", "localhost:2003", "Where is Graphite metrics ingestion server running.")
	host           = flag.String("host", "localhost", "HTTP service host")
	port           = flag.String("port", ":8002", "HTTP service port (e.g., ':8002')")
	local          = flag.Bool("local", false, "Running locally if true. As opposed to in production.")
	workdir        = flag.String("workdir", ".", "Directory to use for scratch work.")
	resourcesDir   = flag.String("resources_dir", "", "The directory to find templates, JS, and CSS files. If blank the current directory will be used.")
)

type CommonCols struct {
	Id          int64         `db:"id"`
	TsAdded     sql.NullInt64 `db:"ts_added"`
	TsStarted   sql.NullInt64 `db:"ts_started"`
	TsCompleted sql.NullInt64 `db:"ts_completed"`
	Username    string        `db:"username"`
	Failure     sql.NullBool  `db:"failure"`
}

type Task interface {
	GetAddedTimestamp() int64
	GetTaskName() string
	TableName() string
	// Returns a slice of the struct type.
	Select(query string, args ...interface{}) (interface{}, error)
}

func (dbrow *CommonCols) GetAddedTimestamp() int64 {
	return dbrow.TsAdded.Int64
}

type ChromiumPerfDBTask struct {
	CommonCols

	Benchmark            string         `db:"benchmark"`
	Platform             string         `db:"platform"`
	PageSets             string         `db:"page_sets"`
	RepeatRuns           int64          `db:"repeat_runs"`
	BenchmarkArgs        string         `db:"benchmark_args"`
	BrowserArgsNoPatch   string         `db:"browser_args_nopatch"`
	BrowserArgsWithPatch string         `db:"browser_args_withpatch"`
	Description          string         `db:"description"`
	ChromiumPatch        string         `db:"chromium_patch"`
	BlinkPatch           string         `db:"blink_patch"`
	SkiaPatch            string         `db:"skia_patch"`
	Results              sql.NullString `db:"results"`
	NoPatchRawOutput     sql.NullString `db:"nopatch_raw_output"`
	WithPatchRawOutput   sql.NullString `db:"withpatch_raw_output"`
}

func (task ChromiumPerfDBTask) GetTaskName() string {
	return "ChromiumPerf"
}

func (task ChromiumPerfDBTask) TableName() string {
	return db.TABLE_CHROMIUM_PERF_TASKS
}

func (task ChromiumPerfDBTask) Select(query string, args ...interface{}) (interface{}, error) {
	result := []ChromiumPerfDBTask{}
	err := db.DB.Select(&result, query, args...)
	return result, err
}

type RecreatePageSetsDBTask struct {
	CommonCols

	PageSets string `db:"page_sets"`
}

func (task RecreatePageSetsDBTask) GetTaskName() string {
	return "RecreatePageSets"
}

func (task RecreatePageSetsDBTask) TableName() string {
	return db.TABLE_RECREATE_PAGE_SETS_TASKS
}

func (task RecreatePageSetsDBTask) Select(query string, args ...interface{}) (interface{}, error) {
	result := []RecreatePageSetsDBTask{}
	err := db.DB.Select(&result, query, args...)
	return result, err
}

type RecreateWebpageArchivesDBTask struct {
	CommonCols

	PageSets      string `db:"page_sets"`
	ChromiumBuild string `db:"chromium_build"`
}

func (task RecreateWebpageArchivesDBTask) GetTaskName() string {
	return "RecreateWebpageArchives"
}

func (task RecreateWebpageArchivesDBTask) TableName() string {
	return db.TABLE_RECREATE_WEBPAGE_ARCHIVES_TASKS
}

func (task RecreateWebpageArchivesDBTask) Select(query string, args ...interface{}) (interface{}, error) {
	result := []RecreateWebpageArchivesDBTask{}
	err := db.DB.Select(&result, query, args...)
	return result, err
}

func reloadTemplates() {
	if *resourcesDir == "" {
		// If resourcesDir is not specified then consider the directory two directories up from this
		// source file as the resourcesDir.
		_, filename, _, _ := runtime.Caller(0)
		*resourcesDir = filepath.Join(filepath.Dir(filename), "../..")
	}
	chromiumPerfTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/chromium_perf.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))
	chromiumPerfRunsHistoryTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/chromium_perf_runs_history.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))

	adminTasksTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/admin_tasks.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))
	recreatePageSetsRunsHistoryTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/recreate_page_sets_runs_history.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))
	recreateWebpageArchivesRunsHistoryTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/recreate_webpage_archives_runs_history.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))

	runsHistoryTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/runs_history.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))

	pendingTasksTemplate = template.Must(template.ParseFiles(
		filepath.Join(*resourcesDir, "templates/pending_tasks.html"),
		filepath.Join(*resourcesDir, "templates/header.html"),
		filepath.Join(*resourcesDir, "templates/titlebar.html"),
		filepath.Join(*resourcesDir, "templates/drawer.html"),
	))
}

func Init() {
	reloadTemplates()
}

func userHasEditRights(r *http.Request) bool {
	return strings.HasSuffix(login.LoggedInAs(r), "@google.com") || strings.HasSuffix(login.LoggedInAs(r), "@chromium.org")
}

func userHasAdminRights(r *http.Request) bool {
	// TODO(benjaminwagner): Add this list to GCE project level metadata and retrieve from there.
	admins := map[string]bool{
		"benjaminwagner@google.com": true,
		"borenet@google.com":        true,
		"jcgregorio@google.com":     true,
		"rmistry@google.com":        true,
		"stephana@google.com":       true,
	}
	return userHasEditRights(r) && admins[login.LoggedInAs(r)]
}

func getCurrentTs() string {
	return time.Now().UTC().Format("20060102150405")
}

func getIntParam(name string, r *http.Request) (*int, error) {
	raw, ok := r.URL.Query()[name]
	if !ok {
		return nil, nil
	}
	v64, err := strconv.ParseInt(raw[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for parameter %q: %s -- %v", name, raw, err)
	}
	v32 := int(v64)
	return &v32, nil
}

func executeSimpleTemplate(template *template.Template, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Don't use cached templates in local mode.
	if *local {
		reloadTemplates()
	}

	if err := template.Execute(w, struct{}{}); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to expand template: %v", err))
		return
	}
}

// Data included in all tasks; set by addTaskHandler.
type AddTaskCommonVars struct {
	Username string
	TsAdded  string
}

type AddTaskVars interface {
	GetAddTaskCommonVars() *AddTaskCommonVars
	IsAdminTask() bool
	GetInsertQueryAndBinds() (string, []interface{})
}

func (vars *AddTaskCommonVars) GetAddTaskCommonVars() *AddTaskCommonVars {
	return vars
}

func (vars *AddTaskCommonVars) IsAdminTask() bool {
	return false
}

func addTaskHandler(w http.ResponseWriter, r *http.Request, task AddTaskVars) {
	if !userHasEditRights(r) {
		skutil.ReportError(w, r, fmt.Errorf("Must have google or chromium account to add tasks"), "")
		return
	}
	if task.IsAdminTask() && !userHasAdminRights(r) {
		skutil.ReportError(w, r, fmt.Errorf("Must be admin to add admin tasks; contact rmistry@"), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to add %v task: %v", task, err))
		return
	}
	defer skutil.Close(r.Body)

	task.GetAddTaskCommonVars().Username = login.LoggedInAs(r)
	task.GetAddTaskCommonVars().TsAdded = getCurrentTs()

	query, binds := task.GetInsertQueryAndBinds()
	_, err := db.DB.Exec(query, binds...)
	if err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to insert %v task: %v", task, err))
		return
	}
}

func chromiumPerfView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(chromiumPerfTemplate, w, r)
}

func chromiumPerfHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := map[string]interface{}{
		"benchmarks": util.SupportedBenchmarks,
		"platforms":  util.SupportedPlatformsToDesc,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to encode JSON: %v", err))
		return
	}
}

func pageSetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageSetsToDesc := map[string]string{}
	for pageSet := range util.PagesetTypeToInfo {
		pageSetsToDesc[pageSet] = util.PagesetTypeToInfo[pageSet].Description
	}

	if err := json.NewEncoder(w).Encode(pageSetsToDesc); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to encode JSON: %v", err))
		return
	}
}

// ChromiumPerfVars is the type used by the Chromium Perf pages.
type ChromiumPerfVars struct {
	AddTaskCommonVars

	Benchmark            string `json:"benchmark"`
	Platform             string `json:"platform"`
	PageSets             string `json:"page_sets"`
	RepeatRuns           string `json:"repeat_runs"`
	BenchmarkArgs        string `json:"benchmark_args"`
	BrowserArgsNoPatch   string `json:"browser_args_nopatch"`
	BrowserArgsWithPatch string `json:"browser_args_withpatch"`
	Description          string `json:"desc"`
	ChromiumPatch        string `json:"chromium_patch"`
	BlinkPatch           string `json:"blink_patch"`
	SkiaPatch            string `json:"skia_patch"`
}

func (task *ChromiumPerfVars) GetInsertQueryAndBinds() (string, []interface{}) {
	return fmt.Sprintf("INSERT INTO %s (username,benchmark,platform,page_sets,repeat_runs,benchmark_args,browser_args_nopatch,browser_args_withpatch,description,chromium_patch,blink_patch,skia_patch,ts_added) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);",
			db.TABLE_CHROMIUM_PERF_TASKS),
		[]interface{}{
			task.Username,
			task.Benchmark,
			task.Platform,
			task.PageSets,
			task.RepeatRuns,
			task.BenchmarkArgs,
			task.BrowserArgsNoPatch,
			task.BrowserArgsWithPatch,
			task.Description,
			task.ChromiumPatch,
			task.BlinkPatch,
			task.SkiaPatch,
			task.TsAdded,
		}
}

func addChromiumPerfTaskHandler(w http.ResponseWriter, r *http.Request) {
	addTaskHandler(w, r, &ChromiumPerfVars{})
}

func dbTaskQuery(prototype Task, username string, includeCompleted bool, countQuery bool, offset int, size int) (string, []interface{}) {
	args := []interface{}{}
	query := "SELECT "
	if countQuery {
		query += "COUNT(*)"
	} else {
		query += "*"
	}
	query += fmt.Sprintf(" FROM %s", prototype.TableName())
	if username != "" {
		query += " WHERE username=?"
		args = append(args, username)
	} else if !includeCompleted {
		query += " WHERE ts_completed IS NULL"
	}
	if !countQuery {
		query += " ORDER BY id DESC LIMIT ?,?"
		args = append(args, offset, size)
	}
	return query, args
}

func getTasksHandler(prototype Task, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Filter by either username or not started yet.
	username := r.FormValue("username")
	notCompleted := r.FormValue("not_completed")
	offset, size, err := skutil.PaginationParams(r.URL.Query(), 0, DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
	if err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to get pagination params: %v", err))
		return
	}
	query, args := dbTaskQuery(prototype, username, notCompleted == "", false, offset, size)
	glog.Infof("Running %s", query)
	data, err := prototype.Select(query, args...)
	if err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to query %s tasks: %v", prototype.GetTaskName(), err))
		return
	}

	query, args = dbTaskQuery(prototype, username, notCompleted == "", true, 0, 0)
	// Get the total count.
	glog.Infof("Running %s", query)
	countVal := []int{}
	if err := db.DB.Select(&countVal, query, args...); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to query %s tasks: %v", prototype.GetTaskName(), err))
		return
	}

	pagination := &skutil.ResponsePagination{
		Offset: offset,
		Size:   size,
		Total:  countVal[0],
	}
	jsonResponse := map[string]interface{}{
		"data":       data,
		"pagination": pagination,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to encode JSON: %v", err))
		return
	}
}

func getChromiumPerfTasksHandler(w http.ResponseWriter, r *http.Request) {
	getTasksHandler(&ChromiumPerfDBTask{}, w, r)
}

func chromiumPerfRunsHistoryView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(chromiumPerfRunsHistoryTemplate, w, r)
}

func adminTasksView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(adminTasksTemplate, w, r)
}

type AdminTaskVars struct {
	AddTaskCommonVars
}

func (vars *AdminTaskVars) IsAdminTask() bool {
	return true
}

// Represents the parameters sent as JSON to the add_recreate_page_sets_task handler.
type RecreatePageSetsTaskHandlerVars struct {
	AdminTaskVars
	PageSets string `json:"page_sets"`
}

func (task *RecreatePageSetsTaskHandlerVars) GetInsertQueryAndBinds() (string, []interface{}) {
	return fmt.Sprintf("INSERT INTO %s (username,page_sets,ts_added) VALUES (?,?,?);",
			db.TABLE_RECREATE_PAGE_SETS_TASKS),
		[]interface{}{
			task.Username,
			task.PageSets,
			task.TsAdded,
		}
}

func addRecreatePageSetsTaskHandler(w http.ResponseWriter, r *http.Request) {
	addTaskHandler(w, r, &RecreatePageSetsTaskHandlerVars{})
}

// Represents the parameters sent as JSON to the add_recreate_webpage_archives_task handler.
type RecreateWebpageArchivesTaskHandlerVars struct {
	AdminTaskVars
	PageSets      string `json:"page_sets"`
	ChromiumBuild string `json:"chromium_build"`
}

func (task *RecreateWebpageArchivesTaskHandlerVars) GetInsertQueryAndBinds() (string, []interface{}) {
	return fmt.Sprintf("INSERT INTO %s (username,page_sets,chromium_build,ts_added) VALUES (?,?,?,?);",
			db.TABLE_RECREATE_WEBPAGE_ARCHIVES_TASKS),
		[]interface{}{
			task.Username,
			task.PageSets,
			task.ChromiumBuild,
			task.TsAdded,
		}
}

func addRecreateWebpageArchivesTaskHandler(w http.ResponseWriter, r *http.Request) {
	addTaskHandler(w, r, &RecreateWebpageArchivesTaskHandlerVars{})
}

func recreatePageSetsRunsHistoryView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(recreatePageSetsRunsHistoryTemplate, w, r)
}

func recreateWebpageArchivesRunsHistoryView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(recreateWebpageArchivesRunsHistoryTemplate, w, r)
}

func getRecreatePageSetsTasksHandler(w http.ResponseWriter, r *http.Request) {
	getTasksHandler(&RecreatePageSetsDBTask{}, w, r)
}

func getRecreateWebpageArchivesTasksHandler(w http.ResponseWriter, r *http.Request) {
	getTasksHandler(&RecreateWebpageArchivesDBTask{}, w, r)
}

func chromiumBuildsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO(benjaminwagner): load real data
	resp := []struct {
		Key         string `json:"key"`
		Description string `json:"description"`
	}{
		{Key: "abc", Description: "first build"},
		{Key: "def", Description: "second build"},
		{Key: "ghi", Description: "third build"},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to encode JSON: %v", err))
		return
	}
}

func runsHistoryView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(runsHistoryTemplate, w, r)
}

func getAllPendingTasks() ([]Task, error) {
	tasks := []Task{}
	for _, tableName := range taskTables {
		var task Task
		query := fmt.Sprintf("SELECT * FROM %s WHERE ts_completed IS NULL ORDER BY ts_added LIMIT 1;", tableName)
		switch tableName {
		case db.TABLE_CHROMIUM_PERF_TASKS:
			task = &ChromiumPerfDBTask{}
		case db.TABLE_RECREATE_PAGE_SETS_TASKS:
			task = &RecreatePageSetsDBTask{}
		case db.TABLE_RECREATE_WEBPAGE_ARCHIVES_TASKS:
			task = &RecreateWebpageArchivesDBTask{}
		default:
			panic("Unknown table " + tableName)
		}

		if err := db.DB.Get(task, query); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to query DB: %v", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func getOldestPendingTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := getAllPendingTasks()
	if err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to get all pending tasks: %v", err))
		return
	}

	var oldestTask Task
	for _, task := range tasks {
		if oldestTask == nil {
			oldestTask = task
		} else if oldestTask.GetAddedTimestamp() < task.GetAddedTimestamp() {
			oldestTask = task
		}
	}

	oldestTaskJsonRepr := map[string]Task{}
	if oldestTask != nil {
		oldestTaskJsonRepr[oldestTask.GetTaskName()] = oldestTask
	}
	if err := json.NewEncoder(w).Encode(oldestTaskJsonRepr); err != nil {
		skutil.ReportError(w, r, err, fmt.Sprintf("Failed to encode JSON: %v", err))
		return
	}
}

func pendingTasksView(w http.ResponseWriter, r *http.Request) {
	executeSimpleTemplate(pendingTasksTemplate, w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, login.LoginURL(w, r), http.StatusFound)
	return
}

func runServer(serverURL string) {
	r := mux.NewRouter()
	r.PathPrefix("/res/").HandlerFunc(skutil.MakeResourceHandler(*resourcesDir))

	// Chromium Perf handlers.
	r.HandleFunc("/", chromiumPerfView).Methods("GET")
	r.HandleFunc("/chromium_perf/", chromiumPerfView).Methods("GET")
	r.HandleFunc("/chromium_perf_runs/", chromiumPerfRunsHistoryView).Methods("GET")
	r.HandleFunc("/_/chromium_perf/", chromiumPerfHandler).Methods("POST")
	r.HandleFunc("/_/add_chromium_perf_task", addChromiumPerfTaskHandler).Methods("POST")
	r.HandleFunc("/_/get_chromium_perf_tasks", getChromiumPerfTasksHandler).Methods("POST")

	// Admin Tasks handlers.
	r.HandleFunc("/admin_tasks/", adminTasksView).Methods("GET")
	r.HandleFunc("/recreate_page_sets_runs/", recreatePageSetsRunsHistoryView).Methods("GET")
	r.HandleFunc("/recreate_webpage_archives_runs/", recreateWebpageArchivesRunsHistoryView).Methods("GET")
	r.HandleFunc("/_/add_recreate_page_sets_task", addRecreatePageSetsTaskHandler).Methods("POST")
	r.HandleFunc("/_/add_recreate_webpage_archives_task", addRecreateWebpageArchivesTaskHandler).Methods("POST")
	r.HandleFunc("/_/get_recreate_page_sets_tasks", getRecreatePageSetsTasksHandler).Methods("POST")
	r.HandleFunc("/_/get_recreate_webpage_archives_tasks", getRecreateWebpageArchivesTasksHandler).Methods("POST")

	// Runs history handlers.
	r.HandleFunc("/history/", runsHistoryView).Methods("GET")

	// Task Queue handlers.
	r.HandleFunc("/queue/", pendingTasksView).Methods("GET")
	r.HandleFunc("/_/get_oldest_pending_task", getOldestPendingTaskHandler).Methods("GET")

	// Common handlers used by different pages.
	r.HandleFunc("/_/page_sets/", pageSetsHandler).Methods("POST")
	r.HandleFunc("/_/chromium_builds/", chromiumBuildsHandler).Methods("POST")
	r.HandleFunc("/json/version", skiaversion.JsonHandler)
	r.HandleFunc("/oauth2callback/", login.OAuth2CallbackHandler)
	r.HandleFunc("/login/", loginHandler)
	r.HandleFunc("/logout/", login.LogoutHandler)
	r.HandleFunc("/loginstatus/", login.StatusHandler)
	http.Handle("/", skutil.LoggingGzipRequestResponse(r))
	glog.Infof("Ready to serve on %s", serverURL)
	glog.Fatal(http.ListenAndServe(*port, nil))
}

func main() {
	// Setup flags.
	dbConf := db.DBConfigFromFlags()
	influxdb.SetupFlags()

	common.InitWithMetrics("ctfe", graphiteServer)
	v, err := skiaversion.GetVersion()
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("Version %s, built at %s", v.Commit, v.Date)

	Init()
	serverURL := "https://" + *host
	if *local {
		serverURL = "http://" + *host + *port
	}

	// Setup InfluxDB client.
	dbClient, err = influxdb.NewClientFromFlagsAndMetadata(*local)
	if err != nil {
		glog.Fatal(err)
	}

	// By default use a set of credentials setup for localhost access.
	var cookieSalt = "notverysecret"
	var clientID = "31977622648-1873k0c1e5edaka4adpv1ppvhr5id3qm.apps.googleusercontent.com"
	var clientSecret = "cw0IosPu4yjaG2KWmppj2guj"
	var redirectURL = serverURL + "/oauth2callback/"
	if !*local {
		cookieSalt = metadata.Must(metadata.ProjectGet(metadata.COOKIESALT))
		clientID = metadata.Must(metadata.ProjectGet(metadata.CLIENT_ID))
		clientSecret = metadata.Must(metadata.ProjectGet(metadata.CLIENT_SECRET))
	}
	login.Init(clientID, clientSecret, redirectURL, cookieSalt, login.DEFAULT_SCOPE, login.DEFAULT_DOMAIN_WHITELIST, *local)

	glog.Info("CloneOrUpdate complete")

	// Initialize the ctfe database.
	if !*local {
		if err := dbConf.GetPasswordFromMetadata(); err != nil {
			glog.Fatal(err)
		}
	}
	if err := dbConf.InitDB(); err != nil {
		glog.Fatal(err)
	}

	runServer(serverURL)
}
