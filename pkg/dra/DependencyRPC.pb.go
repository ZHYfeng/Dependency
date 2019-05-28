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

type Empty struct {
	Address              uint32   `protobuf:"varint,1,opt,name=address,proto3" json:"address,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func (m *Empty) GetAddress() uint32 {
	if m != nil {
		return m.Address
	}
	return 0
}

func (m *Empty) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RelatedInput struct {
	Sig                  string   `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	Prog                 []byte   `protobuf:"bytes,3,opt,name=prog,proto3" json:"prog,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelatedInput) Reset()         { *m = RelatedInput{} }
func (m *RelatedInput) String() string { return proto.CompactTextString(m) }
func (*RelatedInput) ProtoMessage()    {}
func (*RelatedInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{1}
}

func (m *RelatedInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelatedInput.Unmarshal(m, b)
}
func (m *RelatedInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelatedInput.Marshal(b, m, deterministic)
}
func (m *RelatedInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelatedInput.Merge(m, src)
}
func (m *RelatedInput) XXX_Size() int {
	return xxx_messageInfo_RelatedInput.Size(m)
}
func (m *RelatedInput) XXX_DiscardUnknown() {
	xxx_messageInfo_RelatedInput.DiscardUnknown(m)
}

var xxx_messageInfo_RelatedInput proto.InternalMessageInfo

func (m *RelatedInput) GetSig() string {
	if m != nil {
		return m.Sig
	}
	return ""
}

func (m *RelatedInput) GetProg() []byte {
	if m != nil {
		return m.Prog
	}
	return nil
}

type RelatedSyscall struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Number               uint64   `protobuf:"varint,3,opt,name=number,proto3" json:"number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelatedSyscall) Reset()         { *m = RelatedSyscall{} }
func (m *RelatedSyscall) String() string { return proto.CompactTextString(m) }
func (*RelatedSyscall) ProtoMessage()    {}
func (*RelatedSyscall) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{2}
}

func (m *RelatedSyscall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelatedSyscall.Unmarshal(m, b)
}
func (m *RelatedSyscall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelatedSyscall.Marshal(b, m, deterministic)
}
func (m *RelatedSyscall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelatedSyscall.Merge(m, src)
}
func (m *RelatedSyscall) XXX_Size() int {
	return xxx_messageInfo_RelatedSyscall.Size(m)
}
func (m *RelatedSyscall) XXX_DiscardUnknown() {
	xxx_messageInfo_RelatedSyscall.DiscardUnknown(m)
}

var xxx_messageInfo_RelatedSyscall proto.InternalMessageInfo

func (m *RelatedSyscall) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RelatedSyscall) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

type RelatedAddress struct {
	Address              uint32            `protobuf:"varint,2,opt,name=address,proto3" json:"address,omitempty"`
	RelatedInput         []*RelatedInput   `protobuf:"bytes,4,rep,name=related_input,json=relatedInput,proto3" json:"related_input,omitempty"`
	RelatedSyscall       []*RelatedSyscall `protobuf:"bytes,5,rep,name=related_syscall,json=relatedSyscall,proto3" json:"related_syscall,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RelatedAddress) Reset()         { *m = RelatedAddress{} }
func (m *RelatedAddress) String() string { return proto.CompactTextString(m) }
func (*RelatedAddress) ProtoMessage()    {}
func (*RelatedAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{3}
}

func (m *RelatedAddress) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelatedAddress.Unmarshal(m, b)
}
func (m *RelatedAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelatedAddress.Marshal(b, m, deterministic)
}
func (m *RelatedAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelatedAddress.Merge(m, src)
}
func (m *RelatedAddress) XXX_Size() int {
	return xxx_messageInfo_RelatedAddress.Size(m)
}
func (m *RelatedAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_RelatedAddress.DiscardUnknown(m)
}

var xxx_messageInfo_RelatedAddress proto.InternalMessageInfo

func (m *RelatedAddress) GetAddress() uint32 {
	if m != nil {
		return m.Address
	}
	return 0
}

func (m *RelatedAddress) GetRelatedInput() []*RelatedInput {
	if m != nil {
		return m.RelatedInput
	}
	return nil
}

func (m *RelatedAddress) GetRelatedSyscall() []*RelatedSyscall {
	if m != nil {
		return m.RelatedSyscall
	}
	return nil
}

type UncoveredAddress struct {
	Address              uint32            `protobuf:"varint,1,opt,name=address,proto3" json:"address,omitempty"`
	Idx                  uint32            `protobuf:"varint,2,opt,name=idx,proto3" json:"idx,omitempty"`
	ConditionAddress     uint32            `protobuf:"varint,3,opt,name=condition_address,json=conditionAddress,proto3" json:"condition_address,omitempty"`
	RelatedAddress       []*RelatedAddress `protobuf:"bytes,4,rep,name=related_address,json=relatedAddress,proto3" json:"related_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UncoveredAddress) Reset()         { *m = UncoveredAddress{} }
func (m *UncoveredAddress) String() string { return proto.CompactTextString(m) }
func (*UncoveredAddress) ProtoMessage()    {}
func (*UncoveredAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{4}
}

func (m *UncoveredAddress) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UncoveredAddress.Unmarshal(m, b)
}
func (m *UncoveredAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UncoveredAddress.Marshal(b, m, deterministic)
}
func (m *UncoveredAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UncoveredAddress.Merge(m, src)
}
func (m *UncoveredAddress) XXX_Size() int {
	return xxx_messageInfo_UncoveredAddress.Size(m)
}
func (m *UncoveredAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_UncoveredAddress.DiscardUnknown(m)
}

