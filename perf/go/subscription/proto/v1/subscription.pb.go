// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: subscription.proto

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

// A subscription defines alerting configurations for anomalies detected.
type Subscription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique name identifying subscription.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// infra_internal Git hash on which a subscription is based on.
	Revision string `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
	// Labels to attach to bugs associated with a subscription.
	BugLabels []string `protobuf:"bytes,3,rep,name=bug_labels,json=bugLabels,proto3" json:"bug_labels,omitempty"`
	// Hotlists to add to bugs associated with a subscription.
	Hotlists []string `protobuf:"bytes,4,rep,name=hotlists,proto3" json:"hotlists,omitempty"`
	// Component in which to file bugs associated with a subscription.
	BugComponent string `protobuf:"bytes,5,opt,name=bug_component,json=bugComponent,proto3" json:"bug_component,omitempty"`
	// Priority of bugs associated with a subscription. Must be between 0-4.
	BugPriority int32 `protobuf:"varint,8,opt,name=bug_priority,json=bugPriority,proto3" json:"bug_priority,omitempty"`
	// Severity of bugs associated with a subscription. Must be between 0-4.
	BugSeverity int32 `protobuf:"varint,9,opt,name=bug_severity,json=bugSeverity,proto3" json:"bug_severity,omitempty"`
	// Emails to CC in bugs associated with a subscription.
	BugCcEmails []string `protobuf:"bytes,6,rep,name=bug_cc_emails,json=bugCcEmails,proto3" json:"bug_cc_emails,omitempty"`
	// Owner of subscription. Used for contact purposes.
	ContactEmail string `protobuf:"bytes,7,opt,name=contact_email,json=contactEmail,proto3" json:"contact_email,omitempty"`
}

func (x *Subscription) Reset() {
	*x = Subscription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_subscription_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscription) ProtoMessage() {}

func (x *Subscription) ProtoReflect() protoreflect.Message {
	mi := &file_subscription_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscription.ProtoReflect.Descriptor instead.
func (*Subscription) Descriptor() ([]byte, []int) {
	return file_subscription_proto_rawDescGZIP(), []int{0}
}

func (x *Subscription) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Subscription) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

func (x *Subscription) GetBugLabels() []string {
	if x != nil {
		return x.BugLabels
	}
	return nil
}

func (x *Subscription) GetHotlists() []string {
	if x != nil {
		return x.Hotlists
	}
	return nil
}

func (x *Subscription) GetBugComponent() string {
	if x != nil {
		return x.BugComponent
	}
	return ""
}

func (x *Subscription) GetBugPriority() int32 {
	if x != nil {
		return x.BugPriority
	}
	return 0
}

func (x *Subscription) GetBugSeverity() int32 {
	if x != nil {
		return x.BugSeverity
	}
	return 0
}

func (x *Subscription) GetBugCcEmails() []string {
	if x != nil {
		return x.BugCcEmails
	}
	return nil
}

func (x *Subscription) GetContactEmail() string {
	if x != nil {
		return x.ContactEmail
	}
	return ""
}

var File_subscription_proto protoreflect.FileDescriptor

var file_subscription_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0xad, 0x02, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x75, 0x67, 0x5f, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x62, 0x75, 0x67, 0x4c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x74, 0x6c, 0x69, 0x73, 0x74,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x74, 0x6c, 0x69, 0x73, 0x74,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x75, 0x67, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x62, 0x75, 0x67, 0x43, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x75, 0x67, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x62, 0x75,
	0x67, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x75, 0x67,
	0x5f, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x62, 0x75, 0x67, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x22, 0x0a, 0x0d,
	0x62, 0x75, 0x67, 0x5f, 0x63, 0x63, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x75, 0x67, 0x43, 0x63, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x5f, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x6f, 0x2e, 0x73, 0x6b, 0x69, 0x61,
	0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x70, 0x65, 0x72, 0x66, 0x2f,
	0x67, 0x6f, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_subscription_proto_rawDescOnce sync.Once
	file_subscription_proto_rawDescData = file_subscription_proto_rawDesc
)

func file_subscription_proto_rawDescGZIP() []byte {
	file_subscription_proto_rawDescOnce.Do(func() {
		file_subscription_proto_rawDescData = protoimpl.X.CompressGZIP(file_subscription_proto_rawDescData)
	})
	return file_subscription_proto_rawDescData
}

var file_subscription_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_subscription_proto_goTypes = []interface{}{
	(*Subscription)(nil), // 0: subscription.v1.Subscription
}
var file_subscription_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_subscription_proto_init() }
func file_subscription_proto_init() {
	if File_subscription_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_subscription_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscription); i {
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
			RawDescriptor: file_subscription_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_subscription_proto_goTypes,
		DependencyIndexes: file_subscription_proto_depIdxs,
		MessageInfos:      file_subscription_proto_msgTypes,
	}.Build()
	File_subscription_proto = out.File
	file_subscription_proto_rawDesc = nil
	file_subscription_proto_goTypes = nil
	file_subscription_proto_depIdxs = nil
}
