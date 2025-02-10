// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.21.12
// source: culprit_service.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request object for PersistCulprit rpc.
type PersistCulpritRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of commits identified as culprits.
	Commits []*Commit `protobuf:"bytes,1,rep,name=commits,proto3" json:"commits,omitempty"`
	// ID of the anomaly group corresponding to the bisection.
	AnomalyGroupId string `protobuf:"bytes,2,opt,name=anomaly_group_id,json=anomalyGroupId,proto3" json:"anomaly_group_id,omitempty"`
}

func (x *PersistCulpritRequest) Reset() {
	*x = PersistCulpritRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistCulpritRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistCulpritRequest) ProtoMessage() {}

func (x *PersistCulpritRequest) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistCulpritRequest.ProtoReflect.Descriptor instead.
func (*PersistCulpritRequest) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{0}
}

func (x *PersistCulpritRequest) GetCommits() []*Commit {
	if x != nil {
		return x.Commits
	}
	return nil
}

func (x *PersistCulpritRequest) GetAnomalyGroupId() string {
	if x != nil {
		return x.AnomalyGroupId
	}
	return ""
}

// Response object for PersistCulprit rpc.
type PersistCulpritResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of culprit ids created.
	CulpritIds []string `protobuf:"bytes,1,rep,name=culprit_ids,json=culpritIds,proto3" json:"culprit_ids,omitempty"`
}

func (x *PersistCulpritResponse) Reset() {
	*x = PersistCulpritResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistCulpritResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistCulpritResponse) ProtoMessage() {}

func (x *PersistCulpritResponse) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistCulpritResponse.ProtoReflect.Descriptor instead.
func (*PersistCulpritResponse) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{1}
}

func (x *PersistCulpritResponse) GetCulpritIds() []string {
	if x != nil {
		return x.CulpritIds
	}
	return nil
}

// Request object for GetCulprit rpc.
type GetCulpritRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CulpritIds []string `protobuf:"bytes,1,rep,name=culprit_ids,json=culpritIds,proto3" json:"culprit_ids,omitempty"`
}

func (x *GetCulpritRequest) Reset() {
	*x = GetCulpritRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCulpritRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCulpritRequest) ProtoMessage() {}

func (x *GetCulpritRequest) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCulpritRequest.ProtoReflect.Descriptor instead.
func (*GetCulpritRequest) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetCulpritRequest) GetCulpritIds() []string {
	if x != nil {
		return x.CulpritIds
	}
	return nil
}

// Response object for GetCulprit rpc.
type GetCulpritResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Culprits []*Culprit `protobuf:"bytes,1,rep,name=culprits,proto3" json:"culprits,omitempty"`
}

func (x *GetCulpritResponse) Reset() {
	*x = GetCulpritResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCulpritResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCulpritResponse) ProtoMessage() {}

func (x *GetCulpritResponse) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCulpritResponse.ProtoReflect.Descriptor instead.
func (*GetCulpritResponse) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetCulpritResponse) GetCulprits() []*Culprit {
	if x != nil {
		return x.Culprits
	}
	return nil
}

// Request object for NotifyUserOfAnomaly rpc.
type NotifyUserOfAnomalyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AnomalyGroupId string     `protobuf:"bytes,1,opt,name=anomaly_group_id,json=anomalyGroupId,proto3" json:"anomaly_group_id,omitempty"`
	Anomaly        []*Anomaly `protobuf:"bytes,2,rep,name=anomaly,proto3" json:"anomaly,omitempty"`
}

func (x *NotifyUserOfAnomalyRequest) Reset() {
	*x = NotifyUserOfAnomalyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserOfAnomalyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserOfAnomalyRequest) ProtoMessage() {}

func (x *NotifyUserOfAnomalyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserOfAnomalyRequest.ProtoReflect.Descriptor instead.
func (*NotifyUserOfAnomalyRequest) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{4}
}

func (x *NotifyUserOfAnomalyRequest) GetAnomalyGroupId() string {
	if x != nil {
		return x.AnomalyGroupId
	}
	return ""
}

func (x *NotifyUserOfAnomalyRequest) GetAnomaly() []*Anomaly {
	if x != nil {
		return x.Anomaly
	}
	return nil
}

// Response object for NotifyUserOfAnomaly rpc.
type NotifyUserOfAnomalyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id of issue created
	IssueId string `protobuf:"bytes,1,opt,name=issue_id,json=issueId,proto3" json:"issue_id,omitempty"`
}

func (x *NotifyUserOfAnomalyResponse) Reset() {
	*x = NotifyUserOfAnomalyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserOfAnomalyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserOfAnomalyResponse) ProtoMessage() {}

func (x *NotifyUserOfAnomalyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserOfAnomalyResponse.ProtoReflect.Descriptor instead.
func (*NotifyUserOfAnomalyResponse) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{5}
}

