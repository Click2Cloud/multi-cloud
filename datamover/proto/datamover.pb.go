// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: datamover.proto

package datamover

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AbortJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *AbortJobRequest) Reset() {
	*x = AbortJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortJobRequest) ProtoMessage() {}

func (x *AbortJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortJobRequest.ProtoReflect.Descriptor instead.
func (*AbortJobRequest) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{0}
}

func (x *AbortJobRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type AbortJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *AbortJobResponse) Reset() {
	*x = AbortJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortJobResponse) ProtoMessage() {}

func (x *AbortJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortJobResponse.ProtoReflect.Descriptor instead.
func (*AbortJobResponse) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{1}
}

func (x *AbortJobResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *AbortJobResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AbortJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type PauseJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PauseJobRequest) Reset() {
	*x = PauseJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PauseJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PauseJobRequest) ProtoMessage() {}

func (x *PauseJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PauseJobRequest.ProtoReflect.Descriptor instead.
func (*PauseJobRequest) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{2}
}

func (x *PauseJobRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PauseJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err    string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PauseJobResponse) Reset() {
	*x = PauseJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PauseJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PauseJobResponse) ProtoMessage() {}

func (x *PauseJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PauseJobResponse.ProtoReflect.Descriptor instead.
func (*PauseJobResponse) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{3}
}

func (x *PauseJobResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *PauseJobResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PauseJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KV) Reset() {
	*x = KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KV) ProtoMessage() {}

func (x *KV) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KV.ProtoReflect.Descriptor instead.
func (*KV) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{4}
}

func (x *KV) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KV) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prefix string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Tag    []*KV  `protobuf:"bytes,2,rep,name=tag,proto3" json:"tag,omitempty"`
}

func (x *Filter) Reset() {
	*x = Filter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Filter) ProtoMessage() {}

func (x *Filter) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Filter.ProtoReflect.Descriptor instead.
func (*Filter) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{5}
}

func (x *Filter) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *Filter) GetTag() []*KV {
	if x != nil {
		return x.Tag
	}
	return nil
}

type Connector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"` //opensds,aws,azure,hwcloud,etc.
	BucketName  string `protobuf:"bytes,2,opt,name=BucketName,proto3" json:"BucketName,omitempty"`
	ConnConfig  []*KV  `protobuf:"bytes,3,rep,name=ConnConfig,proto3" json:"ConnConfig,omitempty"`
	BackendName string `protobuf:"bytes,4,opt,name=BackendName,proto3" json:"BackendName,omitempty"`
}

func (x *Connector) Reset() {
	*x = Connector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connector) ProtoMessage() {}

func (x *Connector) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connector.ProtoReflect.Descriptor instead.
func (*Connector) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{6}
}

func (x *Connector) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Connector) GetBucketName() string {
	if x != nil {
		return x.BucketName
	}
	return ""
}

func (x *Connector) GetConnConfig() []*KV {
	if x != nil {
		return x.ConnConfig
	}
	return nil
}

func (x *Connector) GetBackendName() string {
	if x != nil {
		return x.BackendName
	}
	return ""
}

type RunJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TenanId      string     `protobuf:"bytes,2,opt,name=tenanId,proto3" json:"tenanId,omitempty"`
	UserId       string     `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	SourceConn   *Connector `protobuf:"bytes,4,opt,name=sourceConn,proto3" json:"sourceConn,omitempty"`
	DestConn     *Connector `protobuf:"bytes,5,opt,name=destConn,proto3" json:"destConn,omitempty"`
	Filt         *Filter    `protobuf:"bytes,6,opt,name=filt,proto3" json:"filt,omitempty"`
	RemainSource bool       `protobuf:"varint,7,opt,name=remainSource,proto3" json:"remainSource,omitempty"`
}

func (x *RunJobRequest) Reset() {
	*x = RunJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunJobRequest) ProtoMessage() {}

func (x *RunJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunJobRequest.ProtoReflect.Descriptor instead.
func (*RunJobRequest) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{7}
}

func (x *RunJobRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RunJobRequest) GetTenanId() string {
	if x != nil {
		return x.TenanId
	}
	return ""
}

func (x *RunJobRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RunJobRequest) GetSourceConn() *Connector {
	if x != nil {
		return x.SourceConn
	}
	return nil
}

func (x *RunJobRequest) GetDestConn() *Connector {
	if x != nil {
		return x.DestConn
	}
	return nil
}

func (x *RunJobRequest) GetFilt() *Filter {
	if x != nil {
		return x.Filt
	}
	return nil
}

func (x *RunJobRequest) GetRemainSource() bool {
	if x != nil {
		return x.RemainSource
	}
	return false
}

type RunJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *RunJobResponse) Reset() {
	*x = RunJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunJobResponse) ProtoMessage() {}

func (x *RunJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunJobResponse.ProtoReflect.Descriptor instead.
func (*RunJobResponse) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{8}
}

func (x *RunJobResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type LifecycleActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjKey        string `protobuf:"bytes,1,opt,name=objKey,proto3" json:"objKey,omitempty"`                // for transition and expiration
	BucketName    string `protobuf:"bytes,2,opt,name=bucketName,proto3" json:"bucketName,omitempty"`        // for transition and expiration
	VersionId     string `protobuf:"bytes,3,opt,name=versionId,proto3" json:"versionId,omitempty"`          // for transition and expiration
	StorageMeta   string `protobuf:"bytes,4,opt,name=storageMeta,proto3" json:"storageMeta,omitempty"`      // for transition and expiration
	ObjectId      string `protobuf:"bytes,5,opt,name=objectId,proto3" json:"objectId,omitempty"`            // for transition and expiration
	Action        int32  `protobuf:"varint,6,opt,name=action,proto3" json:"action,omitempty"`               // 0-Expiration, 1-IncloudTransition, 2-CrossCloudTransition, 3-AbortMultipartUpload
	SourceTier    int32  `protobuf:"varint,7,opt,name=sourceTier,proto3" json:"sourceTier,omitempty"`       // only for transition
	TargetTier    int32  `protobuf:"varint,8,opt,name=targetTier,proto3" json:"targetTier,omitempty"`       // only for transtion
	SourceBackend string `protobuf:"bytes,9,opt,name=sourceBackend,proto3" json:"sourceBackend,omitempty"`  //for transition and expiration
	TargetBackend string `protobuf:"bytes,10,opt,name=targetBackend,proto3" json:"targetBackend,omitempty"` //for transition and abort incomplete multipart upload
	TargetBucket  string `protobuf:"bytes,11,opt,name=targetBucket,proto3" json:"targetBucket,omitempty"`   // for transition to specific backend
	ObjSize       int64  `protobuf:"varint,12,opt,name=objSize,proto3" json:"objSize,omitempty"`            // for transition
	UploadId      string `protobuf:"bytes,13,opt,name=uploadId,proto3" json:"uploadId,omitempty"`           // only for abort incomplete multipart upload
}

func (x *LifecycleActionRequest) Reset() {
	*x = LifecycleActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LifecycleActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LifecycleActionRequest) ProtoMessage() {}

func (x *LifecycleActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LifecycleActionRequest.ProtoReflect.Descriptor instead.
func (*LifecycleActionRequest) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{9}
}

func (x *LifecycleActionRequest) GetObjKey() string {
	if x != nil {
		return x.ObjKey
	}
	return ""
}

func (x *LifecycleActionRequest) GetBucketName() string {
	if x != nil {
		return x.BucketName
	}
	return ""
}

func (x *LifecycleActionRequest) GetVersionId() string {
	if x != nil {
		return x.VersionId
	}
	return ""
}

func (x *LifecycleActionRequest) GetStorageMeta() string {
	if x != nil {
		return x.StorageMeta
	}
	return ""
}

func (x *LifecycleActionRequest) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

func (x *LifecycleActionRequest) GetAction() int32 {
	if x != nil {
		return x.Action
	}
	return 0
}

func (x *LifecycleActionRequest) GetSourceTier() int32 {
	if x != nil {
		return x.SourceTier
	}
	return 0
}

func (x *LifecycleActionRequest) GetTargetTier() int32 {
	if x != nil {
		return x.TargetTier
	}
	return 0
}

func (x *LifecycleActionRequest) GetSourceBackend() string {
	if x != nil {
		return x.SourceBackend
	}
	return ""
}

func (x *LifecycleActionRequest) GetTargetBackend() string {
	if x != nil {
		return x.TargetBackend
	}
	return ""
}

func (x *LifecycleActionRequest) GetTargetBucket() string {
	if x != nil {
		return x.TargetBucket
	}
	return ""
}

func (x *LifecycleActionRequest) GetObjSize() int64 {
	if x != nil {
		return x.ObjSize
	}
	return 0
}

func (x *LifecycleActionRequest) GetUploadId() string {
	if x != nil {
		return x.UploadId
	}
	return ""
}

type LifecycleActionResonse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *LifecycleActionResonse) Reset() {
	*x = LifecycleActionResonse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datamover_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LifecycleActionResonse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LifecycleActionResonse) ProtoMessage() {}

func (x *LifecycleActionResonse) ProtoReflect() protoreflect.Message {
	mi := &file_datamover_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LifecycleActionResonse.ProtoReflect.Descriptor instead.
func (*LifecycleActionResonse) Descriptor() ([]byte, []int) {
	return file_datamover_proto_rawDescGZIP(), []int{10}
}

func (x *LifecycleActionResonse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_datamover_proto protoreflect.FileDescriptor

var file_datamover_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x64, 0x61, 0x74, 0x61, 0x6d, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x21, 0x0a, 0x0f, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x4c, 0x0a, 0x10, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x21, 0x0a, 0x0f, 0x50, 0x61, 0x75, 0x73, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4c, 0x0a, 0x10, 0x50, 0x61, 0x75, 0x73, 0x65, 0x4a, 0x6f,
	0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x2c, 0x0a, 0x02, 0x4b, 0x56, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x37, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x70,
	0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x12, 0x15, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x03, 0x2e, 0x4b, 0x56, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0x86, 0x01, 0x0a, 0x09, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0a,
	0x43, 0x6f, 0x6e, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x03, 0x2e, 0x4b, 0x56, 0x52, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x20, 0x0a, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0xe6, 0x01, 0x0a, 0x0d, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43,
	0x6f, 0x6e, 0x6e, 0x12, 0x26, 0x0a, 0x08, 0x64, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x08, 0x64, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x1b, 0x0a, 0x04, 0x66,
	0x69, 0x6c, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x6d, 0x61,
	0x69, 0x6e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x22, 0x0a, 0x0e,
	0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72,
	0x22, 0xaa, 0x03, 0x0a, 0x16, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f,
	0x62, 0x6a, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x62, 0x6a,
	0x4b, 0x65, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x54, 0x69, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x54, 0x69, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x54, 0x69, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x54, 0x69, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x12, 0x24, 0x0a,
	0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x62, 0x6a, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x22, 0x2a, 0x0a,
	0x16, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x32, 0xe7, 0x01, 0x0a, 0x09, 0x64, 0x61,
	0x74, 0x61, 0x6d, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x06, 0x52, 0x75, 0x6e, 0x6a, 0x6f,
	0x62, 0x12, 0x0e, 0x2e, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x11, 0x44, 0x6f, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79,
	0x63, 0x6c, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x2e, 0x4c, 0x69, 0x66, 0x65,
	0x63, 0x79, 0x63, 0x6c, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x4c, 0x69, 0x66, 0x65, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x31, 0x0a,
	0x08, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x10, 0x2e, 0x41, 0x62, 0x6f, 0x72,
	0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x41, 0x62,
	0x6f, 0x72, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x31, 0x0a, 0x08, 0x50, 0x61, 0x75, 0x73, 0x65, 0x4a, 0x6f, 0x62, 0x12, 0x10, 0x2e, 0x50,
	0x61, 0x75, 0x73, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11,
	0x2e, 0x50, 0x61, 0x75, 0x73, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x6d, 0x6f, 0x76,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_datamover_proto_rawDescOnce sync.Once
	file_datamover_proto_rawDescData = file_datamover_proto_rawDesc
)

func file_datamover_proto_rawDescGZIP() []byte {
	file_datamover_proto_rawDescOnce.Do(func() {
		file_datamover_proto_rawDescData = protoimpl.X.CompressGZIP(file_datamover_proto_rawDescData)
	})
	return file_datamover_proto_rawDescData
}

var file_datamover_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_datamover_proto_goTypes = []interface{}{
	(*AbortJobRequest)(nil),        // 0: AbortJobRequest
	(*AbortJobResponse)(nil),       // 1: AbortJobResponse
	(*PauseJobRequest)(nil),        // 2: PauseJobRequest
	(*PauseJobResponse)(nil),       // 3: PauseJobResponse
	(*KV)(nil),                     // 4: KV
	(*Filter)(nil),                 // 5: Filter
	(*Connector)(nil),              // 6: Connector
	(*RunJobRequest)(nil),          // 7: RunJobRequest
	(*RunJobResponse)(nil),         // 8: RunJobResponse
	(*LifecycleActionRequest)(nil), // 9: LifecycleActionRequest
	(*LifecycleActionResonse)(nil), // 10: LifecycleActionResonse
}
var file_datamover_proto_depIdxs = []int32{
	4,  // 0: Filter.tag:type_name -> KV
	4,  // 1: Connector.ConnConfig:type_name -> KV
	6,  // 2: RunJobRequest.sourceConn:type_name -> Connector
	6,  // 3: RunJobRequest.destConn:type_name -> Connector
	5,  // 4: RunJobRequest.filt:type_name -> Filter
	7,  // 5: datamover.Runjob:input_type -> RunJobRequest
	9,  // 6: datamover.DoLifecycleAction:input_type -> LifecycleActionRequest
	0,  // 7: datamover.AbortJob:input_type -> AbortJobRequest
	2,  // 8: datamover.PauseJob:input_type -> PauseJobRequest
	8,  // 9: datamover.Runjob:output_type -> RunJobResponse
	10, // 10: datamover.DoLifecycleAction:output_type -> LifecycleActionResonse
	1,  // 11: datamover.AbortJob:output_type -> AbortJobResponse
	3,  // 12: datamover.PauseJob:output_type -> PauseJobResponse
	9,  // [9:13] is the sub-list for method output_type
	5,  // [5:9] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_datamover_proto_init() }
func file_datamover_proto_init() {
	if File_datamover_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_datamover_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortJobRequest); i {
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
		file_datamover_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortJobResponse); i {
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
		file_datamover_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PauseJobRequest); i {
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
		file_datamover_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PauseJobResponse); i {
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
		file_datamover_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KV); i {
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
		file_datamover_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Filter); i {
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
		file_datamover_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connector); i {
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
		file_datamover_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunJobRequest); i {
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
		file_datamover_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunJobResponse); i {
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
		file_datamover_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LifecycleActionRequest); i {
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
		file_datamover_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LifecycleActionResonse); i {
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
			RawDescriptor: file_datamover_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_datamover_proto_goTypes,
		DependencyIndexes: file_datamover_proto_depIdxs,
		MessageInfos:      file_datamover_proto_msgTypes,
	}.Build()
	File_datamover_proto = out.File
	file_datamover_proto_rawDesc = nil
	file_datamover_proto_goTypes = nil
	file_datamover_proto_depIdxs = nil
}
