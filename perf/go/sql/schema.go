package sql

// Generated by //go/sql/exporter/
// DO NOT EDIT

const Schema = `CREATE TABLE IF NOT EXISTS Alerts (
  id INT PRIMARY KEY DEFAULT unique_rowid(),
  alert TEXT,
  config_state INT DEFAULT 0,
  last_modified INT
);
CREATE TABLE IF NOT EXISTS AnomalyGroups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  creation_time TIMESTAMPTZ DEFAULT now(),
  anomaly_ids UUID ARRAY,
  group_meta_data JSONB,
  common_rev_start INT,
  common_rev_end INT,
  action TEXT,
  action_time TIMESTAMPTZ,
  bisection_id TEXT,
  reported_issue_id TEXT,
  culprit_ids UUID ARRAY,
  last_modified_time TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS Commits (
  commit_number INT PRIMARY KEY,
  git_hash TEXT UNIQUE NOT NULL,
  commit_time INT,
  author TEXT,
  subject TEXT
);
CREATE TABLE IF NOT EXISTS Culprits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  host STRING,
  project STRING,
  ref STRING,
  revision STRING,
  last_modified INT,
  anomaly_group_ids STRING ARRAY,
  issue_ids STRING ARRAY,
  UNIQUE INDEX by_revision (revision, host, project, ref)
);
CREATE TABLE IF NOT EXISTS GraphsShortcuts (
  id TEXT UNIQUE NOT NULL PRIMARY KEY,
  graphs TEXT
);
CREATE TABLE IF NOT EXISTS ParamSets (
  tile_number INT,
  param_key STRING,
  param_value STRING,
  PRIMARY KEY (tile_number, param_key, param_value),
  INDEX by_tile_number (tile_number DESC)
);
CREATE TABLE IF NOT EXISTS Postings (
  tile_number INT,
  key_value STRING NOT NULL,
  trace_id BYTES,
  PRIMARY KEY (tile_number, key_value, trace_id),
  INDEX by_trace_id (tile_number, trace_id, key_value),
  INDEX by_key_value (tile_number, key_value)
);
CREATE TABLE IF NOT EXISTS Regressions (
  commit_number INT,
  alert_id INT,
  regression TEXT,
  PRIMARY KEY (commit_number, alert_id)
);
CREATE TABLE IF NOT EXISTS Regressions2 (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  commit_number INT,
  prev_commit_number INT,
  alert_id INT,
  creation_time TIMESTAMPTZ DEFAULT now(),
  median_before REAL,
  median_after REAL,
  is_improvement BOOL,
  cluster_type TEXT,
  cluster_summary JSONB,
  frame JSONB,
  triage_status TEXT,
  triage_message TEXT,
  INDEX by_alert_id (alert_id),
  INDEX by_commit_alert (commit_number, alert_id)
);
CREATE TABLE IF NOT EXISTS Shortcuts (
  id TEXT UNIQUE NOT NULL PRIMARY KEY,
  trace_ids TEXT
);
CREATE TABLE IF NOT EXISTS SourceFiles (
  source_file_id INT PRIMARY KEY DEFAULT unique_rowid(),
  source_file STRING UNIQUE NOT NULL,
  INDEX by_source_file (source_file, source_file_id)
);
CREATE TABLE IF NOT EXISTS Subscriptions (
  name STRING UNIQUE NOT NULL,
  revision STRING NOT NULL,
  bug_labels STRING ARRAY,
  hotlists STRING ARRAY,
  bug_component STRING,
  bug_cc_emails STRING ARRAY,
  contact_email STRING,
  PRIMARY KEY(name, revision)
);
CREATE TABLE IF NOT EXISTS TraceValues (
  trace_id BYTES,
  commit_number INT,
  val REAL,
  source_file_id INT,
  PRIMARY KEY (trace_id, commit_number),
  INDEX by_source_file_id (source_file_id, trace_id)
);
`

var Alerts = []string{
	"id",
	"alert",
	"config_state",
	"last_modified",
}

var AnomalyGroups = []string{
	"id",
	"creation_time",
	"anomaly_ids",
	"group_meta_data",
	"common_rev_start",
	"common_rev_end",
	"action",
	"action_time",
	"bisection_id",
	"reported_issue_id",
	"culprit_ids",
	"last_modified_time",
}

var Commits = []string{
	"commit_number",
	"git_hash",
	"commit_time",
	"author",
	"subject",
}

var Culprits = []string{
	"id",
	"host",
	"project",
	"ref",
	"revision",
	"last_modified",
	"anomaly_group_ids",
	"issue_ids",
	"UNIQUE",
}

var GraphsShortcuts = []string{
	"id",
	"graphs",
}

var ParamSets = []string{
	"tile_number",
	"param_key",
	"param_value",
}

var Postings = []string{
	"tile_number",
	"key_value",
	"trace_id",
}

var Regressions = []string{
	"commit_number",
	"alert_id",
	"regression",
}

var Regressions2 = []string{
	"id",
	"commit_number",
	"prev_commit_number",
	"alert_id",
	"creation_time",
	"median_before",
	"median_after",
	"is_improvement",
	"cluster_type",
	"cluster_summary",
	"frame",
	"triage_status",
	"triage_message",
}

var Shortcuts = []string{
	"id",
	"trace_ids",
}

var SourceFiles = []string{
	"source_file_id",
	"source_file",
}

var Subscriptions = []string{
	"name",
	"revision",
	"bug_labels",
	"hotlists",
	"bug_component",
	"bug_cc_emails",
	"contact_email",
}

var TraceValues = []string{
	"trace_id",
	"commit_number",
	"val",
	"source_file_id",
}
