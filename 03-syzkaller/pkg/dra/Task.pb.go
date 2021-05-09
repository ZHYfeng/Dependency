// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.11.0
// source: Task.proto

package dra

import (
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

type TaskStatus int32

const (
	TaskStatus_untested                  TaskStatus = 0
	TaskStatus_testing                   TaskStatus = 1
	TaskStatus_not_find_input            TaskStatus = -5
	TaskStatus_not_find_write_address    TaskStatus = -4
	TaskStatus_not_find_write_input      TaskStatus = -3
	TaskStatus_unstable_write            TaskStatus = 11
	TaskStatus_stable_write              TaskStatus = 12
	TaskStatus_unstable_condition        TaskStatus = 13
	TaskStatus_stable_condition          TaskStatus = 14
	TaskStatus_unstable_insert_write     TaskStatus = 15
	TaskStatus_stable_insert_write       TaskStatus = 16
	TaskStatus_unstable_insert_condition TaskStatus = 17
	TaskStatus_stable_insert_condition   TaskStatus = 18
	TaskStatus_unstable                  TaskStatus = 19
	TaskStatus_tested                    TaskStatus = 21
	TaskStatus_covered                   TaskStatus = 22
	TaskStatus_recursive                 TaskStatus = 31
	TaskStatus_out                       TaskStatus = 32
)

// Enum value maps for TaskStatus.
var (
	TaskStatus_name = map[int32]string{
		0:  "untested",
		1:  "testing",
		-5: "not_find_input",
		-4: "not_find_write_address",
		-3: "not_find_write_input",
		11: "unstable_write",
		12: "stable_write",
		13: "unstable_condition",
		14: "stable_condition",
		15: "unstable_insert_write",
		16: "stable_insert_write",
		17: "unstable_insert_condition",
		18: "stable_insert_condition",
		19: "unstable",
		21: "tested",
		22: "covered",
		31: "recursive",
		32: "out",
	}
	TaskStatus_value = map[string]int32{
		"untested":                  0,
		"testing":                   1,
		"not_find_input":            -5,
		"not_find_write_address":    -4,
		"not_find_write_input":      -3,
		"unstable_write":            11,
		"stable_write":              12,
		"unstable_condition":        13,
		"stable_condition":          14,
		"unstable_insert_write":     15,
		"stable_insert_write":       16,
		"unstable_insert_condition": 17,
		"stable_insert_condition":   18,
		"unstable":                  19,
		"tested":                    21,
		"covered":                   22,
		"recursive":                 31,
		"out":                       32,
	}
)

func (x TaskStatus) Enum() *TaskStatus {
	p := new(TaskStatus)
	*p = x
	return p
}

func (x TaskStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_Task_proto_enumTypes[0].Descriptor()
}

func (TaskStatus) Type() protoreflect.EnumType {
	return &file_Task_proto_enumTypes[0]
}

func (x TaskStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskStatus.Descriptor instead.
func (TaskStatus) EnumDescriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{0}
}

type WriteStatementKind int32

const (
	WriteStatementKind_WriteStatementConstant             WriteStatementKind = 0
	WriteStatementKind_WriteStatementNonconstant          WriteStatementKind = 1
	WriteStatementKind_WriteStatementDependencyRelated    WriteStatementKind = 2
	WriteStatementKind_WriteStatementNotDependencyRelated WriteStatementKind = 3
)

// Enum value maps for WriteStatementKind.
var (
	WriteStatementKind_name = map[int32]string{
		0: "WriteStatementConstant",
		1: "WriteStatementNonconstant",
		2: "WriteStatementDependencyRelated",
		3: "WriteStatementNotDependencyRelated",
	}
	WriteStatementKind_value = map[string]int32{
		"WriteStatementConstant":             0,
		"WriteStatementNonconstant":          1,
		"WriteStatementDependencyRelated":    2,
		"WriteStatementNotDependencyRelated": 3,
	}
)

func (x WriteStatementKind) Enum() *WriteStatementKind {
	p := new(WriteStatementKind)
	*p = x
	return p
}

func (x WriteStatementKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WriteStatementKind) Descriptor() protoreflect.EnumDescriptor {
	return file_Task_proto_enumTypes[1].Descriptor()
}

func (WriteStatementKind) Type() protoreflect.EnumType {
	return &file_Task_proto_enumTypes[1]
}

func (x WriteStatementKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WriteStatementKind.Descriptor instead.
func (WriteStatementKind) EnumDescriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{1}
}

type TaskKind int32

const (
	TaskKind_Boot   TaskKind = 0
	TaskKind_High   TaskKind = 1
	TaskKind_Ckeck  TaskKind = 3
	TaskKind_Normal TaskKind = 5
)

// Enum value maps for TaskKind.
var (
	TaskKind_name = map[int32]string{
		0: "Boot",
		1: "High",
		3: "Ckeck",
		5: "Normal",
	}
	TaskKind_value = map[string]int32{
		"Boot":   0,
		"High":   1,
		"Ckeck":  3,
		"Normal": 5,
	}
)

func (x TaskKind) Enum() *TaskKind {
	p := new(TaskKind)
	*p = x
	return p
}

func (x TaskKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskKind) Descriptor() protoreflect.EnumDescriptor {
	return file_Task_proto_enumTypes[2].Descriptor()
}

func (TaskKind) Type() protoreflect.EnumType {
	return &file_Task_proto_enumTypes[2]
}

func (x TaskKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskKind.Descriptor instead.
func (TaskKind) EnumDescriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{2}
}

type RunTimeData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Priority                uint32     `protobuf:"varint,1,opt,name=priority,proto3" json:"priority,omitempty"`
	WriteAddress            uint32     `protobuf:"varint,2,opt,name=write_address,json=writeAddress,proto3" json:"write_address,omitempty"`
	ConditionAddress        uint32     `protobuf:"varint,3,opt,name=condition_address,json=conditionAddress,proto3" json:"condition_address,omitempty"`
	Address                 uint32     `protobuf:"varint,4,opt,name=address,proto3" json:"address,omitempty"`
	RightBranchAddress      []uint32   `protobuf:"varint,6,rep,packed,name=right_branch_address,json=rightBranchAddress,proto3" json:"right_branch_address,omitempty"`
	TaskStatus              TaskStatus `protobuf:"varint,10,opt,name=task_status,json=taskStatus,proto3,enum=dra.TaskStatus" json:"task_status,omitempty"`
	Program                 []byte     `protobuf:"bytes,12,opt,name=program,proto3" json:"program,omitempty"`
	Idx                     uint32     `protobuf:"varint,13,opt,name=idx,proto3" json:"idx,omitempty"`
	RecursiveCount          uint32     `protobuf:"varint,14,opt,name=recursive_count,json=recursiveCount,proto3" json:"recursive_count,omitempty"`
	CheckWrite              bool       `protobuf:"varint,20,opt,name=checkWrite,proto3" json:"checkWrite,omitempty"`
	CheckCondition          bool       `protobuf:"varint,21,opt,name=checkCondition,proto3" json:"checkCondition,omitempty"`
	CheckAddress            bool       `protobuf:"varint,22,opt,name=checkAddress,proto3" json:"checkAddress,omitempty"`
	CheckRightBranchAddress []bool     `protobuf:"varint,23,rep,packed,name=checkRightBranchAddress,proto3" json:"checkRightBranchAddress,omitempty"` //    map<uint32, uint32> right_branch_address = 16;
}

