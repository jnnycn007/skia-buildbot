// DO NOT EDIT. This file is automatically generated.

export namespace pivot {
	export interface Request {
		group_by: string[] | null;
		operation: pivot.Operation;
		summary: pivot.Operation[] | null;
	}
}

export interface Go2TS {
	GenerateNominalTypes: boolean;
}

export interface Alert {
	id_as_string: string;
	display_name: string;
	query: string;
	alert: string;
	issue_tracker_component: SerializesToString;
	interesting: number;
	bug_uri_template: string;
	algo: ClusterAlgo;
	step: StepDetection;
	state: ConfigState;
	owner: string;
	step_up_only: boolean;
	direction: Direction;
	radius: number;
	k: number;
	group_by: string;
	sparse: boolean;
	minimum_num: number;
	category: string;
	action?: AlertAction;
	sub_name?: string;
	sub_revision?: string;
}

export interface AlertsStatus {
	alerts: number;
}

export interface RevisionInfo {
	master: string;
	bot: string;
	benchmark: string;
	start_revision: number;
	end_revision: number;
	start_time: number;
	end_time: number;
	test: string;
	is_improvement: boolean;
	bug_id: string;
	explore_url: string;
	query: string;
	anomaly_ids: string[] | null;
}

export interface ValuePercent {
	value: string;
	percent: number;
}

export interface StepFit {
	least_squares: number;
	turning_point: number;
	step_size: number;
	regression: number;
	status: StepFitStatus;
}

export interface ColumnHeader {
	offset: CommitNumber;
	timestamp: TimestampSeconds;
}

export interface ClusterSummary {
	centroid: number[] | null;
	shortcut: string;
	param_summaries2: ValuePercent[] | null;
	step_fit: StepFit | null;
	step_point: ColumnHeader | null;
	num: number;
	ts: string;
	notification_id?: string;
}

export interface FavoritesSectionLinkConfig {
	id?: string;
	text: string;
	href: string;
	description: string;
}

export interface FavoritesSectionConfig {
	name: string;
	links: FavoritesSectionLinkConfig[] | null;
}

export interface Favorites {
	sections: FavoritesSectionConfig[] | null;
}

export interface QueryCacheConfig {
	type: CacheType;
	level1_cache_key?: string;
	level1_cache_values?: string[] | null;
	level2_cache_key?: string;
	level2_cache_values?: string[] | null;
	enabled?: boolean;
}

export interface RedisConfig {
	project?: string;
	zone?: string;
	instance?: string;
	cache_expiration_minutes?: number;
}

export interface QueryConfig {
	include_params?: string[] | null;
	default_param_selections?: { [key: string]: string[] | null } | null;
	default_url_values?: { [key: string]: string } | null;
	cache_config?: QueryCacheConfig;
	redis_config?: RedisConfig;
}

export interface Commit {
	offset: CommitNumber;
	hash: string;
	ts: number;
	author: string;
	message: string;
	url: string;
	body: string;
}

export interface DataFrame {
	traceset: TraceSet;
	header: (ColumnHeader | null)[] | null;
	paramset: ReadOnlyParamSet;
	skip: number;
}

export interface Anomaly {
	id: number;
	test_path: string;
	bug_id: number;
	start_revision: number;
	end_revision: number;
	start_revision_hash?: string;
	end_revision_hash?: string;
	is_improvement: boolean;
	recovered: boolean;
	state: string;
	statistic: string;
	units: string;
	degrees_of_freedom: number;
	median_before_anomaly: number;
	median_after_anomaly: number;
	p_value: number;
	segment_size_after: number;
	segment_size_before: number;
	std_dev_before_anomaly: number;
	t_statistic: number;
	subscription_name: string;
	bug_component: string;
	bug_labels: string[] | null;
	bug_cc_emails: string[] | null;
	bisect_ids: string[] | null;
}

export interface FrameResponse {
	dataframe: DataFrame | null;
	skps: number[] | null;
	msg: string;
	display_mode: FrameResponseDisplayMode;
	anomalymap: AnomalyMap;
}

export interface TriageStatus {
	status: Status;
	message: string;
}

export interface Regression {
	low: ClusterSummary | null;
	high: ClusterSummary | null;
	frame: FrameResponse | null;
	low_status: TriageStatus;
	high_status: TriageStatus;
	id: string;
	commit_number: CommitNumber;
	prev_commit_number: CommitNumber;
	alert_id: number;
	creation_time: string;
	median_before: number;
	median_after: number;
	is_improvement: boolean;
	cluster_type: string;
}

export interface RegressionAtCommit {
	cid: Commit;
	regression: Regression | null;
}

