// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Statistic.proto

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

type Statistic struct {
	Name                 FuzzingStat `protobuf:"varint,1,opt,name=name,proto3,enum=dra.FuzzingStat" json:"name,omitempty"`
	ExecuteNum           uint64      `protobuf:"varint,11,opt,name=executeNum,proto3" json:"executeNum,omitempty"`
	Time                 float64     `protobuf:"fixed64,12,opt,name=time,proto3" json:"time,omitempty"`
	NewTestCaseNum       uint64      `protobuf:"varint,13,opt,name=newTestCaseNum,proto3" json:"newTestCaseNum,omitempty"`
	NewAddressNum        uint64      `protobuf:"varint,14,opt,name=newAddressNum,proto3" json:"newAddressNum,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Statistic) Reset()         { *m = Statistic{} }
func (m *Statistic) String() string { return proto.CompactTextString(m) }
func (*Statistic) ProtoMessage()    {}
func (*Statistic) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8fb87f92dbb7b88, []int{0}
}

func (m *Statistic) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Statistic.Unmarshal(m, b)
}
func (m *Statistic) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Statistic.Marshal(b, m, deterministic)
}
func (m *Statistic) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Statistic.Merge(m, src)
}
func (m *Statistic) XXX_Size() int {
	return xxx_messageInfo_Statistic.Size(m)
}
func (m *Statistic) XXX_DiscardUnknown() {
	xxx_messageInfo_Statistic.DiscardUnknown(m)
}

var xxx_messageInfo_Statistic proto.InternalMessageInfo

func (m *Statistic) GetName() FuzzingStat {
	if m != nil {
		return m.Name
	}
	return FuzzingStat_StatGenerate
}

func (m *Statistic) GetExecuteNum() uint64 {
	if m != nil {
		return m.ExecuteNum
	}
	return 0
}

func (m *Statistic) GetTime() float64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Statistic) GetNewTestCaseNum() uint64 {
	if m != nil {
		return m.NewTestCaseNum
	}
	return 0
}

func (m *Statistic) GetNewAddressNum() uint64 {
	if m != nil {
		return m.NewAddressNum
	}
	return 0
}

type Time struct {
	Time                 float64  `protobuf:"fixed64,1,opt,name=time,proto3" json:"time,omitempty"`
	Num                  int64    `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`
	ExecuteNum           int64    `protobuf:"varint,3,opt,name=executeNum,proto3" json:"executeNum,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Time) Reset()         { *m = Time{} }
func (m *Time) String() string { return proto.CompactTextString(m) }
func (*Time) ProtoMessage()    {}
func (*Time) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8fb87f92dbb7b88, []int{1}
}

func (m *Time) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Time.Unmarshal(m, b)
}
func (m *Time) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Time.Marshal(b, m, deterministic)
}
func (m *Time) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Time.Merge(m, src)
}
func (m *Time) XXX_Size() int {
	return xxx_messageInfo_Time.Size(m)
}
func (m *Time) XXX_DiscardUnknown() {
	xxx_messageInfo_Time.DiscardUnknown(m)
}

var xxx_messageInfo_Time proto.InternalMessageInfo

func (m *Time) GetTime() float64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Time) GetNum() int64 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *Time) GetExecuteNum() int64 {
	if m != nil {
		return m.ExecuteNum
	}
	return 0
}

type Coverage struct {
	Coverage             map[uint32]uint32 `protobuf:"bytes,1,rep,name=coverage,proto3" json:"coverage,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Time                 []*Time           `protobuf:"bytes,2,rep,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Coverage) Reset()         { *m = Coverage{} }
func (m *Coverage) String() string { return proto.CompactTextString(m) }
func (*Coverage) ProtoMessage()    {}
func (*Coverage) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8fb87f92dbb7b88, []int{2}
}

func (m *Coverage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Coverage.Unmarshal(m, b)
}
func (m *Coverage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Coverage.Marshal(b, m, deterministic)
}
func (m *Coverage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Coverage.Merge(m, src)
}
func (m *Coverage) XXX_Size() int {
	return xxx_messageInfo_Coverage.Size(m)
}
func (m *Coverage) XXX_DiscardUnknown() {
	xxx_messageInfo_Coverage.DiscardUnknown(m)
}

var xxx_messageInfo_Coverage proto.InternalMessageInfo

func (m *Coverage) GetCoverage() map[uint32]uint32 {
	if m != nil {
		return m.Coverage
	}
	return nil
}