func (x *RunTimeData) Reset() {
	*x = RunTimeData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Task_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunTimeData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunTimeData) ProtoMessage() {}

func (x *RunTimeData) ProtoReflect() protoreflect.Message {
	mi := &file_Task_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunTimeData.ProtoReflect.Descriptor instead.
func (*RunTimeData) Descriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{0}
}

func (x *RunTimeData) GetPriority() uint32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *RunTimeData) GetWriteAddress() uint32 {
	if x != nil {
		return x.WriteAddress
	}
	return 0
}

func (x *RunTimeData) GetConditionAddress() uint32 {
	if x != nil {
		return x.ConditionAddress
	}
	return 0
}

func (x *RunTimeData) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *RunTimeData) GetRightBranchAddress() []uint32 {
	if x != nil {
		return x.RightBranchAddress
	}
	return nil
}

func (x *RunTimeData) GetTaskStatus() TaskStatus {
	if x != nil {
		return x.TaskStatus
	}
	return TaskStatus_untested
}

func (x *RunTimeData) GetProgram() []byte {
	if x != nil {
		return x.Program
	}
	return nil
}

func (x *RunTimeData) GetIdx() uint32 {
	if x != nil {
		return x.Idx
	}
	return 0
}

func (x *RunTimeData) GetRecursiveCount() uint32 {
	if x != nil {
		return x.RecursiveCount
	}
	return 0
}

