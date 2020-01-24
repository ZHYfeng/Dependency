// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Input.proto

package dra

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type FuzzingStat int32

const (
	FuzzingStat_StatDefault        FuzzingStat = 0
	FuzzingStat_StatGenerate       FuzzingStat = 1
	FuzzingStat_StatFuzz           FuzzingStat = 2
	FuzzingStat_StatCandidate      FuzzingStat = 3
	FuzzingStat_StatTriage         FuzzingStat = 4
	FuzzingStat_StatMinimize       FuzzingStat = 5
	FuzzingStat_StatSmash          FuzzingStat = 6
	FuzzingStat_StatHint           FuzzingStat = 7
	FuzzingStat_StatSeed           FuzzingStat = 8
	FuzzingStat_StatDependency     FuzzingStat = 9
	FuzzingStat_StatDependencyBoot FuzzingStat = 10
)

var FuzzingStat_name = map[int32]string{
	0:  "StatDefault",
	1:  "StatGenerate",
	2:  "StatFuzz",
	3:  "StatCandidate",
	4:  "StatTriage",
	5:  "StatMinimize",
	6:  "StatSmash",
	7:  "StatHint",
	8:  "StatSeed",
	9:  "StatDependency",
	10: "StatDependencyBoot",
}

var FuzzingStat_value = map[string]int32{
	"StatDefault":        0,
	"StatGenerate":       1,
	"StatFuzz":           2,
	"StatCandidate":      3,
	"StatTriage":         4,
	"StatMinimize":       5,
	"StatSmash":          6,
	"StatHint":           7,
	"StatSeed":           8,
	"StatDependency":     9,
	"StatDependencyBoot": 10,
}

func (x FuzzingStat) String() string {
	return proto.EnumName(FuzzingStat_name, int32(x))
}

func (FuzzingStat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{0}
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
	return fileDescriptor_e4691306fcd7be97, []int{0}
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
	// for program
	Sig     string           `protobuf:"bytes,11,opt,name=sig,proto3" json:"sig,omitempty"`
	Program []byte           `protobuf:"bytes,12,opt,name=program,proto3" json:"program,omitempty"`
	Call    map[uint32]*Call `protobuf:"bytes,13,rep,name=call,proto3" json:"call,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Paths   []*Paths         `protobuf:"bytes,1,rep,name=paths,proto3" json:"paths,omitempty"`
	Stable  uint32           `protobuf:"varint,14,opt,name=stable,proto3" json:"stable,omitempty"`
	Total   uint32           `protobuf:"varint,15,opt,name=total,proto3" json:"total,omitempty"`
	// for dependency
	Stat FuzzingStat `protobuf:"varint,21,opt,name=stat,proto3,enum=dra.FuzzingStat" json:"stat,omitempty"`
	// uncovered address, index by bits
	UncoveredAddress map[uint32]uint32 `protobuf:"bytes,22,rep,name=uncovered_address,json=uncoveredAddress,proto3" json:"uncovered_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// write address, index by bits
	WriteAddress         map[uint32]uint32 `protobuf:"bytes,25,rep,name=write_address,json=writeAddress,proto3" json:"write_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	ProgramBeforeMini    []byte            `protobuf:"bytes,30,opt,name=program_before_mini,json=programBeforeMini,proto3" json:"program_before_mini,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Input) Reset()         { *m = Input{} }
func (m *Input) String() string { return proto.CompactTextString(m) }
func (*Input) ProtoMessage()    {}
func (*Input) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{1}
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

func (m *Input) GetProgram() []byte {
	if m != nil {
		return m.Program
	}
	return nil
}

func (m *Input) GetCall() map[uint32]*Call {
	if m != nil {
		return m.Call
	}
	return nil
}

func (m *Input) GetPaths() []*Paths {
	if m != nil {
		return m.Paths
	}
	return nil
}

func (m *Input) GetStable() uint32 {
	if m != nil {
		return m.Stable
	}
	return 0
}

func (m *Input) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *Input) GetStat() FuzzingStat {
	if m != nil {
		return m.Stat
	}
	return FuzzingStat_StatDefault
}

func (m *Input) GetUncoveredAddress() map[uint32]uint32 {
	if m != nil {
		return m.UncoveredAddress
	}
	return nil
}

