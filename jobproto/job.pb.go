// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: jobproto/job.proto

package jobproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TerminateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Request bool `protobuf:"varint,1,opt,name=request,proto3" json:"request,omitempty"`
}

func (x *TerminateRequest) Reset() {
	*x = TerminateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminateRequest) ProtoMessage() {}

func (x *TerminateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminateRequest.ProtoReflect.Descriptor instead.
func (*TerminateRequest) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{0}
}

func (x *TerminateRequest) GetRequest() bool {
	if x != nil {
		return x.Request
	}
	return false
}

type TerminateConfirmation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Confirmed         bool  `protobuf:"varint,1,opt,name=confirmed,proto3" json:"confirmed,omitempty"`
	NumOfJobCompleted int32 `protobuf:"varint,2,opt,name=numOfJobCompleted,proto3" json:"numOfJobCompleted,omitempty"`
}

func (x *TerminateConfirmation) Reset() {
	*x = TerminateConfirmation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminateConfirmation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminateConfirmation) ProtoMessage() {}

func (x *TerminateConfirmation) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminateConfirmation.ProtoReflect.Descriptor instead.
func (*TerminateConfirmation) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{1}
}

func (x *TerminateConfirmation) GetConfirmed() bool {
	if x != nil {
		return x.Confirmed
	}
	return false
}

func (x *TerminateConfirmation) GetNumOfJobCompleted() int32 {
	if x != nil {
		return x.NumOfJobCompleted
	}
	return 0
}

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Datafile string `protobuf:"bytes,1,opt,name=datafile,proto3" json:"datafile,omitempty"`
	Start    int32  `protobuf:"varint,2,opt,name=start,proto3" json:"start,omitempty"`
	Length   int32  `protobuf:"varint,3,opt,name=length,proto3" json:"length,omitempty"`
	CValue   int32  `protobuf:"varint,4,opt,name=cValue,proto3" json:"cValue,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{2}
}

func (x *Job) GetDatafile() string {
	if x != nil {
		return x.Datafile
	}
	return ""
}

func (x *Job) GetStart() int32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *Job) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *Job) GetCValue() int32 {
	if x != nil {
		return x.CValue
	}
	return 0
}

type JobResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobFound    *Job  `protobuf:"bytes,1,opt,name=jobFound,proto3" json:"jobFound,omitempty"`
	NumOfPrimes int32 `protobuf:"varint,2,opt,name=numOfPrimes,proto3" json:"numOfPrimes,omitempty"`
}

func (x *JobResult) Reset() {
	*x = JobResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobResult) ProtoMessage() {}

func (x *JobResult) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobResult.ProtoReflect.Descriptor instead.
func (*JobResult) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{3}
}

func (x *JobResult) GetJobFound() *Job {
	if x != nil {
		return x.JobFound
	}
	return nil
}

func (x *JobResult) GetNumOfPrimes() int32 {
	if x != nil {
		return x.NumOfPrimes
	}
	return 0
}

type PushInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result              *JobResult             `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	RequestConfirmation *TerminateConfirmation `protobuf:"bytes,2,opt,name=requestConfirmation,proto3" json:"requestConfirmation,omitempty"`
}

func (x *PushInfo) Reset() {
	*x = PushInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushInfo) ProtoMessage() {}

func (x *PushInfo) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushInfo.ProtoReflect.Descriptor instead.
func (*PushInfo) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{4}
}

func (x *PushInfo) GetResult() *JobResult {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *PushInfo) GetRequestConfirmation() *TerminateConfirmation {
	if x != nil {
		return x.RequestConfirmation
	}
	return nil
}

type Connected struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Connection bool `protobuf:"varint,1,opt,name=connection,proto3" json:"connection,omitempty"`
}

func (x *Connected) Reset() {
	*x = Connected{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jobproto_job_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connected) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connected) ProtoMessage() {}