func (x *NotifyUserOfAnomalyResponse) GetIssueId() string {
	if x != nil {
		return x.IssueId
	}
	return ""
}

// Request object for NotifyUserOfCulprit rpc.
type NotifyUserOfCulpritRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of culprit ids.
	CulpritIds []string `protobuf:"bytes,1,rep,name=culprit_ids,json=culpritIds,proto3" json:"culprit_ids,omitempty"`
	// ID of the anomaly group corresponding to the bisection.
	AnomalyGroupId string `protobuf:"bytes,2,opt,name=anomaly_group_id,json=anomalyGroupId,proto3" json:"anomaly_group_id,omitempty"`
}

func (x *NotifyUserOfCulpritRequest) Reset() {
	*x = NotifyUserOfCulpritRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserOfCulpritRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserOfCulpritRequest) ProtoMessage() {}

func (x *NotifyUserOfCulpritRequest) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserOfCulpritRequest.ProtoReflect.Descriptor instead.
func (*NotifyUserOfCulpritRequest) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{6}
}

func (x *NotifyUserOfCulpritRequest) GetCulpritIds() []string {
	if x != nil {
		return x.CulpritIds
	}
	return nil
}

func (x *NotifyUserOfCulpritRequest) GetAnomalyGroupId() string {
	if x != nil {
		return x.AnomalyGroupId
	}
	return ""
}

// Response object for NotifyUserOfCulprit rpc.
type NotifyUserOfCulpritResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of issue ids created
	IssueIds []string `protobuf:"bytes,1,rep,name=issue_ids,json=issueIds,proto3" json:"issue_ids,omitempty"`
}

func (x *NotifyUserOfCulpritResponse) Reset() {
	*x = NotifyUserOfCulpritResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserOfCulpritResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserOfCulpritResponse) ProtoMessage() {}

func (x *NotifyUserOfCulpritResponse) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserOfCulpritResponse.ProtoReflect.Descriptor instead.
func (*NotifyUserOfCulpritResponse) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{7}
}

func (x *NotifyUserOfCulpritResponse) GetIssueIds() []string {
	if x != nil {
		return x.IssueIds
	}
	return nil
}

// Anomaly detected in a test
// Note: This is (right now) same as Anomaly object in anomalygroup_service.proto. But this has been
// duplicated because the two protos are used in two different service(culprit & anomlaygroup), and
// these services can evolve independently
type Anomaly struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the start commit position of the detected anomaly
	StartCommit int64 `protobuf:"varint,1,opt,name=start_commit,json=startCommit,proto3" json:"start_commit,omitempty"`
	// the end commit position of the detected anomaly
	EndCommit int64 `protobuf:"varint,2,opt,name=end_commit,json=endCommit,proto3" json:"end_commit,omitempty"`
	// the paramset from the regression detected in Skia. The parameters
	// are used in Skia alerts to define which tests to apply the deteciton
	// rules.
	// In chromeperf's context, it should include the following keys:
	//   - bot:
	//     the name of the bot (a.k.a, 'builder' in waterfall, and
	//     'configuration' in pinpoint job page.)
	//   - benchmark:
	//     the name of the benchmark
	//   - story:
	//     the name of the story (a.k.a., test)
	//   - measurement:
	//     the metric to look at. (a.k.a., 'test' in skia query ui,
	//     and 'chart' in pinpoint job page)
	//   - stat:
	//     the aggregation method on the data points
	Paramset map[string]string `protobuf:"bytes,3,rep,name=paramset,proto3" json:"paramset,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// indicate the direction towards which the change should be
	// considered as regression.
	// The possible values are: UP, DOWN or UNKNOWN
	ImprovementDirection string `protobuf:"bytes,4,opt,name=improvement_direction,json=improvementDirection,proto3" json:"improvement_direction,omitempty"`
}

func (x *Anomaly) Reset() {
	*x = Anomaly{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Anomaly) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Anomaly) ProtoMessage() {}

func (x *Anomaly) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Anomaly.ProtoReflect.Descriptor instead.
func (*Anomaly) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{8}
}

func (x *Anomaly) GetStartCommit() int64 {
	if x != nil {
		return x.StartCommit
	}
	return 0
}

func (x *Anomaly) GetEndCommit() int64 {
	if x != nil {
		return x.EndCommit
	}
	return 0
}

func (x *Anomaly) GetParamset() map[string]string {
	if x != nil {
		return x.Paramset
	}
	return nil
}

func (x *Anomaly) GetImprovementDirection() string {
	if x != nil {
		return x.ImprovementDirection
	}
	return ""
}

// Represents the change which has been identified as a culprit.
// TODO(wenbinzhang): remove anomaly group ids and issue ids as we have
// the info needed the group issue map
type Culprit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Commit          *Commit           `protobuf:"bytes,2,opt,name=commit,proto3" json:"commit,omitempty"`
	AnomalyGroupIds []string          `protobuf:"bytes,3,rep,name=anomaly_group_ids,json=anomalyGroupIds,proto3" json:"anomaly_group_ids,omitempty"`
	IssueIds        []string          `protobuf:"bytes,4,rep,name=issue_ids,json=issueIds,proto3" json:"issue_ids,omitempty"`
	GroupIssueMap   map[string]string `protobuf:"bytes,5,rep,name=group_issue_map,json=groupIssueMap,proto3" json:"group_issue_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Culprit) Reset() {
	*x = Culprit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Culprit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Culprit) ProtoMessage() {}