var xxx_messageInfo_UncoveredAddress proto.InternalMessageInfo

func (m *UncoveredAddress) GetAddress() uint32 {
	if m != nil {
		return m.Address
	}
	return 0
}

func (m *UncoveredAddress) GetIdx() uint32 {
	if m != nil {
		return m.Idx
	}
	return 0
}

func (m *UncoveredAddress) GetConditionAddress() uint32 {
	if m != nil {
		return m.ConditionAddress
	}
	return 0
}

func (m *UncoveredAddress) GetRelatedAddress() []*RelatedAddress {
	if m != nil {
		return m.RelatedAddress
	}
	return nil
}

type DependencyInput struct {
	Sig                  string              `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	UncoveredAddress     []*UncoveredAddress `protobuf:"bytes,2,rep,name=uncovered_address,json=uncoveredAddress,proto3" json:"uncovered_address,omitempty"`
	Prog                 []byte              `protobuf:"bytes,3,opt,name=prog,proto3" json:"prog,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *DependencyInput) Reset()         { *m = DependencyInput{} }
func (m *DependencyInput) String() string { return proto.CompactTextString(m) }
func (*DependencyInput) ProtoMessage()    {}
func (*DependencyInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{5}
}

func (m *DependencyInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DependencyInput.Unmarshal(m, b)
}
func (m *DependencyInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DependencyInput.Marshal(b, m, deterministic)
}
func (m *DependencyInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DependencyInput.Merge(m, src)
}
func (m *DependencyInput) XXX_Size() int {
	return xxx_messageInfo_DependencyInput.Size(m)
}
func (m *DependencyInput) XXX_DiscardUnknown() {
	xxx_messageInfo_DependencyInput.DiscardUnknown(m)
}

var xxx_messageInfo_DependencyInput proto.InternalMessageInfo

func (m *DependencyInput) GetSig() string {
	if m != nil {
		return m.Sig
	}
	return ""
}

func (m *DependencyInput) GetUncoveredAddress() []*UncoveredAddress {
	if m != nil {
		return m.UncoveredAddress
	}
	return nil
}

func (m *DependencyInput) GetProg() []byte {
	if m != nil {
		return m.Prog
	}
	return nil
}

type NewDependencyInput struct {
	DependencyInput      []*DependencyInput `protobuf:"bytes,1,rep,name=dependencyInput,proto3" json:"dependencyInput,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NewDependencyInput) Reset()         { *m = NewDependencyInput{} }
func (m *NewDependencyInput) String() string { return proto.CompactTextString(m) }
func (*NewDependencyInput) ProtoMessage()    {}
func (*NewDependencyInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{6}
}

func (m *NewDependencyInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewDependencyInput.Unmarshal(m, b)
}
func (m *NewDependencyInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewDependencyInput.Marshal(b, m, deterministic)
}
func (m *NewDependencyInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewDependencyInput.Merge(m, src)
}
func (m *NewDependencyInput) XXX_Size() int {
	return xxx_messageInfo_NewDependencyInput.Size(m)
}
func (m *NewDependencyInput) XXX_DiscardUnknown() {
	xxx_messageInfo_NewDependencyInput.DiscardUnknown(m)
}

var xxx_messageInfo_NewDependencyInput proto.InternalMessageInfo

func (m *NewDependencyInput) GetDependencyInput() []*DependencyInput {
	if m != nil {
		return m.DependencyInput
	}
	return nil
}

type Call struct {
	Idx                  uint32            `protobuf:"varint,1,opt,name=idx,proto3" json:"idx,omitempty"`
	Address              map[uint32]uint32 `protobuf:"bytes,2,rep,name=address,proto3" json:"address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Call) Reset()         { *m = Call{} }
func (m *Call) String() string { return proto.CompactTextString(m) }
func (*Call) ProtoMessage()    {}
func (*Call) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{7}
}

func (m *Call) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Call.Unmarshal(m, b)
}
func (m *Call) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Call.Marshal(b, m, deterministic)
}
func (m *Call) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Call.Merge(m, src)
}
func (m *Call) XXX_Size() int {
	return xxx_messageInfo_Call.Size(m)
}
func (m *Call) XXX_DiscardUnknown() {
	xxx_messageInfo_Call.DiscardUnknown(m)
}