func (x *RunTimeData) GetCheckWrite() bool {
	if x != nil {
		return x.CheckWrite
	}
	return false
}

func (x *RunTimeData) GetCheckCondition() bool {
	if x != nil {
		return x.CheckCondition
	}
	return false
}

func (x *RunTimeData) GetCheckAddress() bool {
	if x != nil {
		return x.CheckAddress
	}
	return false
}

func (x *RunTimeData) GetCheckRightBranchAddress() []bool {
	if x != nil {
		return x.CheckRightBranchAddress
	}
	return nil
}

type TaskRunTimeData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash             string                  `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Program          []byte                  `protobuf:"bytes,2,opt,name=program,proto3" json:"program,omitempty"`
	WriteIdx         uint32                  `protobuf:"varint,5,opt,name=write_idx,json=writeIdx,proto3" json:"write_idx,omitempty"`
	ConditionIdx     uint32                  `protobuf:"varint,6,opt,name=condition_idx,json=conditionIdx,proto3" json:"condition_idx,omitempty"`
	Check            bool                    `protobuf:"varint,10,opt,name=check,proto3" json:"check,omitempty"`
	UncoveredAddress map[uint32]*RunTimeData `protobuf:"bytes,21,rep,name=uncovered_address,json=uncoveredAddress,proto3" json:"uncovered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CoveredAddress   map[uint32]*RunTimeData `protobuf:"bytes,23,rep,name=covered_address,json=coveredAddress,proto3" json:"covered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TaskRunTimeData) Reset() {
	*x = TaskRunTimeData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Task_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskRunTimeData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskRunTimeData) ProtoMessage() {}

func (x *TaskRunTimeData) ProtoReflect() protoreflect.Message {
	mi := &file_Task_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskRunTimeData.ProtoReflect.Descriptor instead.
func (*TaskRunTimeData) Descriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{1}
}

func (x *TaskRunTimeData) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TaskRunTimeData) GetProgram() []byte {
	if x != nil {
		return x.Program
	}
	return nil
}

func (x *TaskRunTimeData) GetWriteIdx() uint32 {
	if x != nil {
		return x.WriteIdx
	}
	return 0
}

func (x *TaskRunTimeData) GetConditionIdx() uint32 {
	if x != nil {
		return x.ConditionIdx
	}
	return 0
}

func (x *TaskRunTimeData) GetCheck() bool {
	if x != nil {
		return x.Check
	}
	return false
}

func (x *TaskRunTimeData) GetUncoveredAddress() map[uint32]*RunTimeData {
	if x != nil {
		return x.UncoveredAddress
	}
	return nil
}

func (x *TaskRunTimeData) GetCoveredAddress() map[uint32]*RunTimeData {
	if x != nil {
		return x.CoveredAddress
	}
	return nil
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sig          string             `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	Index        uint32             `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Program      []byte             `protobuf:"bytes,3,opt,name=program,proto3" json:"program,omitempty"`
	Kind         WriteStatementKind `protobuf:"varint,4,opt,name=kind,proto3,enum=dra.WriteStatementKind" json:"kind,omitempty"`
	Priority     int32              `protobuf:"varint,5,opt,name=priority,proto3" json:"priority,omitempty"`
	Hash         string             `protobuf:"bytes,6,opt,name=hash,proto3" json:"hash,omitempty"`
	Count        uint32             `protobuf:"varint,7,opt,name=count,proto3" json:"count,omitempty"`
	WriteSig     string             `protobuf:"bytes,11,opt,name=write_sig,json=writeSig,proto3" json:"write_sig,omitempty"`
	WriteIndex   uint32             `protobuf:"varint,12,opt,name=write_index,json=writeIndex,proto3" json:"write_index,omitempty"`
	WriteProgram []byte             `protobuf:"bytes,13,opt,name=write_program,json=writeProgram,proto3" json:"write_program,omitempty"`
	TaskStatus   TaskStatus         `protobuf:"varint,24,opt,name=task_status,json=taskStatus,proto3,enum=dra.TaskStatus" json:"task_status,omitempty"`
	Check        bool               `protobuf:"varint,25,opt,name=check,proto3" json:"check,omitempty"`
	// uncovered address, priority
	UncoveredAddress map[uint32]*RunTimeData `protobuf:"bytes,21,rep,name=uncovered_address,json=uncoveredAddress,proto3" json:"uncovered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CoveredAddress   map[uint32]*RunTimeData `protobuf:"bytes,23,rep,name=covered_address,json=coveredAddress,proto3" json:"covered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TaskRunTimeData  []*TaskRunTimeData      `protobuf:"bytes,31,rep,name=task_run_time_data,json=taskRunTimeData,proto3" json:"task_run_time_data,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Task_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_Task_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetSig() string {
	if x != nil {
		return x.Sig
	}
	return ""
}

func (x *Task) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Task) GetProgram() []byte {
	if x != nil {
		return x.Program
	}
	return nil
}