func (x *Culprit) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Culprit.ProtoReflect.Descriptor instead.
func (*Culprit) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{9}
}

func (x *Culprit) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Culprit) GetCommit() *Commit {
	if x != nil {
		return x.Commit
	}
	return nil
}

func (x *Culprit) GetAnomalyGroupIds() []string {
	if x != nil {
		return x.AnomalyGroupIds
	}
	return nil
}

func (x *Culprit) GetIssueIds() []string {
	if x != nil {
		return x.IssueIds
	}
	return nil
}

func (x *Culprit) GetGroupIssueMap() map[string]string {
	if x != nil {
		return x.GroupIssueMap
	}
	return nil
}

// Represents a commit which has been identified as a culprit.
type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Repo host e.g. chromium.googlesource.com
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	// Project inside the repo e.g. chromium/src
	Project string `protobuf:"bytes,2,opt,name=project,proto3" json:"project,omitempty"`
	// Repo ref e.g. "refs/heads/main
	Ref string `protobuf:"bytes,3,opt,name=ref,proto3" json:"ref,omitempty"`
	// Commit Hash
	Revision string `protobuf:"bytes,4,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_culprit_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_culprit_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_culprit_service_proto_rawDescGZIP(), []int{10}
}

func (x *Commit) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Commit) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *Commit) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *Commit) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

var File_culprit_service_proto protoreflect.FileDescriptor

var file_culprit_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x2e, 0x76, 0x31, 0x22, 0x6f, 0x0a, 0x15, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75,
	0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x6e,
	0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x22, 0x39, 0x0a, 0x16, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43,
	0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x49, 0x64, 0x73, 0x22,
	0x34, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x5f,
	0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x6c, 0x70, 0x72,
	0x69, 0x74, 0x49, 0x64, 0x73, 0x22, 0x45, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x75, 0x6c, 0x70,
	0x72, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x63,
	0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x75, 0x6c, 0x70, 0x72,
	0x69, 0x74, 0x52, 0x08, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x73, 0x22, 0x75, 0x0a, 0x1a,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x41, 0x6e, 0x6f, 0x6d,
	0x61, 0x6c, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x6e,
	0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x07, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x52, 0x07, 0x61, 0x6e, 0x6f, 0x6d,
	0x61, 0x6c, 0x79, 0x22, 0x38, 0x0a, 0x1b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65,
	0x72, 0x4f, 0x66, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x73, 0x73, 0x75, 0x65, 0x49, 0x64, 0x22, 0x67, 0x0a,
	0x1a, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x43, 0x75, 0x6c,
	0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0a, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x49, 0x64, 0x73, 0x12, 0x28, 0x0a, 0x10,
	0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x1b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x69,
	0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x49,
	0x64, 0x73, 0x22, 0xfc, 0x01, 0x0a, 0x07, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x12, 0x21,
	0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x6e, 0x64, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x12, 0x3d, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x65, 0x74,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x65, 0x74, 0x12,
	0x33, 0x0a, 0x15, 0x69, 0x6d, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14,
	0x69, 0x6d, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x3b, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x65, 0x74,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0xa0, 0x02, 0x0a, 0x07, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2a, 0x0a,
	0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x61, 0x6e, 0x6f,
	0x6d, 0x61, 0x6c, 0x79, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x69,
	0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x49,
	0x64, 0x73, 0x12, 0x4e, 0x0a, 0x0f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x73, 0x73, 0x75,
	0x65, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63, 0x75,
	0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x73, 0x73, 0x75, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0d, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x73, 0x73, 0x75, 0x65, 0x4d,
	0x61, 0x70, 0x1a, 0x40, 0x0a, 0x12, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x64, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x72, 0x65, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0x8e, 0x03, 0x0a, 0x0e, 0x43,
	0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a,
	0x0e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12,
	0x21, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x72,
	0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43,
	0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12, 0x1d, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x68, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x12, 0x26,
	0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66,
	0x41, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x68, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f,
	0x66, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12, 0x26, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72,
	0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x4f, 0x66, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x27, 0x2e, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x66, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a, 0x67,
	0x6f, 0x2e, 0x73, 0x6b, 0x69, 0x61, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61,
	0x2f, 0x70, 0x65, 0x72, 0x66, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_culprit_service_proto_rawDescOnce sync.Once
	file_culprit_service_proto_rawDescData = file_culprit_service_proto_rawDesc
)

func file_culprit_service_proto_rawDescGZIP() []byte {
	file_culprit_service_proto_rawDescOnce.Do(func() {
		file_culprit_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_culprit_service_proto_rawDescData)
	})
	return file_culprit_service_proto_rawDescData
}

var file_culprit_service_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_culprit_service_proto_goTypes = []interface{}{
	(*PersistCulpritRequest)(nil),       // 0: culprit.v1.PersistCulpritRequest
	(*PersistCulpritResponse)(nil),      // 1: culprit.v1.PersistCulpritResponse
	(*GetCulpritRequest)(nil),           // 2: culprit.v1.GetCulpritRequest
	(*GetCulpritResponse)(nil),          // 3: culprit.v1.GetCulpritResponse
	(*NotifyUserOfAnomalyRequest)(nil),  // 4: culprit.v1.NotifyUserOfAnomalyRequest
	(*NotifyUserOfAnomalyResponse)(nil), // 5: culprit.v1.NotifyUserOfAnomalyResponse
	(*NotifyUserOfCulpritRequest)(nil),  // 6: culprit.v1.NotifyUserOfCulpritRequest
	(*NotifyUserOfCulpritResponse)(nil), // 7: culprit.v1.NotifyUserOfCulpritResponse
	(*Anomaly)(nil),                     // 8: culprit.v1.Anomaly
	(*Culprit)(nil),                     // 9: culprit.v1.Culprit
	(*Commit)(nil),                      // 10: culprit.v1.Commit
	nil,                                 // 11: culprit.v1.Anomaly.ParamsetEntry
	nil,                                 // 12: culprit.v1.Culprit.GroupIssueMapEntry
}
var file_culprit_service_proto_depIdxs = []int32{
	10, // 0: culprit.v1.PersistCulpritRequest.commits:type_name -> culprit.v1.Commit
	9,  // 1: culprit.v1.GetCulpritResponse.culprits:type_name -> culprit.v1.Culprit
	8,  // 2: culprit.v1.NotifyUserOfAnomalyRequest.anomaly:type_name -> culprit.v1.Anomaly
	11, // 3: culprit.v1.Anomaly.paramset:type_name -> culprit.v1.Anomaly.ParamsetEntry
	10, // 4: culprit.v1.Culprit.commit:type_name -> culprit.v1.Commit
	12, // 5: culprit.v1.Culprit.group_issue_map:type_name -> culprit.v1.Culprit.GroupIssueMapEntry
	0,  // 6: culprit.v1.CulpritService.PersistCulprit:input_type -> culprit.v1.PersistCulpritRequest
	2,  // 7: culprit.v1.CulpritService.GetCulprit:input_type -> culprit.v1.GetCulpritRequest
	4,  // 8: culprit.v1.CulpritService.NotifyUserOfAnomaly:input_type -> culprit.v1.NotifyUserOfAnomalyRequest
	6,  // 9: culprit.v1.CulpritService.NotifyUserOfCulprit:input_type -> culprit.v1.NotifyUserOfCulpritRequest
	1,  // 10: culprit.v1.CulpritService.PersistCulprit:output_type -> culprit.v1.PersistCulpritResponse
	3,  // 11: culprit.v1.CulpritService.GetCulprit:output_type -> culprit.v1.GetCulpritResponse
	5,  // 12: culprit.v1.CulpritService.NotifyUserOfAnomaly:output_type -> culprit.v1.NotifyUserOfAnomalyResponse
	7,  // 13: culprit.v1.CulpritService.NotifyUserOfCulprit:output_type -> culprit.v1.NotifyUserOfCulpritResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_culprit_service_proto_init() }
func file_culprit_service_proto_init() {
	if File_culprit_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_culprit_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistCulpritRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistCulpritResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCulpritRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCulpritResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserOfAnomalyRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserOfAnomalyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserOfCulpritRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserOfCulpritResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Anomaly); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Culprit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_culprit_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_culprit_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_culprit_service_proto_goTypes,
		DependencyIndexes: file_culprit_service_proto_depIdxs,
		MessageInfos:      file_culprit_service_proto_msgTypes,
	}.Build()
	File_culprit_service_proto = out.File
	file_culprit_service_proto_rawDesc = nil
	file_culprit_service_proto_goTypes = nil
	file_culprit_service_proto_depIdxs = nil
}