func (x *Connected) ProtoReflect() protoreflect.Message {
	mi := &file_jobproto_job_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connected.ProtoReflect.Descriptor instead.
func (*Connected) Descriptor() ([]byte, []int) {
	return file_jobproto_job_proto_rawDescGZIP(), []int{5}
}

func (x *Connected) GetConnection() bool {
	if x != nil {
		return x.Connection
	}
	return false
}

var File_jobproto_job_proto protoreflect.FileDescriptor

var file_jobproto_job_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6a, 0x6f, 0x62, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a,
	0x10, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x63, 0x0a, 0x15, 0x54,
	0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d,
	0x65, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x4a, 0x6f, 0x62, 0x43, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x6e,
	0x75, 0x6d, 0x4f, 0x66, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x22, 0x67, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x66,
	0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5a, 0x0a, 0x09, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x46, 0x6f, 0x75,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x46, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x50, 0x72, 0x69, 0x6d,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x50,
	0x72, 0x69, 0x6d, 0x65, 0x73, 0x22, 0x8e, 0x01, 0x0a, 0x08, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x2d, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x53, 0x0a, 0x13, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x65, 0x72, 0x6d,
	0x69, 0x6e, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x13, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2b, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x32, 0xcb, 0x01, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x35, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a, 0x6f, 0x62,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4a, 0x6f, 0x62, 0x12, 0x40, 0x0a, 0x0a, 0x50, 0x75, 0x73,
	0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x1c, 0x2e,
	0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x65, 0x72, 0x6d, 0x69,
	0x6e, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x44, 0x0a, 0x13, 0x45,
	0x73, 0x74, 0x61, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x15, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x75, 0x6b, 0x74, 0x61, 0x72, 0x2d, 0x67, 0x69, 0x66, 0x2f, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x2d, 0x32, 0x2d, 0x36, 0x31, 0x31, 0x2f, 0x6a, 0x6f, 0x62, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_jobproto_job_proto_rawDescOnce sync.Once
	file_jobproto_job_proto_rawDescData = file_jobproto_job_proto_rawDesc
)

func file_jobproto_job_proto_rawDescGZIP() []byte {
	file_jobproto_job_proto_rawDescOnce.Do(func() {
		file_jobproto_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_jobproto_job_proto_rawDescData)
	})
	return file_jobproto_job_proto_rawDescData
}

var file_jobproto_job_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_jobproto_job_proto_goTypes = []interface{}{
	(*TerminateRequest)(nil),      // 0: jobservice.TerminateRequest
	(*TerminateConfirmation)(nil), // 1: jobservice.TerminateConfirmation
	(*Job)(nil),                   // 2: jobservice.Job
	(*JobResult)(nil),             // 3: jobservice.JobResult
	(*PushInfo)(nil),              // 4: jobservice.PushInfo
	(*Connected)(nil),             // 5: jobservice.Connected
	(*emptypb.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_jobproto_job_proto_depIdxs = []int32{
	2, // 0: jobservice.JobResult.jobFound:type_name -> jobservice.Job
	3, // 1: jobservice.PushInfo.result:type_name -> jobservice.JobResult
	1, // 2: jobservice.PushInfo.requestConfirmation:type_name -> jobservice.TerminateConfirmation
	6, // 3: jobservice.JobService.RequestJob:input_type -> google.protobuf.Empty
	4, // 4: jobservice.JobService.PushResult:input_type -> jobservice.PushInfo
	5, // 5: jobservice.JobService.EstablishConnection:input_type -> jobservice.Connected
	2, // 6: jobservice.JobService.RequestJob:output_type -> jobservice.Job
	0, // 7: jobservice.JobService.PushResult:output_type -> jobservice.TerminateRequest
	6, // 8: jobservice.JobService.EstablishConnection:output_type -> google.protobuf.Empty
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_jobproto_job_proto_init() }
func file_jobproto_job_proto_init() {
	if File_jobproto_job_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_jobproto_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminateRequest); i {
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
		file_jobproto_job_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminateConfirmation); i {
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
		file_jobproto_job_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_jobproto_job_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobResult); i {
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
		file_jobproto_job_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushInfo); i {
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
		file_jobproto_job_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connected); i {
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
			RawDescriptor: file_jobproto_job_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jobproto_job_proto_goTypes,
		DependencyIndexes: file_jobproto_job_proto_depIdxs,
		MessageInfos:      file_jobproto_job_proto_msgTypes,
	}.Build()
	File_jobproto_job_proto = out.File
	file_jobproto_job_proto_rawDesc = nil
	file_jobproto_job_proto_goTypes = nil
	file_jobproto_job_proto_depIdxs = nil
}
