// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: proto/service.proto

package proto

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
	Culprits []*Culprit `protobuf:"bytes,1,rep,name=culprits,proto3" json:"culprits,omitempty"`
	// ID of the anomaly group corresponding to the bisection.
	AnomalyGroupId string `protobuf:"bytes,2,opt,name=anomaly_group_id,json=anomalyGroupId,proto3" json:"anomaly_group_id,omitempty"`
}

func (x *PersistCulpritRequest) Reset() {
	*x = PersistCulpritRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistCulpritRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistCulpritRequest) ProtoMessage() {}

func (x *PersistCulpritRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[0]
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
	return file_proto_service_proto_rawDescGZIP(), []int{0}
}

func (x *PersistCulpritRequest) GetCulprits() []*Culprit {
	if x != nil {
		return x.Culprits
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

	// List of issue ids created
	CulpritIds []string `protobuf:"bytes,1,rep,name=culprit_ids,json=culpritIds,proto3" json:"culprit_ids,omitempty"`
	IssueIds   []string `protobuf:"bytes,2,rep,name=issue_ids,json=issueIds,proto3" json:"issue_ids,omitempty"`
}

func (x *PersistCulpritResponse) Reset() {
	*x = PersistCulpritResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistCulpritResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistCulpritResponse) ProtoMessage() {}

func (x *PersistCulpritResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[1]
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
	return file_proto_service_proto_rawDescGZIP(), []int{1}
}

func (x *PersistCulpritResponse) GetCulpritIds() []string {
	if x != nil {
		return x.CulpritIds
	}
	return nil
}

func (x *PersistCulpritResponse) GetIssueIds() []string {
	if x != nil {
		return x.IssueIds
	}
	return nil
}

// Request object for NotifyUser rpc.
type NotifyUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of commits identified as culprits.
	Culprits []*Culprit `protobuf:"bytes,1,rep,name=culprits,proto3" json:"culprits,omitempty"`
	// ID of the anomaly group corresponding to the bisection.
	AnomalyGroupId string `protobuf:"bytes,2,opt,name=anomaly_group_id,json=anomalyGroupId,proto3" json:"anomaly_group_id,omitempty"`
}

func (x *NotifyUserRequest) Reset() {
	*x = NotifyUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserRequest) ProtoMessage() {}

func (x *NotifyUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserRequest.ProtoReflect.Descriptor instead.
func (*NotifyUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{2}
}

func (x *NotifyUserRequest) GetCulprits() []*Culprit {
	if x != nil {
		return x.Culprits
	}
	return nil
}

func (x *NotifyUserRequest) GetAnomalyGroupId() string {
	if x != nil {
		return x.AnomalyGroupId
	}
	return ""
}

// Response object for NotifyUser rpc.
type NotifyUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of issue ids created
	IssueIds []string `protobuf:"bytes,1,rep,name=issue_ids,json=issueIds,proto3" json:"issue_ids,omitempty"`
}

func (x *NotifyUserResponse) Reset() {
	*x = NotifyUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyUserResponse) ProtoMessage() {}

func (x *NotifyUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyUserResponse.ProtoReflect.Descriptor instead.
func (*NotifyUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{3}
}

func (x *NotifyUserResponse) GetIssueIds() []string {
	if x != nil {
		return x.IssueIds
	}
	return nil
}

// Represents the change which has been identified as a culprit.
type Culprit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Commit *Commit `protobuf:"bytes,1,opt,name=commit,proto3" json:"commit,omitempty"`
}

func (x *Culprit) Reset() {
	*x = Culprit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Culprit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Culprit) ProtoMessage() {}

func (x *Culprit) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[4]
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
	return file_proto_service_proto_rawDescGZIP(), []int{4}
}

func (x *Culprit) GetCommit() *Commit {
	if x != nil {
		return x.Commit
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
		mi := &file_proto_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[5]
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
	return file_proto_service_proto_rawDescGZIP(), []int{5}
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

var File_proto_service_proto protoreflect.FileDescriptor

var file_proto_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x15,
	0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x08, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x73, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x6e, 0x6f,
	0x6d, 0x61, 0x6c, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x16, 0x50,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x6c, 0x70,
	0x72, 0x69, 0x74, 0x49, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f,
	0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x49, 0x64, 0x73, 0x22, 0x69, 0x0a, 0x11, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x63, 0x75, 0x6c, 0x70,
	0x72, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x08, 0x63, 0x75, 0x6c, 0x70,
	0x72, 0x69, 0x74, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x5f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x61, 0x6e, 0x6f, 0x6d, 0x61, 0x6c, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x22, 0x31,
	0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x69, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x49, 0x64,
	0x73, 0x22, 0x30, 0x0a, 0x07, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12, 0x25, 0x0a, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x06, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x22, 0x64, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x72,
	0x65, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x1a, 0x0a,
	0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0xa6, 0x01, 0x0a, 0x0e, 0x43, 0x75,
	0x6c, 0x70, 0x72, 0x69, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0e,
	0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x12, 0x1c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75,
	0x6c, 0x70, 0x72, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x43, 0x75, 0x6c, 0x70,
	0x72, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x43, 0x0a,
	0x0a, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x6f, 0x2e, 0x73, 0x6b, 0x69, 0x61, 0x2e, 0x6f, 0x72,
	0x67, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x70, 0x65, 0x72, 0x66, 0x2f, 0x67, 0x6f, 0x2f,
	0x63, 0x75, 0x6c, 0x70, 0x72, 0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_service_proto_rawDescOnce sync.Once
	file_proto_service_proto_rawDescData = file_proto_service_proto_rawDesc
)

func file_proto_service_proto_rawDescGZIP() []byte {
	file_proto_service_proto_rawDescOnce.Do(func() {
		file_proto_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_service_proto_rawDescData)
	})
	return file_proto_service_proto_rawDescData
}

var file_proto_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_service_proto_goTypes = []interface{}{
	(*PersistCulpritRequest)(nil),  // 0: proto.PersistCulpritRequest
	(*PersistCulpritResponse)(nil), // 1: proto.PersistCulpritResponse
	(*NotifyUserRequest)(nil),      // 2: proto.NotifyUserRequest
	(*NotifyUserResponse)(nil),     // 3: proto.NotifyUserResponse
	(*Culprit)(nil),                // 4: proto.Culprit
	(*Commit)(nil),                 // 5: proto.Commit
}
var file_proto_service_proto_depIdxs = []int32{
	4, // 0: proto.PersistCulpritRequest.culprits:type_name -> proto.Culprit
	4, // 1: proto.NotifyUserRequest.culprits:type_name -> proto.Culprit
	5, // 2: proto.Culprit.commit:type_name -> proto.Commit
	0, // 3: proto.CulpritService.PersistCulprit:input_type -> proto.PersistCulpritRequest
	2, // 4: proto.CulpritService.NotifyUser:input_type -> proto.NotifyUserRequest
	1, // 5: proto.CulpritService.PersistCulprit:output_type -> proto.PersistCulpritResponse
	3, // 6: proto.CulpritService.NotifyUser:output_type -> proto.NotifyUserResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_service_proto_init() }
func file_proto_service_proto_init() {
	if File_proto_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserRequest); i {
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
		file_proto_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyUserResponse); i {
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
		file_proto_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_proto_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_service_proto_goTypes,
		DependencyIndexes: file_proto_service_proto_depIdxs,
		MessageInfos:      file_proto_service_proto_msgTypes,
	}.Build()
	File_proto_service_proto = out.File
	file_proto_service_proto_rawDesc = nil
	file_proto_service_proto_goTypes = nil
	file_proto_service_proto_depIdxs = nil
}