func (x *Task) GetKind() WriteStatementKind {
	if x != nil {
		return x.Kind
	}
	return WriteStatementKind_WriteStatementConstant
}

func (x *Task) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *Task) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Task) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Task) GetWriteSig() string {
	if x != nil {
		return x.WriteSig
	}
	return ""
}

func (x *Task) GetWriteIndex() uint32 {
	if x != nil {
		return x.WriteIndex
	}
	return 0
}

func (x *Task) GetWriteProgram() []byte {
	if x != nil {
		return x.WriteProgram
	}
	return nil
}

func (x *Task) GetTaskStatus() TaskStatus {
	if x != nil {
		return x.TaskStatus
	}
	return TaskStatus_untested
}

func (x *Task) GetCheck() bool {
	if x != nil {
		return x.Check
	}
	return false
}

func (x *Task) GetUncoveredAddress() map[uint32]*RunTimeData {
	if x != nil {
		return x.UncoveredAddress
	}
	return nil
}

func (x *Task) GetCoveredAddress() map[uint32]*RunTimeData {
	if x != nil {
		return x.CoveredAddress
	}
	return nil
}

func (x *Task) GetTaskRunTimeData() []*TaskRunTimeData {
	if x != nil {
		return x.TaskRunTimeData
	}
	return nil
}