export interface FrameRequest {
	begin: number;
	end: number;
	formulas?: string[] | null;
	queries?: string[] | null;
	keys?: string;
	tz: string;
	num_commits?: number;
	request_type?: RequestType;
	disable_filter_parent_traces?: boolean;
	pivot?: pivot.Request | null;
}

export interface AlertUpdateResponse {
	IDAsString: string;
}

export interface CIDHandlerResponse {
	commitSlice: Commit[] | null;
	logEntry: string;
}

export interface ClusterStartResponse {
	id: string;
}

export interface CommitDetailsRequest {
	cid: CommitNumber;
	traceid: string;
}

export interface CountHandlerRequest {
	q: string;
	begin: number;
	end: number;
}

export interface CountHandlerResponse {
	count: number;
	paramset: ReadOnlyParamSet;
}

export interface Subscription {
	name?: string;
	revision?: string;
	bug_labels?: string[] | null;
	hotlists?: string[] | null;
	bug_component?: string;
	bug_priority?: number;
	bug_severity?: number;
	bug_cc_emails?: string[] | null;
	contact_email?: string;
}

export interface GetAnomaliesResponse {
	subscription: Subscription | null;
	alerts: Alert[] | null;
	anomaly_list: Anomaly[] | null;
	anomaly_cursor: string;
	error: string;
}

export interface Timerange {
	begin: number;
	end: number;
}

export interface GetGroupReportResponse {
	anomaly_list: Anomaly[] | null;
	sid: string;
	selected_keys: string[] | null;
	error: string;
	timerange_map: { [key: number]: Timerange } | null;
}

export interface GetGraphsShortcutRequest {
	id: string;
}

export interface GetSheriffListResponse {
	sheriff_list: string[] | null;
	error: string;
}

export interface NextParamListHandlerRequest {
	q: string;
}

export interface NextParamListHandlerResponse {
	count: number;
	paramset: ReadOnlyParamSet;
}

export interface RangeRequest {
	offset: CommitNumber;
	begin: number;
	end: number;
}

export interface RegressionRangeRequest {
	begin: number;
	end: number;
	subset: Subset;
	alert_filter: string;
}

export interface RegressionRow {
	cid: Commit;
	columns: (Regression | null)[] | null;
}

export interface RegressionRangeResponse {
	header: (Alert | null)[] | null;
	table: (RegressionRow | null)[] | null;
	categories: string[] | null;
}

export interface ShiftRequest {
	begin: CommitNumber;
	end: CommitNumber;
}

export interface ShiftResponse {
	begin: number;
	end: number;
}

export interface SkPerfConfig {
	radius: number;
	key_order: string[] | null;
	num_shift: number;
	interesting: number;
	step_up_only: boolean;
	commit_range_url: string;
	demo: boolean;
	display_group_by: boolean;
	hide_list_of_commits_on_explore: boolean;
	notifications: NotifierTypes;
	fetch_chrome_perf_anomalies: boolean;
	feedback_url: string;
	chat_url: string;
	help_url_override: string;
	trace_format: TraceFormat;
	need_alert_action: boolean;
	bug_host_url: string;
	git_repo_url: string;
	keys_for_commit_range: string[] | null;
	keys_for_useful_links: string[] | null;
	skip_commit_detail_display: boolean;
	image_tag: string;
}

export interface TriageRequest {
	cid: CommitNumber;
	alert: Alert;
	triage: TriageStatus;
	cluster_type: string;
}

export interface TriageResponse {
	bug: string;
}

export interface TryBugRequest {
	bug_uri_template: string;
}

export interface TryBugResponse {
	url: string;
}

export interface SourceIssues {
	sourceIssueIds?: Int64s;
}

export interface AttachmentDataRef {
	resourceName?: string;
}

export interface Attachment {
	attachmentDataRef?: AttachmentDataRef | null;
	attachmentId?: string;
	contentType?: string;
	etag?: string;
	filename?: string;
	length?: string;
}

export interface Date {
	day?: number;
	month?: number;
	year?: number;
}

export interface RepeatedDate {
	values?: (Date | null)[] | null;
}

export interface RepeatedString {
	values?: string[] | null;
}

export interface RepeatedDouble {
	values?: number[] | null;
}

export interface CustomFieldValue {
	customFieldId?: string;
	dateValue?: Date | null;
	displayString?: string;
	enumValue?: string;
	numericValue?: number;
	repeatedDateValue?: RepeatedDate | null;
	repeatedEnumValue?: RepeatedString | null;
	repeatedNumericValue?: RepeatedDouble | null;
	repeatedTextValue?: RepeatedString | null;
	textValue?: string;
}

