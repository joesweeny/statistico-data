// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/competition/competition.proto

package competition // import "github.com/statistico/statistico-data/proto/competition"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Competition struct {
	Id                   int64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string              `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	IsCup                *wrappers.BoolValue `protobuf:"bytes,3,opt,name=is_cup,json=isCup,proto3" json:"is_cup,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Competition) Reset()         { *m = Competition{} }
func (m *Competition) String() string { return proto.CompactTextString(m) }
func (*Competition) ProtoMessage()    {}
func (*Competition) Descriptor() ([]byte, []int) {
	return fileDescriptor_competition_a3a3a3d6c228d864, []int{0}
}
func (m *Competition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Competition.Unmarshal(m, b)
}
func (m *Competition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Competition.Marshal(b, m, deterministic)
}
func (dst *Competition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Competition.Merge(dst, src)
}
func (m *Competition) XXX_Size() int {
	return xxx_messageInfo_Competition.Size(m)
}
func (m *Competition) XXX_DiscardUnknown() {
	xxx_messageInfo_Competition.DiscardUnknown(m)
}

var xxx_messageInfo_Competition proto.InternalMessageInfo

func (m *Competition) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Competition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Competition) GetIsCup() *wrappers.BoolValue {
	if m != nil {
		return m.IsCup
	}
	return nil
}

func init() {
	proto.RegisterType((*Competition)(nil), "competition.Competition")
}

func init() {
	proto.RegisterFile("proto/competition/competition.proto", fileDescriptor_competition_a3a3a3d6c228d864)
}

var fileDescriptor_competition_a3a3a3d6c228d864 = []byte{
	// 196 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x8f, 0x31, 0x6b, 0x85, 0x30,
	0x14, 0x85, 0x89, 0xb6, 0x42, 0x23, 0x74, 0xc8, 0x24, 0x0e, 0x45, 0xda, 0xc5, 0xa5, 0x09, 0x6d,
	0x87, 0xd2, 0x55, 0xff, 0x81, 0x43, 0x87, 0x2e, 0x25, 0x26, 0xa9, 0xbd, 0xa0, 0xde, 0x60, 0x6e,
	0x78, 0x7f, 0xff, 0x41, 0xe4, 0xf1, 0x84, 0xb7, 0x7d, 0x9c, 0x73, 0xe0, 0xe3, 0xf0, 0x17, 0xbf,
	0x21, 0xa1, 0x32, 0xb8, 0x78, 0x47, 0x40, 0x80, 0xeb, 0x91, 0x65, 0x6a, 0x45, 0x79, 0x88, 0xea,
	0xa7, 0x09, 0x71, 0x9a, 0x9d, 0x4a, 0xd5, 0x18, 0xff, 0xd4, 0x69, 0xd3, 0xde, 0xbb, 0x2d, 0xec,
	0xe3, 0x67, 0xcb, 0xcb, 0xfe, 0x3a, 0x17, 0x8f, 0x3c, 0x03, 0x5b, 0xb1, 0x86, 0xb5, 0xf9, 0x90,
	0x81, 0x15, 0x82, 0xdf, 0xad, 0x7a, 0x71, 0x55, 0xd6, 0xb0, 0xf6, 0x61, 0x48, 0x2c, 0xde, 0x78,
	0x01, 0xe1, 0xd7, 0x44, 0x5f, 0xe5, 0x0d, 0x6b, 0xcb, 0xf7, 0x5a, 0xee, 0x0e, 0x79, 0x71, 0xc8,
	0x0e, 0x71, 0xfe, 0xd6, 0x73, 0x74, 0xc3, 0x3d, 0x84, 0x3e, 0xfa, 0xee, 0xeb, 0xe7, 0x73, 0x02,
	0xfa, 0x8f, 0xa3, 0x34, 0xb8, 0xa8, 0x40, 0x9a, 0x20, 0x10, 0x18, 0x3c, 0xe0, 0xab, 0xd5, 0xa4,
	0xd5, 0xcd, 0xbf, 0xb1, 0x48, 0xd1, 0xc7, 0x39, 0x00, 0x00, 0xff, 0xff, 0x5e, 0x65, 0x74, 0x74,
	0xfb, 0x00, 0x00, 0x00,
}