func (m *Coverage) GetTime() []*Time {
	if m != nil {
		return m.Time
	}
	return nil
}

type UsefulInput struct {
	Input                *Input   `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	Time                 float64  `protobuf:"fixed64,2,opt,name=time,proto3" json:"time,omitempty"`
	Num                  uint64   `protobuf:"varint,3,opt,name=num,proto3" json:"num,omitempty"`
	NewAddress           []uint32 `protobuf:"varint,4,rep,packed,name=new_address,json=newAddress,proto3" json:"new_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UsefulInput) Reset()         { *m = UsefulInput{} }
func (m *UsefulInput) String() string { return proto.CompactTextString(m) }
func (*UsefulInput) ProtoMessage()    {}
func (*UsefulInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8fb87f92dbb7b88, []int{3}
}

func (m *UsefulInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UsefulInput.Unmarshal(m, b)
}
func (m *UsefulInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UsefulInput.Marshal(b, m, deterministic)
}
func (m *UsefulInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UsefulInput.Merge(m, src)
}
func (m *UsefulInput) XXX_Size() int {
	return xxx_messageInfo_UsefulInput.Size(m)
}
func (m *UsefulInput) XXX_DiscardUnknown() {
	xxx_messageInfo_UsefulInput.DiscardUnknown(m)
}

var xxx_messageInfo_UsefulInput proto.InternalMessageInfo

func (m *UsefulInput) GetInput() *Input {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *UsefulInput) GetTime() float64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *UsefulInput) GetNum() uint64 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *UsefulInput) GetNewAddress() []uint32 {
	if m != nil {
		return m.NewAddress
	}
	return nil
}

type Statistics struct {
	SignalNum                 uint64               `protobuf:"varint,1,opt,name=signalNum,proto3" json:"signalNum,omitempty"`
	NumberBasicBlock          uint32               `protobuf:"varint,3,opt,name=number_basic_block,json=numberBasicBlock,proto3" json:"number_basic_block,omitempty"`
	NumberBasicBlockReal      uint32               `protobuf:"varint,4,opt,name=number_basic_block_real,json=numberBasicBlockReal,proto3" json:"number_basic_block_real,omitempty"`
	NumberBasicBlockCovered   uint32               `protobuf:"varint,5,opt,name=number_basic_block_covered,json=numberBasicBlockCovered,proto3" json:"number_basic_block_covered,omitempty"`
	NumberBasicBlockUncovered uint32               `protobuf:"varint,6,opt,name=number_basic_block_uncovered,json=numberBasicBlockUncovered,proto3" json:"number_basic_block_uncovered,omitempty"`
	Coverage                  *Coverage            `protobuf:"bytes,8,opt,name=coverage,proto3" json:"coverage,omitempty"`
	Stat                      map[int32]*Statistic `protobuf:"bytes,11,rep,name=stat,proto3" json:"stat,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UsefulInput               []*UsefulInput       `protobuf:"bytes,12,rep,name=useful_input,json=usefulInput,proto3" json:"useful_input,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}             `json:"-"`
	XXX_unrecognized          []byte               `json:"-"`
	XXX_sizecache             int32                `json:"-"`
}

func (m *Statistics) Reset()         { *m = Statistics{} }
func (m *Statistics) String() string { return proto.CompactTextString(m) }
func (*Statistics) ProtoMessage()    {}
func (*Statistics) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8fb87f92dbb7b88, []int{4}
}

func (m *Statistics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Statistics.Unmarshal(m, b)
}
func (m *Statistics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Statistics.Marshal(b, m, deterministic)
}
func (m *Statistics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Statistics.Merge(m, src)
}
func (m *Statistics) XXX_Size() int {
	return xxx_messageInfo_Statistics.Size(m)
}
func (m *Statistics) XXX_DiscardUnknown() {
	xxx_messageInfo_Statistics.DiscardUnknown(m)
}

var xxx_messageInfo_Statistics proto.InternalMessageInfo

func (m *Statistics) GetSignalNum() uint64 {
	if m != nil {
		return m.SignalNum
	}
	return 0
}

func (m *Statistics) GetNumberBasicBlock() uint32 {
	if m != nil {
		return m.NumberBasicBlock
	}
	return 0
}

func (m *Statistics) GetNumberBasicBlockReal() uint32 {
	if m != nil {
		return m.NumberBasicBlockReal
	}
	return 0
}