type Tasks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string           `protobuf:"bytes,41,opt,name=name,proto3" json:"name,omitempty"`
	Kind      TaskKind         `protobuf:"varint,42,opt,name=kind,proto3,enum=dra.TaskKind" json:"kind,omitempty"`
	TaskMap   map[string]*Task `protobuf:"bytes,1,rep,name=task_map,json=taskMap,proto3" json:"task_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TaskArray []*Task          `protobuf:"bytes,2,rep,name=task_array,json=taskArray,proto3" json:"task_array,omitempty"`
}

func (x *Tasks) Reset() {
	*x = Tasks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Task_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tasks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tasks) ProtoMessage() {}

func (x *Tasks) ProtoReflect() protoreflect.Message {
	mi := &file_Task_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tasks.ProtoReflect.Descriptor instead.
func (*Tasks) Descriptor() ([]byte, []int) {
	return file_Task_proto_rawDescGZIP(), []int{3}
}

func (x *Tasks) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tasks) GetKind() TaskKind {
	if x != nil {
		return x.Kind
	}
	return TaskKind_Boot
}

func (x *Tasks) GetTaskMap() map[string]*Task {
	if x != nil {
		return x.TaskMap
	}
	return nil
}

func (x *Tasks) GetTaskArray() []*Task {
	if x != nil {
		return x.TaskArray
	}
	return nil
}

var File_Task_proto protoreflect.FileDescriptor

var file_Task_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x64, 0x72,
	0x61, 0x22, 0xf4, 0x03, 0x0a, 0x0b, 0x72, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x23, 0x0a,
	0x0d, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x77, 0x72, 0x69, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x72, 0x69, 0x67,
	0x68, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x12, 0x72, 0x69, 0x67, 0x68, 0x74, 0x42, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x30, 0x0a, 0x0b, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0f, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x78, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x69, 0x64, 0x78, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x63,
	0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x18, 0x14, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x15, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x16, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0c, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x38,
	0x0a, 0x17, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x69, 0x67, 0x68, 0x74, 0x42, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x17, 0x20, 0x03, 0x28, 0x08, 0x52,
	0x17, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x69, 0x67, 0x68, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63,
	0x68, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0xef, 0x03, 0x0a, 0x0f, 0x54, 0x61, 0x73,
	0x6b, 0x52, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x77,
	0x72, 0x69, 0x74, 0x65, 0x49, 0x64, 0x78, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x78, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c,
	0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x78, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x12, 0x57, 0x0a, 0x11, 0x75, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e,
	0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x2e, 0x55, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x75, 0x6e, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x51, 0x0a, 0x0f, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x17,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x43, 0x6f, 0x76, 0x65, 0x72,
	0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x1a, 0x55,
	0x0a, 0x15, 0x55, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x72,
	0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x53, 0x0a, 0x13, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x64, 0x72, 0x61, 0x2e, 0x72, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xeb, 0x05, 0x0a, 0x04, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x72,
	0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x2b, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x5f, 0x73, 0x69, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x72, 0x69,
	0x74, 0x65, 0x53, 0x69, 0x67, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x23, 0x0a, 0x0d, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f,
	0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x77,
	0x72, 0x69, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x30, 0x0a, 0x0b, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x18, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0f, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x19, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x12, 0x4c, 0x0a, 0x11, 0x75, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x55, 0x6e, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x10, 0x75, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x46, 0x0a, 0x0f, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x17, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x72, 0x61,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x41, 0x0a, 0x12, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x1f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0f, 0x74, 0x61, 0x73,
	0x6b, 0x52, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x55, 0x0a, 0x15,
	0x55, 0x6e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x72, 0x75, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x1a, 0x53, 0x0a, 0x13, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x72,
	0x61, 0x2e, 0x72, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xe3, 0x01, 0x0a, 0x05, 0x54, 0x61, 0x73,
	0x6b, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x29, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x2a,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4b,
	0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x72,
	0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x61, 0x70, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x4d, 0x61, 0x70, 0x12, 0x28, 0x0a,
	0x0a, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x61, 0x72, 0x72, 0x61, 0x79, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x09, 0x74, 0x61,
	0x73, 0x6b, 0x41, 0x72, 0x72, 0x61, 0x79, 0x1a, 0x45, 0x0a, 0x0c, 0x54, 0x61, 0x73, 0x6b, 0x4d,
	0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x64, 0x72, 0x61, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x8f,
	0x03, 0x0a, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0c, 0x0a,
	0x08, 0x75, 0x6e, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x74,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x0e, 0x6e, 0x6f, 0x74, 0x5f,
	0x66, 0x69, 0x6e, 0x64, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x10, 0xfb, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x01, 0x12, 0x23, 0x0a, 0x16, 0x6e, 0x6f, 0x74, 0x5f, 0x66, 0x69, 0x6e,
	0x64, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x10,
	0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x12, 0x21, 0x0a, 0x14, 0x6e, 0x6f,
	0x74, 0x5f, 0x66, 0x69, 0x6e, 0x64, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x12, 0x12, 0x0a,
	0x0e, 0x75, 0x6e, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x10,
	0x0b, 0x12, 0x10, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x10, 0x0c, 0x12, 0x16, 0x0a, 0x12, 0x75, 0x6e, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x0d, 0x12, 0x14, 0x0a, 0x10, 0x73,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x10,
	0x0e, 0x12, 0x19, 0x0a, 0x15, 0x75, 0x6e, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x69, 0x6e,
	0x73, 0x65, 0x72, 0x74, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x10, 0x0f, 0x12, 0x17, 0x0a, 0x13,
	0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x77, 0x72,
	0x69, 0x74, 0x65, 0x10, 0x10, 0x12, 0x1d, 0x0a, 0x19, 0x75, 0x6e, 0x73, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x10, 0x11, 0x12, 0x1b, 0x0a, 0x17, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x69,
	0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x10,
	0x12, 0x12, 0x0c, 0x0a, 0x08, 0x75, 0x6e, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x13, 0x12,
	0x0a, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x10, 0x15, 0x12, 0x0b, 0x0a, 0x07, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x65, 0x64, 0x10, 0x16, 0x12, 0x0d, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x75,
	0x72, 0x73, 0x69, 0x76, 0x65, 0x10, 0x1f, 0x12, 0x07, 0x0a, 0x03, 0x6f, 0x75, 0x74, 0x10, 0x20,
	0x2a, 0x9c, 0x01, 0x0a, 0x12, 0x57, 0x72, 0x69, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x16, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x57, 0x72, 0x69, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x6e, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x10, 0x01, 0x12, 0x23, 0x0a, 0x1f, 0x57, 0x72, 0x69, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65,
	0x6c, 0x61, 0x74, 0x65, 0x64, 0x10, 0x02, 0x12, 0x26, 0x0a, 0x22, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x74, 0x44, 0x65, 0x70, 0x65,
	0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x10, 0x03, 0x2a,
	0x35, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x08, 0x0a, 0x04, 0x42,
	0x6f, 0x6f, 0x74, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x69, 0x67, 0x68, 0x10, 0x01, 0x12,
	0x09, 0x0a, 0x05, 0x43, 0x6b, 0x65, 0x63, 0x6b, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x6f,
	0x72, 0x6d, 0x61, 0x6c, 0x10, 0x05, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Task_proto_rawDescOnce sync.Once
	file_Task_proto_rawDescData = file_Task_proto_rawDesc
)

func file_Task_proto_rawDescGZIP() []byte {
	file_Task_proto_rawDescOnce.Do(func() {
		file_Task_proto_rawDescData = protoimpl.X.CompressGZIP(file_Task_proto_rawDescData)
	})
	return file_Task_proto_rawDescData
}

var file_Task_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_Task_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_Task_proto_goTypes = []interface{}{
	(TaskStatus)(0),         // 0: dra.taskStatus
	(WriteStatementKind)(0), // 1: dra.WriteStatementKind
	(TaskKind)(0),           // 2: dra.TaskKind
	(*RunTimeData)(nil),     // 3: dra.runTimeData
	(*TaskRunTimeData)(nil), // 4: dra.TaskRunTimeData
	(*Task)(nil),            // 5: dra.Task
	(*Tasks)(nil),           // 6: dra.Tasks
	nil,                     // 7: dra.TaskRunTimeData.UncoveredAddressEntry
	nil,                     // 8: dra.TaskRunTimeData.CoveredAddressEntry
	nil,                     // 9: dra.Task.UncoveredAddressEntry
	nil,                     // 10: dra.Task.CoveredAddressEntry
	nil,                     // 11: dra.Tasks.TaskMapEntry
}
var file_Task_proto_depIdxs = []int32{
	0,  // 0: dra.runTimeData.task_status:type_name -> dra.taskStatus
	7,  // 1: dra.TaskRunTimeData.uncovered_address:type_name -> dra.TaskRunTimeData.UncoveredAddressEntry
	8,  // 2: dra.TaskRunTimeData.covered_address:type_name -> dra.TaskRunTimeData.CoveredAddressEntry
	1,  // 3: dra.Task.kind:type_name -> dra.WriteStatementKind
	0,  // 4: dra.Task.task_status:type_name -> dra.taskStatus
	9,  // 5: dra.Task.uncovered_address:type_name -> dra.Task.UncoveredAddressEntry
	10, // 6: dra.Task.covered_address:type_name -> dra.Task.CoveredAddressEntry
	4,  // 7: dra.Task.task_run_time_data:type_name -> dra.TaskRunTimeData
	2,  // 8: dra.Tasks.kind:type_name -> dra.TaskKind
	11, // 9: dra.Tasks.task_map:type_name -> dra.Tasks.TaskMapEntry
	5,  // 10: dra.Tasks.task_array:type_name -> dra.Task
	3,  // 11: dra.TaskRunTimeData.UncoveredAddressEntry.value:type_name -> dra.runTimeData
	3,  // 12: dra.TaskRunTimeData.CoveredAddressEntry.value:type_name -> dra.runTimeData
	3,  // 13: dra.Task.UncoveredAddressEntry.value:type_name -> dra.runTimeData
	3,  // 14: dra.Task.CoveredAddressEntry.value:type_name -> dra.runTimeData
	5,  // 15: dra.Tasks.TaskMapEntry.value:type_name -> dra.Task
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_Task_proto_init() }
func file_Task_proto_init() {
	if File_Task_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Task_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunTimeData); i {
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
		file_Task_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskRunTimeData); i {
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
		file_Task_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_Task_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tasks); i {
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
			RawDescriptor: file_Task_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Task_proto_goTypes,
		DependencyIndexes: file_Task_proto_depIdxs,
		EnumInfos:         file_Task_proto_enumTypes,
		MessageInfos:      file_Task_proto_msgTypes,
	}.Build()
	File_Task_proto = out.File
	file_Task_proto_rawDesc = nil
	file_Task_proto_goTypes = nil
	file_Task_proto_depIdxs = nil
}