export interface CustomField {
	componentId?: string;
	customFieldId?: string;
	description?: string;
	enumValues?: string[] | null;
	name?: string;
	required?: boolean;
	shared?: boolean;
	type?: string;
	typedDefaultValue?: CustomFieldValue | null;
}

export interface User {
	emailAddress?: string;
	obfuscatedEmailAddress?: string;
}

export interface IssueComment {
	comment?: string;
	commentNumber?: number;
	formattingMode?: string;
	issueId?: string;
	lastEditor?: User | null;
	modifiedTime?: string;
	version?: number;
	HTTPStatusCode: number;
	Header: Header;
}

export interface IssueAccessLimit {
	accessLevel?: string;
}

export interface IssueState {
	accessLimit?: IssueAccessLimit | null;
	assignee?: User | null;
	blockedByIssueIds?: Int64s;
	blockingIssueIds?: Int64s;
	canonicalIssueId?: string;
	ccs?: (User | null)[] | null;
	componentId?: string;
	customFields?: (CustomFieldValue | null)[] | null;
	duplicateIssueIds?: Int64s;
	foundInVersions?: string[] | null;
	hotlistIds?: Int64s;
	inProd?: boolean;
	isArchived?: boolean;
	priority?: string;
	reporter?: User | null;
	severity?: string;
	status?: string;
	targetedToVersions?: string[] | null;
	title?: string;
	type?: string;
	verifiedInVersions?: string[] | null;
	verifier?: User | null;
}

export interface Hyperlink {
	href?: string;
}

export interface IssueReference {
	commentNumber?: number;
	issueId?: string;
}

export interface RelatedLink {
	hyperlink?: Hyperlink | null;
	issueReference?: IssueReference | null;
	sourceCommentNumber?: number;
}

export interface StatusUpdate {
	formattingMode?: string;
	issueId?: string;
	updateText?: string;
}

export interface IssueUserData {
	editableCommentNumbers?: number[] | null;
	hasStarred?: boolean;
	hasUpvoted?: boolean;
}

export interface FieldId {
	customFieldId?: string;
	standardField?: string;
}

export interface Issue {
	ancestors?: { [key: string]: SourceIssues } | null;
	attachments?: (Attachment | null)[] | null;
	childIssueCount?: number;
	createdTime?: string;
	customFields?: (CustomField | null)[] | null;
	description?: IssueComment | null;
	etag?: string;
	isArchived?: boolean;
	issueComment?: IssueComment | null;
	issueId?: string;
	issueState?: IssueState | null;
	lastModifier?: User | null;
	majorModifiedTime?: string;
	minorModifiedTime?: string;
	modifiedTime?: string;
	parentIssueIds?: Int64s;
	relatedLinks?: (RelatedLink | null)[] | null;
	resolvedTime?: string;
	statusUpdate?: StatusUpdate | null;
	trackerId?: string;
	userData?: IssueUserData | null;
	verifiedTime?: string;
	version?: number;
	visibleFields?: (FieldId | null)[] | null;
	voteCount?: string;
	HTTPStatusCode: number;
	Header: Header;
}

export interface ListIssuesResponse {
	issues?: (Issue | null)[] | null;
}

export interface GraphConfig {
	queries: string[] | null;
	formulas: string[] | null;
	keys: string;
}

export interface GraphsShortcut {
	graphs: GraphConfig[] | null;
}

export interface CreateBisectRequest {
	comparison_mode: string;
	start_git_hash: string;
	end_git_hash: string;
	configuration: string;
	benchmark: string;
	story: string;
	chart: string;
	statistic: string;
	comparison_magnitude: string;
	pin: string;
	project: string;
	bug_id: string;
	user: string;
	alert_ids: string;
}

export interface CreatePinpointResponse {
	jobId: string;
	jobUrl: string;
}

export interface FullSummary {
	summary: ClusterSummary;
	triage: TriageStatus;
	frame: FrameResponse;
}

export interface Domain {
	n: number;
	end: string;
	offset: number;
}

export interface RegressionDetectionRequest {
	alert: Alert | null;
	domain: Domain;
	step: number;
	total_queries: number;
}

export interface ClusterSummaries {
	Clusters: (ClusterSummary | null)[] | null;
	StdDevThreshold: number;
	K: number;
}

export interface RegressionDetectionResponse {
	summary: ClusterSummaries | null;
	frame: FrameResponse | null;
}

export interface TryBotRequest {
	kind: TryBotRequestKind;
	cl: CL;
	patch_number: number;
	commit_number: CommitNumber;
	query: string;
}

export interface TryBotResult {
	params: Params;
	median: number;
	lower: number;
	upper: number;
	stddevRatio: number;
	values: number[] | null;
}