func (m *Input) GetWriteAddress() map[uint32]uint32 {
	if m != nil {
		return m.WriteAddress
	}
	return nil
}

func (m *Input) GetProgramBeforeMini() []byte {
	if m != nil {
		return m.ProgramBeforeMini
	}
	return nil
}

type Inputs struct {
	// map<string, Input> input = 1;
	Input                []*Input `protobuf:"bytes,1,rep,name=input,proto3" json:"input,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Inputs) Reset()         { *m = Inputs{} }
func (m *Inputs) String() string { return proto.CompactTextString(m) }
func (*Inputs) ProtoMessage()    {}
func (*Inputs) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{2}
}

func (m *Inputs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Inputs.Unmarshal(m, b)
}
func (m *Inputs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Inputs.Marshal(b, m, deterministic)
}
func (m *Inputs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Inputs.Merge(m, src)
}
func (m *Inputs) XXX_Size() int {
	return xxx_messageInfo_Inputs.Size(m)
}
func (m *Inputs) XXX_DiscardUnknown() {
	xxx_messageInfo_Inputs.DiscardUnknown(m)
}

var xxx_messageInfo_Inputs proto.InternalMessageInfo

func (m *Inputs) GetInput() []*Input {
	if m != nil {
		return m.Input
	}
	return nil
}

type Path struct {
	Address              []uint32 `protobuf:"varint,1,rep,packed,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Path) Reset()         { *m = Path{} }
func (m *Path) String() string { return proto.CompactTextString(m) }
func (*Path) ProtoMessage()    {}
func (*Path) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{3}
}

func (m *Path) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Path.Unmarshal(m, b)
}
func (m *Path) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Path.Marshal(b, m, deterministic)
}
func (m *Path) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Path.Merge(m, src)
}
func (m *Path) XXX_Size() int {
	return xxx_messageInfo_Path.Size(m)
}
func (m *Path) XXX_DiscardUnknown() {
	xxx_messageInfo_Path.DiscardUnknown(m)
}

var xxx_messageInfo_Path proto.InternalMessageInfo

func (m *Path) GetAddress() []uint32 {
	if m != nil {
		return m.Address
	}
	return nil
}