func (m *Statistics) GetNumberBasicBlockCovered() uint32 {
	if m != nil {
		return m.NumberBasicBlockCovered
	}
	return 0
}

func (m *Statistics) GetNumberBasicBlockUncovered() uint32 {
	if m != nil {
		return m.NumberBasicBlockUncovered
	}
	return 0
}

func (m *Statistics) GetCoverage() *Coverage {
	if m != nil {
		return m.Coverage
	}
	return nil
}

func (m *Statistics) GetStat() map[int32]*Statistic {
	if m != nil {
		return m.Stat
	}
	return nil
}

func (m *Statistics) GetUsefulInput() []*UsefulInput {
	if m != nil {
		return m.UsefulInput
	}
	return nil
}

func init() {
	proto.RegisterType((*Statistic)(nil), "dra.Statistic")
	proto.RegisterType((*Time)(nil), "dra.Time")
	proto.RegisterType((*Coverage)(nil), "dra.Coverage")
	proto.RegisterMapType((map[uint32]uint32)(nil), "dra.Coverage.CoverageEntry")
	proto.RegisterType((*UsefulInput)(nil), "dra.UsefulInput")
	proto.RegisterType((*Statistics)(nil), "dra.Statistics")
	proto.RegisterMapType((map[int32]*Statistic)(nil), "dra.Statistics.StatEntry")
}

func init() { proto.RegisterFile("Statistic.proto", fileDescriptor_d8fb87f92dbb7b88) }