var xxx_messageInfo_Call proto.InternalMessageInfo

func (m *Call) GetIdx() uint32 {
	if m != nil {
		return m.Idx
	}
	return 0
}

func (m *Call) GetAddress() map[uint32]uint32 {
	if m != nil {
		return m.Address
	}
	return nil
}

type Input struct {
	Sig                  string           `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	Call                 map[uint32]*Call `protobuf:"bytes,2,rep,name=call,proto3" json:"call,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Prog                 []byte           `protobuf:"bytes,3,opt,name=prog,proto3" json:"prog,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Input) Reset()         { *m = Input{} }
func (m *Input) String() string { return proto.CompactTextString(m) }
func (*Input) ProtoMessage()    {}
func (*Input) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{8}
}

func (m *Input) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Input.Unmarshal(m, b)
}
func (m *Input) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Input.Marshal(b, m, deterministic)
}
func (m *Input) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Input.Merge(m, src)
}
func (m *Input) XXX_Size() int {
	return xxx_messageInfo_Input.Size(m)
}
func (m *Input) XXX_DiscardUnknown() {
	xxx_messageInfo_Input.DiscardUnknown(m)
}

var xxx_messageInfo_Input proto.InternalMessageInfo

func (m *Input) GetSig() string {
	if m != nil {
		return m.Sig
	}
	return ""
}

func (m *Input) GetCall() map[uint32]*Call {
	if m != nil {
		return m.Call
	}
	return nil
}

func (m *Input) GetProg() []byte {
	if m != nil {
		return m.Prog
	}
	return nil
}

