// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/proto/competition.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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
	return fileDescriptor_5e3aac1570e4745d, []int{0}
}

func (m *Competition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Competition.Unmarshal(m, b)
}
func (m *Competition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Competition.Marshal(b, m, deterministic)
}
func (m *Competition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Competition.Merge(m, src)
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
	proto.RegisterType((*Competition)(nil), "proto.Competition")
}

func init() {
	proto.RegisterFile("internal/app/proto/competition.proto", fileDescriptor_5e3aac1570e4745d)
}

var fileDescriptor_5e3aac1570e4745d = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8d, 0xc1, 0xaa, 0x83, 0x30,
	0x10, 0x45, 0x89, 0x3e, 0x85, 0x17, 0xa1, 0x8b, 0xac, 0xc4, 0x45, 0x91, 0xd2, 0x85, 0xab, 0x84,
	0xb6, 0x7f, 0x50, 0xff, 0xc0, 0x45, 0xb7, 0x25, 0x6a, 0x2a, 0x03, 0x31, 0x33, 0xc4, 0x84, 0xfe,
	0x7e, 0x21, 0x52, 0xba, 0x9a, 0xe1, 0xde, 0x73, 0x39, 0xfc, 0x0c, 0x2e, 0x18, 0xef, 0xb4, 0x55,
	0x9a, 0x48, 0x91, 0xc7, 0x80, 0x6a, 0xc2, 0x95, 0x4c, 0x80, 0x00, 0xe8, 0x64, 0x4a, 0x44, 0x91,
	0x4e, 0x73, 0x5c, 0x10, 0x17, 0x6b, 0x76, 0x6c, 0x8c, 0x2f, 0xf5, 0xf6, 0x9a, 0xc8, 0xf8, 0x6d,
	0xc7, 0x4e, 0x33, 0xaf, 0xfa, 0xdf, 0x56, 0x1c, 0x78, 0x06, 0x73, 0xcd, 0x5a, 0xd6, 0xe5, 0x43,
	0x06, 0xb3, 0x10, 0xfc, 0xcf, 0xe9, 0xd5, 0xd4, 0x59, 0xcb, 0xba, 0xff, 0x21, 0xfd, 0xe2, 0xc2,
	0x4b, 0xd8, 0x9e, 0x53, 0xa4, 0x3a, 0x6f, 0x59, 0x57, 0x5d, 0x1b, 0xb9, 0x3b, 0xe4, 0xd7, 0x21,
	0xef, 0x88, 0xf6, 0xa1, 0x6d, 0x34, 0x43, 0x01, 0x5b, 0x1f, 0x69, 0x2c, 0x53, 0x75, 0xfb, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xb1, 0xd4, 0x0f, 0x2b, 0xbb, 0x00, 0x00, 0x00,
}