var fileDescriptor_d8fb87f92dbb7b88 = []byte{
	// 537 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xd5, 0xc6, 0x4e, 0x69, 0xc6, 0x71, 0x88, 0x56, 0x95, 0xea, 0x86, 0x02, 0x56, 0x14, 0x21,
	0x23, 0x41, 0x0e, 0xa9, 0x10, 0x88, 0x1c, 0x10, 0x89, 0x00, 0x21, 0x21, 0x0e, 0x4b, 0x7b, 0xe1,
	0x12, 0x6d, 0xe2, 0x25, 0xb2, 0xea, 0x6c, 0x2a, 0xef, 0x6e, 0x42, 0xfb, 0x33, 0x7c, 0x04, 0xdf,
	0xc4, 0x7f, 0xa0, 0x9d, 0x6d, 0x12, 0xc7, 0xe9, 0x6d, 0x3c, 0xf3, 0x9e, 0xe7, 0xe9, 0xcd, 0xb3,
	0xe1, 0xf1, 0x0f, 0xcd, 0x75, 0xa6, 0x74, 0x36, 0xeb, 0xdf, 0x14, 0x4b, 0xbd, 0xa4, 0x5e, 0x5a,
	0xf0, 0x4e, 0xf0, 0x55, 0xde, 0x18, 0xed, 0x3a, 0xdd, 0xbf, 0x04, 0x1a, 0x5b, 0x14, 0xed, 0x81,
	0x2f, 0xf9, 0x42, 0x44, 0x24, 0x26, 0x49, 0x6b, 0xd0, 0xee, 0xa7, 0x05, 0xef, 0x7f, 0x36, 0x77,
	0x77, 0x99, 0x9c, 0x5b, 0x10, 0xc3, 0x29, 0x7d, 0x06, 0x20, 0x7e, 0x8b, 0x99, 0xd1, 0xe2, 0xbb,
	0x59, 0x44, 0x41, 0x4c, 0x12, 0x9f, 0x95, 0x3a, 0x94, 0x82, 0xaf, 0xb3, 0x85, 0x88, 0x9a, 0x31,
	0x49, 0x08, 0xc3, 0x9a, 0xbe, 0x80, 0x96, 0x14, 0xeb, 0x4b, 0xa1, 0xf4, 0x98, 0x2b, 0xe4, 0x85,
	0xc8, 0xab, 0x74, 0x69, 0x0f, 0x42, 0x29, 0xd6, 0x1f, 0xd3, 0xb4, 0x10, 0x4a, 0x59, 0x58, 0x0b,
	0x61, 0xfb, 0xcd, 0xee, 0x37, 0xf0, 0x2f, 0xed, 0x5b, 0x37, 0x9b, 0x48, 0x69, 0x53, 0x1b, 0x3c,
	0x69, 0x16, 0x51, 0x2d, 0x26, 0x89, 0xc7, 0x6c, 0x59, 0xd1, 0xeb, 0xe1, 0xa0, 0xd4, 0xe9, 0xfe,
	0x21, 0x70, 0x3c, 0x5e, 0xae, 0x44, 0xc1, 0xe7, 0x82, 0xbe, 0x85, 0xe3, 0xd9, 0x7d, 0x1d, 0x91,
	0xd8, 0x4b, 0x82, 0xc1, 0x13, 0xb4, 0x61, 0x03, 0xd8, 0x16, 0x9f, 0xa4, 0x2e, 0x6e, 0xd9, 0x16,
	0x4c, 0x9f, 0xde, 0x6b, 0xa9, 0x21, 0xa9, 0x81, 0x24, 0x2b, 0xd2, 0xc9, 0xea, 0x0c, 0x21, 0xdc,
	0x63, 0x5a, 0x9d, 0xd7, 0xe2, 0x16, 0xa5, 0x87, 0xcc, 0x96, 0xf4, 0x04, 0xea, 0x2b, 0x9e, 0x1b,
	0x81, 0xda, 0x43, 0xe6, 0x1e, 0xde, 0xd7, 0xde, 0x91, 0xee, 0x0a, 0x82, 0x2b, 0x25, 0x7e, 0x99,
	0x1c, 0x4f, 0x47, 0x63, 0xa8, 0x67, 0xb6, 0x40, 0x72, 0x30, 0x00, 0xdc, 0x85, 0x23, 0xe6, 0x06,
	0x5b, 0x63, 0x6a, 0x87, 0xc6, 0x78, 0x68, 0x28, 0x1a, 0xf3, 0x1c, 0x02, 0x29, 0xd6, 0x13, 0xee,
	0x8c, 0x8d, 0xfc, 0xd8, 0x4b, 0x42, 0x06, 0x3b, 0xab, 0xbb, 0xff, 0x3c, 0x80, 0x6d, 0x3a, 0x14,
	0x3d, 0x87, 0x86, 0xca, 0xe6, 0x92, 0xe7, 0xd6, 0x47, 0x82, 0xef, 0xd9, 0x35, 0xe8, 0x2b, 0xa0,
	0xd2, 0x2c, 0xa6, 0xa2, 0x98, 0x4c, 0xb9, 0xca, 0x66, 0x93, 0x69, 0xbe, 0x9c, 0x5d, 0xe3, 0xba,
	0x90, 0xb5, 0xdd, 0x64, 0x64, 0x07, 0x23, 0xdb, 0xa7, 0x6f, 0xe0, 0xf4, 0x10, 0x3d, 0x29, 0x04,
	0xcf, 0x23, 0x1f, 0x29, 0x27, 0x55, 0x0a, 0x13, 0x3c, 0xa7, 0x43, 0xe8, 0x3c, 0x40, 0xc3, 0x23,
	0x88, 0x34, 0xaa, 0x23, 0xf3, 0xb4, 0xca, 0x1c, 0xbb, 0x31, 0xfd, 0x00, 0xe7, 0x0f, 0x90, 0x8d,
	0xdc, 0xd0, 0x8f, 0x90, 0x7e, 0x56, 0xa5, 0x5f, 0x6d, 0x00, 0xf4, 0x65, 0x29, 0x1c, 0xc7, 0xe8,
	0x7d, 0xb8, 0x17, 0x8e, 0x52, 0x1c, 0x5e, 0x83, 0xaf, 0x34, 0xd7, 0x51, 0x80, 0x71, 0x38, 0x43,
	0xd8, 0xce, 0x4a, 0x2c, 0x5d, 0x82, 0x10, 0x46, 0x2f, 0xa0, 0x69, 0xf0, 0xc2, 0x13, 0x77, 0xd9,
	0x26, 0xd2, 0xdc, 0x17, 0x58, 0x3a, 0x3d, 0x0b, 0xcc, 0xee, 0xa1, 0xf3, 0xc5, 0x7d, 0xbb, 0x07,
	0x79, 0xaa, 0xbb, 0x3c, 0xf5, 0xca, 0x79, 0x0a, 0x06, 0xad, 0x7d, 0x0d, 0xa5, 0x7c, 0x8d, 0x1e,
	0xfd, 0xac, 0xf7, 0x87, 0x69, 0xc1, 0xa7, 0x47, 0xf8, 0x57, 0xb8, 0xf8, 0x1f, 0x00, 0x00, 0xff,
	0xff, 0x6f, 0xa4, 0xb5, 0x16, 0x3a, 0x04, 0x00, 0x00,
}