type NewInput struct {
	Input                []*Input `protobuf:"bytes,1,rep,name=input,proto3" json:"input,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewInput) Reset()         { *m = NewInput{} }
func (m *NewInput) String() string { return proto.CompactTextString(m) }
func (*NewInput) ProtoMessage()    {}
func (*NewInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_db4d5fd3d0a7c985, []int{9}
}

func (m *NewInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewInput.Unmarshal(m, b)
}
func (m *NewInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewInput.Marshal(b, m, deterministic)
}
func (m *NewInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewInput.Merge(m, src)
}
func (m *NewInput) XXX_Size() int {
	return xxx_messageInfo_NewInput.Size(m)
}
func (m *NewInput) XXX_DiscardUnknown() {
	xxx_messageInfo_NewInput.DiscardUnknown(m)
}

var xxx_messageInfo_NewInput proto.InternalMessageInfo

func (m *NewInput) GetInput() []*Input {
	if m != nil {
		return m.Input
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "dra.Empty")
	proto.RegisterType((*RelatedInput)(nil), "dra.RelatedInput")
	proto.RegisterType((*RelatedSyscall)(nil), "dra.RelatedSyscall")
	proto.RegisterType((*RelatedAddress)(nil), "dra.RelatedAddress")
	proto.RegisterType((*UncoveredAddress)(nil), "dra.UncoveredAddress")
	proto.RegisterType((*DependencyInput)(nil), "dra.DependencyInput")
	proto.RegisterType((*NewDependencyInput)(nil), "dra.NewDependencyInput")
	proto.RegisterType((*Call)(nil), "dra.Call")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.Call.AddressEntry")
	proto.RegisterType((*Input)(nil), "dra.Input")
	proto.RegisterMapType((map[uint32]*Call)(nil), "dra.Input.CallEntry")
	proto.RegisterType((*NewInput)(nil), "dra.NewInput")
}

func init() { proto.RegisterFile("DependencyRPC.proto", fileDescriptor_db4d5fd3d0a7c985) }

var fileDescriptor_db4d5fd3d0a7c985 = []byte{
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0xeb, 0x26, 0xdd, 0xe8, 0x6b, 0xbb, 0xb6, 0xee, 0x18, 0x51, 0x2f, 0x54, 0x46, 0x48,
	0x11, 0xa0, 0x0a, 0x95, 0x1f, 0x82, 0x31, 0x21, 0xb1, 0x32, 0x55, 0x5c, 0x06, 0xf2, 0x80, 0x6b,
	0x95, 0x35, 0xde, 0x14, 0xad, 0x75, 0x2a, 0xc7, 0xdd, 0x88, 0xb8, 0xf3, 0x4f, 0x70, 0xe1, 0xba,
	0xff, 0x12, 0xc5, 0x89, 0x53, 0xd7, 0x74, 0xdc, 0x9e, 0xeb, 0xef, 0xf7, 0xbd, 0xcf, 0x7b, 0x7e,
	0x0d, 0xf4, 0x3e, 0xb2, 0x25, 0xe3, 0x21, 0xe3, 0xb3, 0x94, 0x7e, 0x19, 0x0f, 0x97, 0x22, 0x96,
	0x31, 0x76, 0x42, 0x11, 0x90, 0x57, 0x50, 0x3b, 0x59, 0x2c, 0x65, 0x8a, 0x3d, 0xd8, 0x0d, 0xc2,
	0x50, 0xb0, 0x24, 0xf1, 0xd0, 0x00, 0xf9, 0x2d, 0xaa, 0x8f, 0x18, 0x83, 0xcb, 0x83, 0x05, 0xf3,
	0xaa, 0x03, 0xe4, 0xd7, 0xa9, 0x8a, 0xc9, 0x4b, 0x68, 0x52, 0x36, 0x0f, 0x24, 0x0b, 0x3f, 0xf1,
	0xe5, 0x4a, 0xe2, 0x0e, 0x38, 0x49, 0x74, 0xa9, 0x9c, 0x75, 0x9a, 0x85, 0x99, 0x6b, 0x29, 0xe2,
	0x4b, 0xcf, 0x19, 0x20, 0xbf, 0x49, 0x55, 0x4c, 0x8e, 0x60, 0xaf, 0x70, 0x9d, 0xa5, 0xc9, 0x2c,
	0x98, 0xcf, 0xcb, 0xdc, 0x68, 0x9d, 0x1b, 0x1f, 0xc0, 0x0e, 0x5f, 0x2d, 0xce, 0x99, 0x50, 0x5e,
	0x97, 0x16, 0x27, 0xf2, 0x07, 0x95, 0xf6, 0x0f, 0x05, 0x9a, 0x01, 0x5d, 0xdd, 0x84, 0x7e, 0x0d,
	0x2d, 0x91, 0x6b, 0xa7, 0x51, 0x46, 0xe8, 0xb9, 0x03, 0xc7, 0x6f, 0x8c, 0xba, 0xc3, 0x50, 0x04,
	0x43, 0x13, 0x9d, 0x36, 0x85, 0xd9, 0xc8, 0x11, 0xb4, 0xb5, 0x2f, 0xc9, 0x19, 0xbd, 0x9a, 0x72,
	0xf6, 0x4c, 0x67, 0x81, 0x4f, 0xf7, 0xc4, 0xc6, 0x99, 0xdc, 0x22, 0xe8, 0x7c, 0xe3, 0xb3, 0xf8,
	0x9a, 0x89, 0xad, 0x90, 0xd6, 0x64, 0x3b, 0xe0, 0x44, 0xe1, 0x8f, 0x02, 0x3d, 0x0b, 0xf1, 0x53,
	0xe8, 0xce, 0x62, 0x1e, 0x46, 0x32, 0x8a, 0xf9, 0x54, 0xbb, 0x1c, 0x75, 0xdf, 0x29, 0x2f, 0x74,
	0x62, 0x83, 0x55, 0x4b, 0xdd, 0x7f, 0x59, 0x0b, 0x75, 0xc9, 0x5a, 0x9c, 0xc9, 0x4f, 0x68, 0xaf,
	0xb7, 0xe2, 0xae, 0x57, 0x3c, 0x86, 0xee, 0x4a, 0xf7, 0x33, 0x5d, 0x8f, 0x3a, 0x2b, 0x72, 0x5f,
	0x15, 0xb1, 0xbb, 0xa5, 0x9d, 0x95, 0xdd, 0xff, 0xb6, 0x4d, 0xf8, 0x0a, 0xf8, 0x94, 0xdd, 0xd8,
	0xf5, 0xdf, 0x43, 0x3b, 0xdc, 0xfc, 0xc9, 0x43, 0xaa, 0xd6, 0xbe, 0xaa, 0x65, 0xc9, 0xa9, 0x2d,
	0x26, 0xbf, 0x10, 0xb8, 0xe3, 0x6c, 0xad, 0x8a, 0xc1, 0xa2, 0xf5, 0x60, 0x9f, 0x9b, 0x9b, 0x92,
	0xa5, 0x3c, 0x50, 0x29, 0x33, 0xf5, 0xb0, 0x00, 0x3d, 0xe1, 0x52, 0xa4, 0xe5, 0xe3, 0xf4, 0x0f,
	0xa1, 0x69, 0x5e, 0x64, 0x39, 0xaf, 0x58, 0xaa, 0x73, 0x5e, 0xb1, 0x14, 0xef, 0x43, 0xed, 0x3a,
	0x98, 0xaf, 0x58, 0xf1, 0x80, 0xf9, 0xe1, 0xb0, 0xfa, 0x06, 0x91, 0xdf, 0x08, 0x6a, 0x77, 0x8d,
	0xd4, 0x07, 0x57, 0xad, 0x55, 0xd5, 0xe8, 0x4c, 0x69, 0x15, 0x4c, 0x0e, 0xe1, 0xea, 0x3f, 0x87,
	0x3d, 0xb8, 0xfe, 0x31, 0xd4, 0x4b, 0xd9, 0x16, 0xa4, 0x87, 0x26, 0x52, 0x63, 0x54, 0x2f, 0x9b,
	0x34, 0xe9, 0x9e, 0xc1, 0xbd, 0x53, 0x76, 0x93, 0xf3, 0x0d, 0xa0, 0x16, 0x19, 0x83, 0x86, 0x35,
	0x0e, 0xcd, 0x2f, 0x46, 0xb7, 0x55, 0x68, 0x6d, 0x7c, 0x3e, 0xb0, 0x0f, 0xcd, 0x09, 0x93, 0xdf,
	0x17, 0x9f, 0x2f, 0x2e, 0x12, 0x26, 0x13, 0x9c, 0x9b, 0xd4, 0x67, 0xa4, 0x6f, 0xc4, 0xa4, 0x82,
	0x9f, 0x40, 0x63, 0xc2, 0x64, 0x59, 0xcc, 0x14, 0xb6, 0x54, 0xac, 0xaf, 0x48, 0x05, 0xbf, 0x85,
	0xde, 0x19, 0xe3, 0xa1, 0xbd, 0x13, 0x5b, 0x9f, 0xde, 0x2a, 0xf3, 0x08, 0x76, 0xc7, 0x31, 0xe7,
	0x6c, 0x26, 0xff, 0xc3, 0xf2, 0x0e, 0xf0, 0x84, 0x49, 0x3b, 0xbd, 0xa9, 0x7f, 0xa0, 0x91, 0x2c,
	0x11, 0xa9, 0xe0, 0xc7, 0x50, 0xcf, 0xe0, 0x4c, 0xcf, 0x16, 0x90, 0xf3, 0x1d, 0xf5, 0x65, 0x7d,
	0xf1, 0x37, 0x00, 0x00, 0xff, 0xff, 0x79, 0xd3, 0x5d, 0x60, 0x70, 0x05, 0x00, 0x00,
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
	// C++ and syz-manager
	GetVmOffsets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	GetNewInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NewInput, error)
	SendDependencyInput(ctx context.Context, in *DependencyInput, opts ...grpc.CallOption) (*Empty, error)
	//syz-fuzzer and syz-manager
	Connect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	GetDependencyInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NewDependencyInput, error)
	SendInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error)
}

type dependencyRPCClient struct {
	cc *grpc.ClientConn
}

func NewDependencyRPCClient(cc *grpc.ClientConn) DependencyRPCClient {
	return &dependencyRPCClient{cc}
}

func (c *dependencyRPCClient) GetVmOffsets(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetVmOffsets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) GetNewInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NewInput, error) {
	out := new(NewInput)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetNewInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendDependencyInput(ctx context.Context, in *DependencyInput, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendDependencyInput", in, out, opts...)
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

func (c *dependencyRPCClient) GetDependencyInput(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NewDependencyInput, error) {
	out := new(NewDependencyInput)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/GetDependencyInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyRPCClient) SendInput(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dra.DependencyRPC/SendInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DependencyRPCServer is the server API for DependencyRPC service.
type DependencyRPCServer interface {
	// C++ and syz-manager
	GetVmOffsets(context.Context, *Empty) (*Empty, error)
	GetNewInput(context.Context, *Empty) (*NewInput, error)
	SendDependencyInput(context.Context, *DependencyInput) (*Empty, error)
	//syz-fuzzer and syz-manager
	Connect(context.Context, *Empty) (*Empty, error)
	GetDependencyInput(context.Context, *Empty) (*NewDependencyInput, error)
	SendInput(context.Context, *Input) (*Empty, error)
}

// UnimplementedDependencyRPCServer can be embedded to have forward compatible implementations.
type UnimplementedDependencyRPCServer struct {
}

func (*UnimplementedDependencyRPCServer) GetVmOffsets(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVmOffsets not implemented")
}
func (*UnimplementedDependencyRPCServer) GetNewInput(ctx context.Context, req *Empty) (*NewInput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewInput not implemented")
}
func (*UnimplementedDependencyRPCServer) SendDependencyInput(ctx context.Context, req *DependencyInput) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDependencyInput not implemented")
}
func (*UnimplementedDependencyRPCServer) Connect(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (*UnimplementedDependencyRPCServer) GetDependencyInput(ctx context.Context, req *Empty) (*NewDependencyInput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDependencyInput not implemented")
}
func (*UnimplementedDependencyRPCServer) SendInput(ctx context.Context, req *Input) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendInput not implemented")
}

func RegisterDependencyRPCServer(s *grpc.Server, srv DependencyRPCServer) {
	s.RegisterService(&_DependencyRPC_serviceDesc, srv)
}

func _DependencyRPC_GetVmOffsets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetVmOffsets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetVmOffsets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetVmOffsets(ctx, req.(*Empty))
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

func _DependencyRPC_SendDependencyInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DependencyInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendDependencyInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendDependencyInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendDependencyInput(ctx, req.(*DependencyInput))
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

func _DependencyRPC_GetDependencyInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).GetDependencyInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/GetDependencyInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).GetDependencyInput(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyRPC_SendInput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyRPCServer).SendInput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dra.DependencyRPC/SendInput",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyRPCServer).SendInput(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

var _DependencyRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dra.DependencyRPC",
	HandlerType: (*DependencyRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVmOffsets",
			Handler:    _DependencyRPC_GetVmOffsets_Handler,
		},
		{
			MethodName: "GetNewInput",
			Handler:    _DependencyRPC_GetNewInput_Handler,
		},
		{
			MethodName: "SendDependencyInput",
			Handler:    _DependencyRPC_SendDependencyInput_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _DependencyRPC_Connect_Handler,
		},
		{
			MethodName: "GetDependencyInput",
			Handler:    _DependencyRPC_GetDependencyInput_Handler,
		},
		{
			MethodName: "SendInput",
			Handler:    _DependencyRPC_SendInput_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "DependencyRPC.proto",
}