export interface TryBotResponse {
	header: (ColumnHeader | null)[] | null;
	results: TryBotResult[] | null;
	paramset: ReadOnlyParamSet;
}

export namespace progress {
	export interface Message {
		key: string;
		value: string;
	}
}

export namespace progress {
	export interface SerializedProgress {
		status: progress.Status;
		messages: progress.Message[];
		results?: any;
		url: string;
	}
}

export namespace ingest {
	export interface SingleMeasurement {
		value: string;
		measurement: number;
		links?: { [key: string]: string } | null;
	}
}

export namespace ingest {
	export interface Result {
		key: { [key: string]: string } | null;
		measurement?: number;
		measurements?: { [key: string]: ingest.SingleMeasurement[] | null } | null;
	}
}

export namespace ingest {
	export interface Format {
		version: number;
		git_hash: string;
		issue?: CL;
		patchset?: string;
		key?: { [key: string]: string } | null;
		results: ingest.Result[] | null;
		links?: { [key: string]: string } | null;
	}
}

export type Params = { [key: string]: string } & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_paramsBrand: 'type alias for { [key: string]: string }'
};

export function Params(v: { [key: string]: string }): Params {
	return v as Params;
};

export type ParamSet = { [key: string]: string[] } & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_paramSetBrand: 'type alias for { [key: string]: string[] }'
};

export function ParamSet(v: { [key: string]: string[] }): ParamSet {
	return v as ParamSet;
};

export type ReadOnlyParamSet = { [key: string]: string[] } & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_readOnlyParamSetBrand: 'type alias for { [key: string]: string[] }'
};

export function ReadOnlyParamSet(v: { [key: string]: string[] }): ReadOnlyParamSet {
	return v as ReadOnlyParamSet;
};

export type Trace = number[] & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_traceBrand: 'type alias for number[]'
};

export function Trace(v: number[]): Trace {
	return v as Trace;
};

export type TraceSet = { [key: string]: Trace } & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_traceSetBrand: 'type alias for { [key: string]: Trace }'
};

export function TraceSet(v: { [key: string]: Trace }): TraceSet {
	return v as TraceSet;
};

export namespace pivot { export type Operation = 'sum' | 'avg' | 'geo' | 'std' | 'count' | 'min' | 'max'; }

export type SerializesToString = string & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_serializesToStringBrand: 'type alias for string'
};

export function SerializesToString(v: string): SerializesToString {
	return v as SerializesToString;
};

export type ClusterAlgo = 'kmeans' | 'stepfit';

export type StepDetection = '' | 'absolute' | 'const' | 'percent' | 'cohen' | 'mannwhitneyu';

export type ConfigState = 'ACTIVE' | 'DELETED';

export type Direction = 'UP' | 'DOWN' | 'BOTH';

export type AlertAction = 'noaction' | 'report' | 'bisect';

export type StepFitStatus = 'Low' | 'High' | 'Uninteresting';

export type CommitNumber = number & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_commitNumberBrand: 'type alias for number'
};

export function CommitNumber(v: number): CommitNumber {
	return v as CommitNumber;
};

export type TimestampSeconds = number & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_timestampSecondsBrand: 'type alias for number'
};

export function TimestampSeconds(v: number): TimestampSeconds {
	return v as TimestampSeconds;
};

export type CacheType = string & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_cacheTypeBrand: 'type alias for string'
};

export function CacheType(v: string): CacheType {
	return v as CacheType;
};

export type FrameResponseDisplayMode = 'display_query_only' | 'display_plot' | 'display_pivot_table' | 'display_pivot_plot' | 'display_spinner';

export type CommitNumberAnomalyMap = { [key: number]: Anomaly } | null;

export type AnomalyMap = { [key: string]: CommitNumberAnomalyMap } | null;

export type Status = '' | 'positive' | 'negative' | 'untriaged';

export type RequestType = 0 | 1;

export type Subset = 'all' | 'regressions' | 'untriaged';

export type NotifierTypes = 'html_email' | 'markdown_issuetracker' | 'none';

export type TraceFormat = 'chrome' | '';

export type Int64s = number[] | null;

export type Header = { [key: string]: string[] | null } | null;

export type TryBotRequestKind = 'trybot' | 'commit';

export type CL = string & {
	/**
	* WARNING: Do not reference this field from application code.
	*
	* This field exists solely to provide nominal typing. For reference, see
	* https://www.typescriptlang.org/play#example/nominal-typing.
	*/
	_cLBrand: 'type alias for string'
};

export function CL(v: string): CL {
	return v as CL;
};

export type ProcessState = 'Running' | 'Success' | 'Error';

export type ProjectId = 'chromium';

export namespace progress { export type Status = 'Running' | 'Finished' | 'Error'; }
