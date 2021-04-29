// Code generated by protoc-gen-go. DO NOT EDIT.
// source: DependencyRPC.proto

package dra

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DataDependency struct {
	Input                map[string]*Input            `protobuf:"bytes,1,rep,name=input,proto3" json:"input,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UncoveredAddress     map[uint32]*UncoveredAddress `protobuf:"bytes,4,rep,name=uncovered_address,json=uncoveredAddress,proto3" json:"uncovered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	WriteAddress         map[uint32]*WriteAddress     `protobuf:"bytes,5,rep,name=write_address,json=writeAddress,proto3" json:"write_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	OtherInput           map[string]*Input            `protobuf:"bytes,11,rep,name=other_input,json=otherInput,proto3" json:"other_input,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *DataDependency) Reset()         { *m = DataDependency{} }
func (m *DataDependency) String() string { return proto.CompactTextString(m) }
func (*DataDependency) ProtoMessage()    {}
func (*DataDependency) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{0}
}

func (m *DataDependency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataDependency.Unmarshal(m, b)
}
func (m *DataDependency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataDependency.Marshal(b, m, deterministic)
}
func (m *DataDependency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataDependency.Merge(m, src)
}
func (m *DataDependency) XXX_Size() int {
	return xxx_messageInfo_DataDependency.Size(m)
}
func (m *DataDependency) XXX_DiscardUnknown() {
	xxx_messageInfo_DataDependency.DiscardUnknown(m)
}

var xxx_messageInfo_DataDependency proto.InternalMessageInfo

func (m *DataDependency) GetInput() map[string]*Input {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *DataDependency) GetUncoveredAddress() map[uint32]*UncoveredAddress {
	if m != nil {
		return m.UncoveredAddress
	}
	return nil
}

func (m *DataDependency) GetWriteAddress() map[uint32]*WriteAddress {
	if m != nil {
		return m.WriteAddress
	}
	return nil
}

func (m *DataDependency) GetOtherInput() map[string]*Input {
	if m != nil {
		return m.OtherInput
	}
	return nil
}

type DataResult struct {
	CoveredAddress       map[uint32]*UncoveredAddress `protobuf:"bytes,2,rep,name=covered_address,json=coveredAddress,proto3" json:"covered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	FileOperations       map[string]*FileOperations   `protobuf:"bytes,6,rep,name=file_operations,json=fileOperations,proto3" json:"file_operations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *DataResult) Reset()         { *m = DataResult{} }
func (m *DataResult) String() string { return proto.CompactTextString(m) }
func (*DataResult) ProtoMessage()    {}
func (*DataResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{1}
}

func (m *DataResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataResult.Unmarshal(m, b)
}
func (m *DataResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataResult.Marshal(b, m, deterministic)
}
func (m *DataResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataResult.Merge(m, src)
}
func (m *DataResult) XXX_Size() int {
	return xxx_messageInfo_DataResult.Size(m)
}
func (m *DataResult) XXX_DiscardUnknown() {
	xxx_messageInfo_DataResult.DiscardUnknown(m)
}

var xxx_messageInfo_DataResult proto.InternalMessageInfo

func (m *DataResult) GetCoveredAddress() map[uint32]*UncoveredAddress {
	if m != nil {
		return m.CoveredAddress
	}
	return nil
}

func (m *DataResult) GetFileOperations() map[string]*FileOperations {
	if m != nil {
		return m.FileOperations
	}
	return nil
}