type Paths struct {
	Path                 map[uint32]*Path `protobuf:"bytes,1,rep,name=path,proto3" json:"path,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Paths) Reset()         { *m = Paths{} }
func (m *Paths) String() string { return proto.CompactTextString(m) }
func (*Paths) ProtoMessage()    {}
func (*Paths) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{4}
}

func (m *Paths) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Paths.Unmarshal(m, b)
}
func (m *Paths) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Paths.Marshal(b, m, deterministic)
}
func (m *Paths) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Paths.Merge(m, src)
}
func (m *Paths) XXX_Size() int {
	return xxx_messageInfo_Paths.Size(m)
}
func (m *Paths) XXX_DiscardUnknown() {
	xxx_messageInfo_Paths.DiscardUnknown(m)
}

var xxx_messageInfo_Paths proto.InternalMessageInfo

func (m *Paths) GetPath() map[uint32]*Path {
	if m != nil {
		return m.Path
	}
	return nil
}

type UnstableInput struct {
	Sig          string   `protobuf:"bytes,1,opt,name=sig,proto3" json:"sig,omitempty"`
	Program      []byte   `protobuf:"bytes,2,opt,name=program,proto3" json:"program,omitempty"`
	UnstablePath []*Paths `protobuf:"bytes,12,rep,name=unstable_path,json=unstablePath,proto3" json:"unstable_path,omitempty"`
	// address, index by bits
	Address              map[uint32]uint32 `protobuf:"bytes,13,rep,name=address,proto3" json:"address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UnstableInput) Reset()         { *m = UnstableInput{} }
func (m *UnstableInput) String() string { return proto.CompactTextString(m) }
func (*UnstableInput) ProtoMessage()    {}
func (*UnstableInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{5}
}

func (m *UnstableInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnstableInput.Unmarshal(m, b)
}
func (m *UnstableInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnstableInput.Marshal(b, m, deterministic)
}
func (m *UnstableInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnstableInput.Merge(m, src)
}
func (m *UnstableInput) XXX_Size() int {
	return xxx_messageInfo_UnstableInput.Size(m)
}
func (m *UnstableInput) XXX_DiscardUnknown() {
	xxx_messageInfo_UnstableInput.DiscardUnknown(m)
}

var xxx_messageInfo_UnstableInput proto.InternalMessageInfo

func (m *UnstableInput) GetSig() string {
	if m != nil {
		return m.Sig
	}
	return ""
}

func (m *UnstableInput) GetProgram() []byte {
	if m != nil {
		return m.Program
	}
	return nil
}

func (m *UnstableInput) GetUnstablePath() []*Paths {
	if m != nil {
		return m.UnstablePath
	}
	return nil
}

func (m *UnstableInput) GetAddress() map[uint32]uint32 {
	if m != nil {
		return m.Address
	}
	return nil
}

type UnstableInputs struct {
	UnstableInput        map[string]*UnstableInput `protobuf:"bytes,1,rep,name=unstable_input,json=unstableInput,proto3" json:"unstable_input,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *UnstableInputs) Reset()         { *m = UnstableInputs{} }
func (m *UnstableInputs) String() string { return proto.CompactTextString(m) }
func (*UnstableInputs) ProtoMessage()    {}
func (*UnstableInputs) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4691306fcd7be97, []int{6}
}

func (m *UnstableInputs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnstableInputs.Unmarshal(m, b)
}
func (m *UnstableInputs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnstableInputs.Marshal(b, m, deterministic)
}
func (m *UnstableInputs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnstableInputs.Merge(m, src)
}
func (m *UnstableInputs) XXX_Size() int {
	return xxx_messageInfo_UnstableInputs.Size(m)
}
func (m *UnstableInputs) XXX_DiscardUnknown() {
	xxx_messageInfo_UnstableInputs.DiscardUnknown(m)
}

var xxx_messageInfo_UnstableInputs proto.InternalMessageInfo

func (m *UnstableInputs) GetUnstableInput() map[string]*UnstableInput {
	if m != nil {
		return m.UnstableInput
	}
	return nil
}

func init() {
	proto.RegisterEnum("dra.FuzzingStat", FuzzingStat_name, FuzzingStat_value)
	proto.RegisterType((*Call)(nil), "dra.Call")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.Call.AddressEntry")
	proto.RegisterType((*Input)(nil), "dra.Input")
	proto.RegisterMapType((map[uint32]*Call)(nil), "dra.Input.CallEntry")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.Input.UncoveredAddressEntry")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.Input.WriteAddressEntry")
	proto.RegisterType((*Inputs)(nil), "dra.Inputs")
	proto.RegisterType((*Path)(nil), "dra.Path")
	proto.RegisterType((*Paths)(nil), "dra.Paths")
	proto.RegisterMapType((map[uint32]*Path)(nil), "dra.Paths.PathEntry")
	proto.RegisterType((*UnstableInput)(nil), "dra.UnstableInput")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.UnstableInput.AddressEntry")
	proto.RegisterType((*UnstableInputs)(nil), "dra.UnstableInputs")
	proto.RegisterMapType((map[string]*UnstableInput)(nil), "dra.UnstableInputs.UnstableInputEntry")
}

func init() { proto.RegisterFile("Input.proto", fileDescriptor_e4691306fcd7be97) }

var fileDescriptor_e4691306fcd7be97 = []byte{
	// 674 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xcf, 0x4e, 0xdb, 0x4e,
	0x10, 0xfe, 0x6d, 0xfe, 0x41, 0x26, 0x76, 0xd8, 0xcc, 0x0f, 0x22, 0x37, 0xaa, 0x8a, 0x15, 0x55,
	0x55, 0xc4, 0x21, 0xad, 0xe8, 0xa5, 0xe5, 0x52, 0x01, 0xfd, 0x7b, 0x40, 0xaa, 0x0c, 0xa8, 0xc7,
	0x68, 0xc1, 0x4b, 0x58, 0xe1, 0xd8, 0x91, 0xbd, 0x86, 0x86, 0x07, 0xe8, 0x1b, 0xf5, 0x25, 0xfa,
	0x12, 0x3d, 0xf4, 0x45, 0xaa, 0x59, 0xc7, 0xc6, 0x6e, 0x38, 0x94, 0x43, 0x2f, 0xe0, 0xf9, 0xe6,
	0xdb, 0xf1, 0x37, 0xb3, 0xdf, 0x38, 0xd0, 0xf9, 0x14, 0xce, 0x53, 0x3d, 0x9e, 0xc7, 0x91, 0x8e,
	0xb0, 0xee, 0xc7, 0x62, 0xf8, 0x8d, 0x41, 0xe3, 0x50, 0x04, 0x01, 0x72, 0xa8, 0x2b, 0xff, 0xab,
	0xc3, 0x5c, 0x36, 0xb2, 0x3d, 0x7a, 0xc4, 0x17, 0xb0, 0x26, 0x7c, 0x3f, 0x96, 0x49, 0xe2, 0xd4,
	0xdc, 0xfa, 0xa8, 0xb3, 0xdb, 0x1f, 0xfb, 0xb1, 0x18, 0x13, 0x7b, 0xbc, 0x9f, 0x25, 0xde, 0x85,
	0x3a, 0x5e, 0x78, 0x39, 0x6d, 0xb0, 0x07, 0x56, 0x39, 0x41, 0x35, 0xaf, 0xe4, 0x22, 0xaf, 0x79,
	0x25, 0x17, 0xb8, 0x09, 0xcd, 0x6b, 0x11, 0xa4, 0xd2, 0xa9, 0x19, 0x2c, 0x0b, 0xf6, 0x6a, 0xaf,
	0xd8, 0xf0, 0x67, 0x03, 0x9a, 0x46, 0x1d, 0x9d, 0x4a, 0xd4, 0xd4, 0xe9, 0xb8, 0x6c, 0xd4, 0xf6,
	0xe8, 0x11, 0x1d, 0x58, 0x9b, 0xc7, 0xd1, 0x34, 0x16, 0x33, 0xc7, 0x72, 0xd9, 0xc8, 0xf2, 0xf2,
	0x10, 0x47, 0xd0, 0x38, 0x17, 0x41, 0xe0, 0xd8, 0x46, 0xe0, 0xa6, 0x11, 0x98, 0xf5, 0x48, 0x32,
	0x33, 0x79, 0x86, 0x81, 0x2e, 0x34, 0xe7, 0x42, 0x5f, 0x26, 0x0e, 0x33, 0x54, 0x30, 0xd4, 0xcf,
	0x84, 0x78, 0x59, 0x02, 0xfb, 0xd0, 0x4a, 0xb4, 0x38, 0x0b, 0xa4, 0xd3, 0x35, 0xe2, 0x96, 0x11,
	0x69, 0xd6, 0x91, 0x16, 0x81, 0xb3, 0x91, 0x69, 0x36, 0x01, 0x3e, 0x85, 0x46, 0xa2, 0x85, 0x76,
	0xb6, 0x5c, 0x36, 0xea, 0xee, 0x72, 0x53, 0xee, 0x7d, 0x7a, 0x7b, 0xab, 0xc2, 0xe9, 0xb1, 0x16,
	0xda, 0x33, 0x59, 0x3c, 0x82, 0x5e, 0x1a, 0x9e, 0x47, 0xd7, 0x32, 0x96, 0xfe, 0x24, 0x9f, 0x66,
	0xdf, 0x28, 0x70, 0x4b, 0x62, 0x4f, 0x73, 0x4e, 0x65, 0xae, 0x3c, 0xfd, 0x03, 0xc6, 0x7d, 0xb0,
	0x6f, 0x62, 0xa5, 0x65, 0x51, 0xea, 0x91, 0x29, 0xf5, 0xb8, 0x54, 0xea, 0x0b, 0xe5, 0x2b, 0x65,
	0xac, 0x9b, 0x12, 0x84, 0x63, 0xf8, 0x7f, 0x39, 0xbc, 0xc9, 0x99, 0xbc, 0x88, 0x62, 0x39, 0x99,
	0xa9, 0x50, 0x39, 0x4f, 0xcc, 0x5c, 0x7b, 0xcb, 0xd4, 0x81, 0xc9, 0x1c, 0xa9, 0x50, 0x0d, 0x0e,
	0xa0, 0x5d, 0x8c, 0xf2, 0x9e, 0x0b, 0xdd, 0x2e, 0x5f, 0x68, 0x67, 0xb7, 0x5d, 0x58, 0xa4, 0x74,
	0xb7, 0x83, 0x43, 0xd8, 0xba, 0xb7, 0xc3, 0x87, 0x18, 0x64, 0xf0, 0x06, 0x7a, 0x2b, 0xbd, 0x3d,
	0xc8, 0x61, 0x3b, 0xd0, 0x32, 0x23, 0x4a, 0xc8, 0x0b, 0x8a, 0x9e, 0x2a, 0x5e, 0x30, 0x39, 0x2f,
	0x4b, 0x0c, 0x5d, 0x68, 0x90, 0x37, 0xc8, 0x79, 0xf9, 0xa8, 0x89, 0x6b, 0x17, 0x5e, 0x1f, 0xa6,
	0xd0, 0x34, 0xee, 0x21, 0x0b, 0x92, 0x7f, 0x96, 0xb5, 0x36, 0xef, 0x7c, 0x65, 0xfe, 0x2e, 0x2d,
	0x48, 0x0c, 0x1a, 0x65, 0x01, 0xfd, 0xed, 0x28, 0xe9, 0x40, 0xb9, 0x89, 0x5f, 0x0c, 0xec, 0xd3,
	0x30, 0x73, 0x66, 0x65, 0x5d, 0xd8, 0xbd, 0xeb, 0x52, 0xab, 0xae, 0xcb, 0x73, 0xb0, 0xd3, 0xe5,
	0xe1, 0x89, 0x11, 0x6d, 0xad, 0x2c, 0x83, 0x95, 0x13, 0x4c, 0xff, 0xaf, 0xef, 0xfa, 0xcf, 0x56,
	0x6c, 0xdb, 0x50, 0x2b, 0x0a, 0xfe, 0xc1, 0xc7, 0xe0, 0x3b, 0x83, 0x6e, 0xe5, 0x1d, 0x09, 0x1e,
	0x41, 0xb7, 0x90, 0x5e, 0xbe, 0xbc, 0x67, 0xab, 0x82, 0x92, 0x6a, 0x98, 0xe9, 0x2a, 0x1a, 0x37,
	0xd8, 0xe0, 0x04, 0x70, 0x95, 0x54, 0xd6, 0xd8, 0xce, 0x34, 0x8e, 0xaa, 0x97, 0x82, 0xab, 0x6f,
	0x2b, 0xe9, 0xde, 0xf9, 0xc1, 0xa0, 0x53, 0xfa, 0x08, 0xe0, 0x06, 0x74, 0xe8, 0xff, 0x5b, 0x79,
	0x21, 0xd2, 0x40, 0xf3, 0xff, 0x90, 0x83, 0x45, 0xc0, 0x07, 0x19, 0xca, 0x58, 0x68, 0xc9, 0x19,
	0x5a, 0xb0, 0x4e, 0x08, 0x9d, 0xe2, 0x35, 0xec, 0x81, 0x4d, 0xd1, 0xa1, 0x08, 0x7d, 0xe5, 0x13,
	0xa1, 0x8e, 0x5d, 0x00, 0x82, 0x4e, 0x62, 0x25, 0xa6, 0x92, 0x37, 0xf2, 0x12, 0xb4, 0x9c, 0x33,
	0x75, 0x2b, 0x79, 0x13, 0x6d, 0x68, 0x13, 0x72, 0x3c, 0x13, 0xc9, 0x25, 0x6f, 0xe5, 0x15, 0x3f,
	0xaa, 0x50, 0xf3, 0xb5, 0x3c, 0x3a, 0x96, 0xd2, 0xe7, 0xeb, 0x88, 0xd0, 0xcd, 0x04, 0xcd, 0x65,
	0xe8, 0xcb, 0xf0, 0x7c, 0xc1, 0xdb, 0xd8, 0x07, 0xac, 0x62, 0x07, 0x51, 0xa4, 0x39, 0x9c, 0xb5,
	0xcc, 0xcf, 0xc4, 0xcb, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x93, 0x08, 0xa7, 0x35, 0x06,
	0x00, 0x00,
}