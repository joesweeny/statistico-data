// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/proto/venue.proto

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

type Venue struct {
	Id                   *wrappers.Int64Value  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 *wrappers.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Venue) Reset()         { *m = Venue{} }
func (m *Venue) String() string { return proto.CompactTextString(m) }
func (*Venue) ProtoMessage()    {}
func (*Venue) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1ef5f3896975845, []int{0}
}

func (m *Venue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Venue.Unmarshal(m, b)
}
func (m *Venue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Venue.Marshal(b, m, deterministic)
}
func (m *Venue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Venue.Merge(m, src)
}
func (m *Venue) XXX_Size() int {
	return xxx_messageInfo_Venue.Size(m)
}
func (m *Venue) XXX_DiscardUnknown() {
	xxx_messageInfo_Venue.DiscardUnknown(m)
}

var xxx_messageInfo_Venue proto.InternalMessageInfo

func (m *Venue) GetId() *wrappers.Int64Value {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Venue) GetName() *wrappers.StringValue {
	if m != nil {
		return m.Name
	}
	return nil
}

func init() {
	proto.RegisterType((*Venue)(nil), "proto.Venue")
}

func init() { proto.RegisterFile("internal/app/proto/venue.proto", fileDescriptor_f1ef5f3896975845) }

var fileDescriptor_f1ef5f3896975845 = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcb, 0xcc, 0x2b, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x4f, 0x2c, 0x28, 0xd0, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f,
	0x4b, 0xcd, 0x2b, 0x4d, 0xd5, 0x03, 0xb3, 0x85, 0x58, 0xc1, 0x94, 0x94, 0x5c, 0x7a, 0x7e, 0x7e,
	0x7a, 0x4e, 0x2a, 0x44, 0x41, 0x52, 0x69, 0x9a, 0x7e, 0x79, 0x51, 0x62, 0x41, 0x41, 0x6a, 0x51,
	0x31, 0x44, 0x99, 0x52, 0x1a, 0x17, 0x6b, 0x18, 0x48, 0x97, 0x90, 0x36, 0x17, 0x53, 0x66, 0x8a,
	0x04, 0xa3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0xb4, 0x1e, 0x44, 0x97, 0x1e, 0x4c, 0x97, 0x9e, 0x67,
	0x5e, 0x89, 0x99, 0x49, 0x58, 0x62, 0x4e, 0x69, 0x6a, 0x10, 0x53, 0x66, 0x8a, 0x90, 0x01, 0x17,
	0x4b, 0x5e, 0x62, 0x6e, 0xaa, 0x04, 0x13, 0x58, 0xb9, 0x0c, 0x86, 0xf2, 0xe0, 0x92, 0xa2, 0xcc,
	0xbc, 0x74, 0x88, 0x7a, 0xb0, 0xca, 0x24, 0x36, 0xb0, 0x9c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x60, 0x54, 0x48, 0xeb, 0xb7, 0x00, 0x00, 0x00,
}
