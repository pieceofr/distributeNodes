// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pubsub.proto

package main

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

type NodeMessage struct {
	Command              string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Parameters           [][]byte `protobuf:"bytes,2,rep,name=Parameters,proto3" json:"Parameters,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeMessage) Reset()         { *m = NodeMessage{} }
func (m *NodeMessage) String() string { return proto.CompactTextString(m) }
func (*NodeMessage) ProtoMessage()    {}
func (*NodeMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_91df006b05e20cf7, []int{0}
}

func (m *NodeMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeMessage.Unmarshal(m, b)
}
func (m *NodeMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeMessage.Marshal(b, m, deterministic)
}
func (m *NodeMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeMessage.Merge(m, src)
}
func (m *NodeMessage) XXX_Size() int {
	return xxx_messageInfo_NodeMessage.Size(m)
}
func (m *NodeMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NodeMessage proto.InternalMessageInfo

func (m *NodeMessage) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *NodeMessage) GetParameters() [][]byte {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func init() {
	proto.RegisterType((*NodeMessage)(nil), "main.NodeMessage")
}

func init() { proto.RegisterFile("pubsub.proto", fileDescriptor_91df006b05e20cf7) }

var fileDescriptor_91df006b05e20cf7 = []byte{
	// 107 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x28, 0x4d, 0x2a,
	0x2e, 0x4d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcc, 0xcc, 0x53, 0x72,
	0xe7, 0xe2, 0xf6, 0xcb, 0x4f, 0x49, 0xf5, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x92, 0xe0,
	0x62, 0x4f, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82,
	0x71, 0x85, 0xe4, 0xb8, 0xb8, 0x02, 0x12, 0x8b, 0x12, 0x73, 0x53, 0x4b, 0x52, 0x8b, 0x8a, 0x25,
	0x98, 0x14, 0x98, 0x35, 0x78, 0x82, 0x90, 0x44, 0x92, 0xd8, 0xc0, 0xa6, 0x1a, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x18, 0x22, 0xa1, 0x9c, 0x65, 0x00, 0x00, 0x00,
}