type DataRunTime struct {
	Tasks                *Tasks   `protobuf:"bytes,11,opt,name=tasks,proto3" json:"tasks,omitempty"`
	Return               *Tasks   `protobuf:"bytes,12,opt,name=return,proto3" json:"return,omitempty"`
	HighTask             *Tasks   `protobuf:"bytes,13,opt,name=high_task,json=highTask,proto3" json:"high_task,omitempty"`
	BootTask             *Tasks   `protobuf:"bytes,20,opt,name=boot_task,json=bootTask,proto3" json:"boot_task,omitempty"`
	ReturnBoot           *Tasks   `protobuf:"bytes,21,opt,name=return_boot,json=returnBoot,proto3" json:"return_boot,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataRunTime) Reset()         { *m = DataRunTime{} }
func (m *DataRunTime) String() string { return proto.CompactTextString(m) }
func (*DataRunTime) ProtoMessage()    {}
func (*DataRunTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{2}
}

func (m *DataRunTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataRunTime.Unmarshal(m, b)
}
func (m *DataRunTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataRunTime.Marshal(b, m, deterministic)
}
func (m *DataRunTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataRunTime.Merge(m, src)
}
func (m *DataRunTime) XXX_Size() int {
	return xxx_messageInfo_DataRunTime.Size(m)
}
func (m *DataRunTime) XXX_DiscardUnknown() {
	xxx_messageInfo_DataRunTime.DiscardUnknown(m)
}

var xxx_messageInfo_DataRunTime proto.InternalMessageInfo

func (m *DataRunTime) GetTasks() *Tasks {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func (m *DataRunTime) GetReturn() *Tasks {
	if m != nil {
		return m.Return
	}
	return nil
}

func (m *DataRunTime) GetHighTask() *Tasks {
	if m != nil {
		return m.HighTask
	}
	return nil
}

func (m *DataRunTime) GetBootTask() *Tasks {
	if m != nil {
		return m.BootTask
	}
	return nil
}

func (m *DataRunTime) GetReturnBoot() *Tasks {
	if m != nil {
		return m.ReturnBoot
	}
	return nil
}

func init() {
	proto.RegisterType((*DataDependency)(nil), "dra.DataDependency")
	proto.RegisterMapType((map[string]*Input)(nil), "dra.DataDependency.InputEntry")
	proto.RegisterMapType((map[string]*Input)(nil), "dra.DataDependency.OtherInputEntry")
	proto.RegisterMapType((map[uint32]*UncoveredAddress)(nil), "dra.DataDependency.UncoveredAddressEntry")
	proto.RegisterMapType((map[uint32]*WriteAddress)(nil), "dra.DataDependency.WriteAddressEntry")
	proto.RegisterType((*DataResult)(nil), "dra.DataResult")
	proto.RegisterMapType((map[uint32]*UncoveredAddress)(nil), "dra.DataResult.CoveredAddressEntry")
	proto.RegisterMapType((map[string]*FileOperations)(nil), "dra.DataResult.FileOperationsEntry")
	proto.RegisterType((*DataRunTime)(nil), "dra.DataRunTime")
}

func init() { proto.RegisterFile("DependencyRPC.proto", fileDescriptor_db4d5fd3d0a7c985) }

var fileDescriptor_db4d5fd3d0a7c985 = []byte{
	// 756 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x96, 0x5b, 0x6f, 0x1a, 0x39,
	0x14, 0xc7, 0x03, 0x09, 0x24, 0x39, 0xc3, 0xd5, 0x24, 0xd2, 0x88, 0x95, 0x56, 0x88, 0xec, 0x2a,
	0x44, 0xd1, 0xb2, 0xab, 0xec, 0x66, 0x5b, 0xb5, 0x4f, 0x05, 0x52, 0x94, 0x2a, 0x4d, 0xaa, 0xc9,
	0xa5, 0x55, 0x5e, 0xd0, 0xc0, 0x1c, 0x92, 0x11, 0x64, 0x8c, 0xc6, 0x9e, 0x44, 0x7c, 0xcd, 0x3e,
	0xf6, 0x03, 0xf4, 0xbd, 0xdf, 0xa0, 0xb2, 0x3d, 0x05, 0xcf, 0x45, 0x44, 0xaa, 0xd4, 0x37, 0x73,
	0xce, 0xff, 0xfc, 0x3c, 0xe7, 0x62, 0x1b, 0xa8, 0xf5, 0x70, 0x86, 0x9e, 0x83, 0xde, 0x68, 0x6e,
	0x7d, 0xe8, 0xb6, 0x67, 0x3e, 0xe5, 0x94, 0xac, 0x3b, 0xbe, 0x5d, 0x87, 0x8e, 0xcd, 0x50, 0x19,
	0xea, 0x70, 0x65, 0xb3, 0x49, 0xb8, 0x36, 0x4e, 0xbd, 0x59, 0xc0, 0xc3, 0x1f, 0xe5, 0x4b, 0x6e,
	0x73, 0x97, 0x71, 0x77, 0x14, 0x1a, 0x2a, 0x4b, 0x9e, 0xb2, 0x34, 0xbf, 0x6d, 0x40, 0xa9, 0x67,
	0x73, 0x7b, 0xe9, 0x20, 0xff, 0x41, 0xce, 0x15, 0x10, 0x33, 0xd3, 0x58, 0x6f, 0x19, 0x47, 0xbf,
	0xb7, 0x1d, 0xdf, 0x6e, 0x47, 0x35, 0x6d, 0xb9, 0xcb, 0x89, 0xc7, 0xfd, 0xb9, 0xa5, 0xc4, 0xe4,
	0x06, 0xaa, 0x81, 0x37, 0xa2, 0x8f, 0xe8, 0xa3, 0x33, 0xb0, 0x1d, 0xc7, 0x47, 0xc6, 0xcc, 0x0d,
	0x49, 0x38, 0x48, 0x23, 0x5c, 0xff, 0x10, 0xbf, 0x51, 0x5a, 0x05, 0xab, 0x04, 0x31, 0x33, 0x79,
	0x07, 0xc5, 0x27, 0xdf, 0xe5, 0xb8, 0x60, 0xe6, 0x24, 0xf3, 0xcf, 0x34, 0xe6, 0x47, 0x21, 0x8c,
	0xf0, 0x0a, 0x4f, 0x9a, 0x89, 0xf4, 0xc0, 0xa0, 0xfc, 0x1e, 0xfd, 0x81, 0xca, 0xcf, 0x90, 0xa4,
	0xbd, 0x34, 0xd2, 0x85, 0x90, 0x69, 0x49, 0x02, 0x5d, 0x18, 0xea, 0x3d, 0x80, 0xa5, 0x87, 0x54,
	0x60, 0x7d, 0x82, 0x73, 0x33, 0xd3, 0xc8, 0xb4, 0xb6, 0x2d, 0xb1, 0x24, 0x0d, 0xc8, 0x3d, 0xda,
	0xd3, 0x00, 0xcd, 0x6c, 0x23, 0xd3, 0x32, 0x8e, 0x40, 0xf2, 0x65, 0x84, 0xa5, 0x1c, 0xaf, 0xb2,
	0x2f, 0x33, 0xf5, 0x5b, 0xd8, 0x4d, 0x2d, 0x81, 0x0e, 0x2c, 0x2a, 0xe0, 0x61, 0x14, 0xb8, 0x2b,
	0x81, 0xf1, 0x60, 0x9d, 0x6d, 0x41, 0x35, 0x51, 0x8a, 0x14, 0xee, 0x7e, 0x94, 0x5b, 0x95, 0x5c,
	0x3d, 0x50, 0x67, 0x9e, 0x42, 0x39, 0x56, 0x94, 0x9f, 0x4d, 0xbd, 0xf9, 0x25, 0x0b, 0x20, 0xea,
	0x6d, 0x21, 0x0b, 0xa6, 0x9c, 0x9c, 0x41, 0x39, 0x3e, 0x37, 0xd9, 0x58, 0x67, 0x94, 0xb2, 0xdd,
	0x4d, 0x99, 0x98, 0x52, 0x6c, 0x5e, 0xce, 0xa0, 0x3c, 0x76, 0xa7, 0x38, 0xa0, 0x33, 0xf4, 0x6d,
	0xee, 0x52, 0x8f, 0x99, 0xf9, 0x74, 0xda, 0x5b, 0x77, 0x8a, 0x17, 0x0b, 0x55, 0x48, 0x1b, 0x47,
	0x8c, 0xf5, 0x4f, 0x50, 0xeb, 0xfe, 0x9a, 0x1e, 0xdd, 0x40, 0x2d, 0xe5, 0x03, 0x52, 0x6a, 0x7a,
	0x10, 0x25, 0xd7, 0x24, 0x39, 0x1a, 0xaa, 0x17, 0xf7, 0x73, 0x06, 0x0c, 0x99, 0x64, 0xe0, 0x5d,
	0xb9, 0x0f, 0x28, 0x5a, 0xc2, 0x6d, 0x36, 0x61, 0xa6, 0xa1, 0xb5, 0x44, 0x5c, 0x18, 0xcc, 0x52,
	0x0e, 0xd2, 0x84, 0xbc, 0x8f, 0x3c, 0xf0, 0x3d, 0xb3, 0x90, 0x90, 0x84, 0x1e, 0xb2, 0x0f, 0xdb,
	0xf7, 0xee, 0xdd, 0xfd, 0x40, 0x44, 0x98, 0xc5, 0x84, 0x6c, 0x4b, 0x38, 0xc5, 0x52, 0x08, 0x87,
	0x94, 0x72, 0x25, 0xdc, 0x49, 0x0a, 0x85, 0x53, 0x0a, 0x0f, 0xc1, 0x50, 0xec, 0x81, 0x30, 0x99,
	0xbb, 0x09, 0x29, 0x28, 0x77, 0x87, 0x52, 0x7e, 0xf4, 0x35, 0x0f, 0xc5, 0xc8, 0x55, 0x48, 0x5a,
	0x50, 0xe8, 0x23, 0xbf, 0x79, 0x7f, 0x31, 0x1e, 0x33, 0xe4, 0x8c, 0xa8, 0xc8, 0x93, 0x87, 0x19,
	0x9f, 0xd7, 0xb5, 0x75, 0x73, 0x8d, 0xfc, 0x03, 0x3b, 0x97, 0xe8, 0x39, 0xe7, 0xc1, 0xc3, 0x10,
	0xfd, 0x8e, 0xcd, 0xdc, 0x51, 0x67, 0x4a, 0x47, 0x93, 0x15, 0x11, 0x2f, 0xe0, 0xb7, 0xb4, 0x88,
	0x70, 0x10, 0x56, 0x04, 0xb6, 0xc0, 0xe8, 0x23, 0x3f, 0xc7, 0x27, 0x39, 0xf2, 0x11, 0xa1, 0xb1,
	0x3c, 0x0a, 0xac, 0xb9, 0x46, 0xfe, 0x86, 0x92, 0xd8, 0x42, 0xbb, 0x75, 0xcb, 0x6a, 0x3c, 0x17,
	0x86, 0x18, 0xfa, 0x2f, 0x99, 0x6f, 0x97, 0x7a, 0x8e, 0x2b, 0x5a, 0x1e, 0x61, 0xab, 0xd0, 0x85,
	0x4f, 0xf0, 0x8f, 0xa1, 0x22, 0xf8, 0xfa, 0x61, 0x26, 0xb5, 0xc4, 0xf9, 0x46, 0x16, 0xdb, 0x65,
	0x0f, 0x36, 0xbb, 0xd4, 0xf3, 0x70, 0xc4, 0x57, 0x64, 0xf9, 0x3f, 0x54, 0xfb, 0xc8, 0x63, 0x8f,
	0x86, 0x2e, 0xaf, 0xa5, 0xdc, 0xa8, 0xb2, 0x3a, 0x05, 0x59, 0xd6, 0x68, 0x79, 0xd4, 0x9d, 0x1a,
	0xdd, 0xe1, 0x0f, 0xd8, 0xea, 0xa3, 0x1c, 0x93, 0xb4, 0xc6, 0x4a, 0xbb, 0xe2, 0xf5, 0x91, 0x77,
	0xc2, 0x81, 0x5a, 0xa5, 0xdc, 0x07, 0xc3, 0x92, 0xc3, 0xa4, 0x0b, 0xe5, 0x3a, 0xb6, 0xf1, 0x01,
	0x14, 0xc5, 0x27, 0x0a, 0xe6, 0x73, 0xdf, 0x78, 0x0c, 0x55, 0x21, 0xbd, 0xf6, 0x18, 0xb7, 0x87,
	0x53, 0x54, 0x72, 0x12, 0x1e, 0x7b, 0xcd, 0x96, 0xac, 0xb0, 0x08, 0x3b, 0xa3, 0x77, 0x2b, 0xe7,
	0x68, 0x8b, 0xa1, 0xe7, 0x88, 0xd7, 0x9b, 0x94, 0xa4, 0x67, 0xf1, 0x90, 0x27, 0x71, 0x72, 0xe2,
	0x56, 0x8e, 0x65, 0x98, 0x95, 0x50, 0x3d, 0x93, 0x55, 0x67, 0xf3, 0x36, 0xd7, 0x7e, 0xed, 0xf8,
	0xf6, 0x30, 0x2f, 0xff, 0x1e, 0xfc, 0xfb, 0x3d, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x21, 0xc7, 0x38,
	0x82, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DependencyRPCClient is the client API for DependencyRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DependencyRPCClient interface {
	// DRA and syz-manager
	GetVMOffsets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	SendNumberBasicBlock(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	SendNumberBasicBlockCovered(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	GetNewInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Inputs, error)
	SendDependency(ctx context.Context, in *Dependency, opts ...grpc.CallOption) (*Empty, error)
	GetCondition(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Conditions, error)
	SendWriteAddress(ctx context.Context, in *WriteAddresses, opts ...grpc.CallOption) (*Empty, error)
	//syz-fuzzer and syz-manager
	Connect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	GetDataDependency(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DataDependency, error)
	SendNewInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error)
	GetTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Tasks, error)
	GetBootTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Tasks, error)
	ReturnTasks(ctx context.Context, in *Tasks, opts ...grpc.CallOption) (*Empty, error)
	SendBootInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error)
	SendUnstableInput(ctx context.Context, in *UnstableInput, opts ...grpc.CallOption) (*Empty, error)
	SendLog(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	SendStat(ctx context.Context, in *Statistic, opts ...grpc.CallOption) (*Empty, error)
	GetNeed(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	SendNeedInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error)
}

type dependencyRPCClient struct {
	cc *grpc.ClientConn
}

func NewDependencyRPCClient(cc *grpc.ClientConn) DependencyRPCClient {
	return &dependencyRPCClient{cc}
}

func (c *dependencyRPCClient) GetVMOffsets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetVMOffsets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendNumberBasicBlock(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendNumberBasicBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendNumberBasicBlockCovered(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendNumberBasicBlockCovered", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetNewInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Inputs, error) {
	out := new(Inputs)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetNewInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendDependency(ctx context.Context, in *Dependency, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendDependency", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetCondition(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Conditions, error) {
	out := new(Conditions)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetCondition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendWriteAddress(ctx context.Context, in *WriteAddresses, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendWriteAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) Connect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetDataDependency(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DataDependency, error) {
	out := new(DataDependency)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetDataDependency", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendNewInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendNewInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Tasks, error) {
	out := new(Tasks)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetBootTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Tasks, error) {
	out := new(Tasks)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetBootTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) ReturnTasks(ctx context.Context, in *Tasks, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/ReturnTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendBootInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendBootInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendUnstableInput(ctx context.Context, in *UnstableInput, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendUnstableInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendLog(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendStat(ctx context.Context, in *Statistic, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/sendStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetNeed(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetNeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendNeedInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendNeedInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DependencyRPCServer is the server API for DependencyRPC service.
type DependencyRPCServer interface {
	// DRA and syz-manager
	GetVMOffsets(context.Context, *Empty) (*Empty, error)
	SendNumberBasicBlock(context.Context, *Empty) (*Empty, error)
	SendNumberBasicBlockCovered(context.Context, *Empty) (*Empty, error)
	GetNewInput(context.Context, *Empty) (*Inputs, error)
	SendDependency(context.Context, *Dependency) (*Empty, error)
	GetCondition(context.Context, *Empty) (*Conditions, error)
	SendWriteAddress(context.Context, *WriteAddresses) (*Empty, error)
	//syz-fuzzer and syz-manager
	Connect(context.Context, *Empty) (*Empty, error)
	GetDataDependency(context.Context, *Empty) (*DataDependency, error)
	SendNewInput(context.Context, *Input) (*Empty, error)
	GetTasks(context.Context, *Empty) (*Tasks, error)
	GetBootTasks(context.Context, *Empty) (*Tasks, error)
	ReturnTasks(context.Context, *Tasks) (*Empty, error)
	SendBootInput(context.Context, *Input) (*Empty, error)
	SendUnstableInput(context.Context, *UnstableInput) (*Empty, error)
	SendLog(context.Context, *Empty) (*Empty, error)
	SendStat(context.Context, *Statistic) (*Empty, error)
	GetNeed(context.Context, *Empty) (*Empty, error)
	SendNeedInput(context.Context, *Input) (*Empty, error)
}

// UnimplementedDependencyRPCServer can be embedded to have forward compatible implementations.
type UnimplementedDependencyRPCServer struct {
}

func (*UnimplementedDependencyRPCServer) GetVMOffsets(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVMOffsets not implemented")
}
func (*UnimplementedDependencyRPCServer) SendNumberBasicBlock(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNumberBasicBlock not implemented")
}
func (*UnimplementedDependencyRPCServer) SendNumberBasicBlockCovered(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNumberBasicBlockCovered not implemented")
}
func (*UnimplementedDependencyRPCServer) GetNewInput(ctx context.Context, req *Empty) (*Inputs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewInput not implemented")
}
func (*UnimplementedDependencyRPCServer) SendDependency(ctx context.Context, req *Dependency) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDependency not implemented")
}
func (*UnimplementedDependencyRPCServer) GetCondition(ctx context.Context, req *Empty) (*Conditions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCondition not implemented")
}
func (*UnimplementedDependencyRPCServer) SendWriteAddress(ctx context.Context, req *WriteAddresses) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendWriteAddress not implemented")
}
func (*UnimplementedDependencyRPCServer) Connect(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (*UnimplementedDependencyRPCServer) GetDataDependency(ctx context.Context, req *Empty) (*DataDependency, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataDependency not implemented")
}
func (*UnimplementedDependencyRPCServer) SendNewInput(ctx context.Context, req *Input) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNewInput not implemented")
}
func (*UnimplementedDependencyRPCServer) GetTasks(ctx context.Context, req *Empty) (*Tasks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTasks not implemented")
}
func (*UnimplementedDependencyRPCServer) GetBootTasks(ctx context.Context, req *Empty) (*Tasks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBootTasks not implemented")
}
func (*UnimplementedDependencyRPCServer) ReturnTasks(ctx context.Context, req *Tasks) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnTasks not implemented")
}
func (*UnimplementedDependencyRPCServer) SendBootInput(ctx context.Context, req *Input) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBootInput not implemented")
}
func (*UnimplementedDependencyRPCServer) SendUnstableInput(ctx context.Context, req *UnstableInput) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUnstableInput not implemented")
}
func (*UnimplementedDependencyRPCServer) SendLog(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendLog not implemented")
}
func (*UnimplementedDependencyRPCServer) SendStat(ctx context.Context, req *Statistic) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendStat not implemented")
}
func (*UnimplementedDependencyRPCServer) GetNeed(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNeed not implemented")
}
func (*UnimplementedDependencyRPCServer) SendNeedInput(ctx context.Context, req *Input) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNeedInput not implemented")
}

func RegisterDependencyRPCServer(s *grpc.Server, srv DependencyRPCServer) {
	s.RegisterService(&_DependencyRPC_serviceDesc, srv)
}

func _DependencyRPC_GetVMOffsets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetVMOffsets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetVMOffsets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetVMOffsets(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendNumberBasicBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendNumberBasicBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendNumberBasicBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendNumberBasicBlock(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendNumberBasicBlockCovered_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendNumberBasicBlockCovered(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendNumberBasicBlockCovered",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendNumberBasicBlockCovered(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetNewInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetNewInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetNewInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetNewInput(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendDependency_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Dependency)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendDependency(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendDependency",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendDependency(ctx, req.(*Dependency))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetCondition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetCondition(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendWriteAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteAddresses)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendWriteAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendWriteAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendWriteAddress(ctx, req.(*WriteAddresses))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).Connect(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetDataDependency_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetDataDependency(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetDataDependency",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetDataDependency(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendNewInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendNewInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendNewInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendNewInput(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetTasks(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetBootTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetBootTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetBootTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetBootTasks(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_ReturnTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tasks)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).ReturnTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/ReturnTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).ReturnTasks(ctx, req.(*Tasks))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendBootInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendBootInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendBootInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendBootInput(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendUnstableInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnstableInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendUnstableInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendUnstableInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendUnstableInput(ctx, req.(*UnstableInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendLog(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Statistic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendStat(ctx, req.(*Statistic))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_GetNeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetNeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetNeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetNeed(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendNeedInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendNeedInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendNeedInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendNeedInput(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

var _DependencyRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dra.DependencyRPC",
	HandlerType: (*DependencyRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVMOffsets",
			Handler:    _DependencyRPC_GetVMOffsets_Handler,
		},
		{
			MethodName: "SendNumberBasicBlock",
			Handler:    _DependencyRPC_SendNumberBasicBlock_Handler,
		},
		{
			MethodName: "SendNumberBasicBlockCovered",
			Handler:    _DependencyRPC_SendNumberBasicBlockCovered_Handler,
		},
		{
			MethodName: "GetNewInput",
			Handler:    _DependencyRPC_GetNewInput_Handler,
		},
		{
			MethodName: "SendDependency",
			Handler:    _DependencyRPC_SendDependency_Handler,
		},
		{
			MethodName: "GetCondition",
			Handler:    _DependencyRPC_GetCondition_Handler,
		},
		{
			MethodName: "SendWriteAddress",
			Handler:    _DependencyRPC_SendWriteAddress_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _DependencyRPC_Connect_Handler,
		},
		{
			MethodName: "GetDataDependency",
			Handler:    _DependencyRPC_GetDataDependency_Handler,
		},
		{
			MethodName: "SendNewInput",
			Handler:    _DependencyRPC_SendNewInput_Handler,
		},
		{
			MethodName: "GetTasks",
			Handler:    _DependencyRPC_GetTasks_Handler,
		},
		{
			MethodName: "GetBootTasks",
			Handler:    _DependencyRPC_GetBootTasks_Handler,
		},
		{
			MethodName: "ReturnTasks",
			Handler:    _DependencyRPC_ReturnTasks_Handler,
		},
		{
			MethodName: "SendBootInput",
			Handler:    _DependencyRPC_SendBootInput_Handler,
		},
		{
			MethodName: "SendUnstableInput",
			Handler:    _DependencyRPC_SendUnstableInput_Handler,
		},
		{
			MethodName: "SendLog",
			Handler:    _DependencyRPC_SendLog_Handler,
		},
		{
			MethodName: "sendStat",
			Handler:    _DependencyRPC_SendStat_Handler,
		},
		{
			MethodName: "GetNeed",
			Handler:    _DependencyRPC_GetNeed_Handler,
		},
		{
			MethodName: "SendNeedInput",
			Handler:    _DependencyRPC_SendNeedInput_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "DependencyRPC.proto",
}
