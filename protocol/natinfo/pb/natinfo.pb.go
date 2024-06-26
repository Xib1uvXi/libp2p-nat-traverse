// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: pb/natinfo.proto

package pb

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

type MsgType int32

const (
	MsgType_TestNatType                  MsgType = 0
	MsgType_PortNegotiation              MsgType = 1
	MsgType_PortNegotiationResponse      MsgType = 2
	MsgType_ServerPortChangeTest         MsgType = 3
	MsgType_ServerPortChangeTestResponse MsgType = 4
	MsgType_NatTypeResult                MsgType = 5
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "TestNatType",
		1: "PortNegotiation",
		2: "PortNegotiationResponse",
		3: "ServerPortChangeTest",
		4: "ServerPortChangeTestResponse",
		5: "NatTypeResult",
	}
	MsgType_value = map[string]int32{
		"TestNatType":                  0,
		"PortNegotiation":              1,
		"PortNegotiationResponse":      2,
		"ServerPortChangeTest":         3,
		"ServerPortChangeTestResponse": 4,
		"NatTypeResult":                5,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_natinfo_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_pb_natinfo_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_pb_natinfo_proto_rawDescGZIP(), []int{0}
}

type NATType int32

const (
	NATType_Unknown              NATType = 0
	NATType_None                 NATType = 1
	NATType_FullOrRestrictedCone NATType = 2
	NATType_PortRestrictedCone   NATType = 3
	NATType_Symmetric            NATType = 4
)

// Enum value maps for NATType.
var (
	NATType_name = map[int32]string{
		0: "Unknown",
		1: "None",
		2: "FullOrRestrictedCone",
		3: "PortRestrictedCone",
		4: "Symmetric",
	}
	NATType_value = map[string]int32{
		"Unknown":              0,
		"None":                 1,
		"FullOrRestrictedCone": 2,
		"PortRestrictedCone":   3,
		"Symmetric":            4,
	}
)

func (x NATType) Enum() *NATType {
	p := new(NATType)
	*p = x
	return p
}

func (x NATType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NATType) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_natinfo_proto_enumTypes[1].Descriptor()
}

func (NATType) Type() protoreflect.EnumType {
	return &file_pb_natinfo_proto_enumTypes[1]
}

func (x NATType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NATType.Descriptor instead.
func (NATType) EnumDescriptor() ([]byte, []int) {
	return file_pb_natinfo_proto_rawDescGZIP(), []int{1}
}

type PortChangeType int32

const (
	PortChangeType_Linear      PortChangeType = 0
	PortChangeType_Random      PortChangeType = 1
	PortChangeType_UnKnownRule PortChangeType = 2
)

// Enum value maps for PortChangeType.
var (
	PortChangeType_name = map[int32]string{
		0: "Linear",
		1: "Random",
		2: "UnKnownRule",
	}
	PortChangeType_value = map[string]int32{
		"Linear":      0,
		"Random":      1,
		"UnKnownRule": 2,
	}
)

func (x PortChangeType) Enum() *PortChangeType {
	p := new(PortChangeType)
	*p = x
	return p
}

func (x PortChangeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PortChangeType) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_natinfo_proto_enumTypes[2].Descriptor()
}

func (PortChangeType) Type() protoreflect.EnumType {
	return &file_pb_natinfo_proto_enumTypes[2]
}

func (x PortChangeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PortChangeType.Descriptor instead.
func (PortChangeType) EnumDescriptor() ([]byte, []int) {
	return file_pb_natinfo_proto_rawDescGZIP(), []int{2}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=natinfo.pb.MsgType" json:"type,omitempty"`
	Identity      string  `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity,omitempty"`
	Data          []byte  `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	ErrorInfo     string  `protobuf:"bytes,4,opt,name=error_info,json=errorInfo,proto3" json:"error_info,omitempty"`
	SrcPublicAddr string  `protobuf:"bytes,5,opt,name=src_public_addr,json=srcPublicAddr,proto3" json:"src_public_addr,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_natinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_pb_natinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_pb_natinfo_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_TestNatType
}

func (x *Message) GetIdentity() string {
	if x != nil {
		return x.Identity
	}
	return ""
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Message) GetErrorInfo() string {
	if x != nil {
		return x.ErrorInfo
	}
	return ""
}

func (x *Message) GetSrcPublicAddr() string {
	if x != nil {
		return x.SrcPublicAddr
	}
	return ""
}

type NATTypeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NatType                  NATType        `protobuf:"varint,1,opt,name=nat_type,json=natType,proto3,enum=natinfo.pb.NATType" json:"nat_type,omitempty"`
	PortInfluencedByProtocol bool           `protobuf:"varint,2,opt,name=port_influenced_by_protocol,json=portInfluencedByProtocol,proto3" json:"port_influenced_by_protocol,omitempty"`
	UdpPortChangeRule        PortChangeType `protobuf:"varint,3,opt,name=udp_port_change_rule,json=udpPortChangeRule,proto3,enum=natinfo.pb.PortChangeType" json:"udp_port_change_rule,omitempty"`
}

func (x *NATTypeInfo) Reset() {
	*x = NATTypeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_natinfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NATTypeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NATTypeInfo) ProtoMessage() {}

func (x *NATTypeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pb_natinfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NATTypeInfo.ProtoReflect.Descriptor instead.
func (*NATTypeInfo) Descriptor() ([]byte, []int) {
	return file_pb_natinfo_proto_rawDescGZIP(), []int{1}
}

func (x *NATTypeInfo) GetNatType() NATType {
	if x != nil {
		return x.NatType
	}
	return NATType_Unknown
}

func (x *NATTypeInfo) GetPortInfluencedByProtocol() bool {
	if x != nil {
		return x.PortInfluencedByProtocol
	}
	return false
}

func (x *NATTypeInfo) GetUdpPortChangeRule() PortChangeType {
	if x != nil {
		return x.UdpPortChangeRule
	}
	return PortChangeType_Linear
}

var File_pb_natinfo_proto protoreflect.FileDescriptor

var file_pb_natinfo_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x62, 0x2f, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x62, 0x22, 0xa9,
	0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x26, 0x0a, 0x0f, 0x73, 0x72, 0x63, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x72, 0x63,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x64, 0x64, 0x72, 0x22, 0xc9, 0x01, 0x0a, 0x0b, 0x4e,
	0x41, 0x54, 0x54, 0x79, 0x70, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2e, 0x0a, 0x08, 0x6e, 0x61,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6e,
	0x61, 0x74, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x62, 0x2e, 0x4e, 0x41, 0x54, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x07, 0x6e, 0x61, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3d, 0x0a, 0x1b, 0x70, 0x6f,
	0x72, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x5f, 0x62, 0x79,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x18, 0x70, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6c, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x64, 0x42,
	0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x4b, 0x0a, 0x14, 0x75, 0x64, 0x70,
	0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x75, 0x6c,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x11, 0x75, 0x64, 0x70, 0x50, 0x6f, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x2a, 0x9b, 0x01, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x4e, 0x61, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x6f, 0x72, 0x74, 0x4e, 0x65, 0x67, 0x6f, 0x74,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x50, 0x6f, 0x72, 0x74,
	0x4e, 0x65, 0x67, 0x6f, 0x74, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50,
	0x6f, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x65, 0x73, 0x74, 0x10, 0x03, 0x12,
	0x20, 0x0a, 0x1c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x10,
	0x04, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x61, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x10, 0x05, 0x2a, 0x61, 0x0a, 0x07, 0x4e, 0x41, 0x54, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x46, 0x75, 0x6c, 0x6c, 0x4f, 0x72,
	0x52, 0x65, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x65, 0x10, 0x02,
	0x12, 0x16, 0x0a, 0x12, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74,
	0x65, 0x64, 0x43, 0x6f, 0x6e, 0x65, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x79, 0x6d, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x04, 0x2a, 0x39, 0x0a, 0x0e, 0x50, 0x6f, 0x72, 0x74, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x69, 0x6e,
	0x65, 0x61, 0x72, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x10,
	0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x4b, 0x6e, 0x6f, 0x77, 0x6e, 0x52, 0x75, 0x6c, 0x65,
	0x10, 0x02, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_natinfo_proto_rawDescOnce sync.Once
	file_pb_natinfo_proto_rawDescData = file_pb_natinfo_proto_rawDesc
)

func file_pb_natinfo_proto_rawDescGZIP() []byte {
	file_pb_natinfo_proto_rawDescOnce.Do(func() {
		file_pb_natinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_natinfo_proto_rawDescData)
	})
	return file_pb_natinfo_proto_rawDescData
}

var file_pb_natinfo_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_pb_natinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_natinfo_proto_goTypes = []interface{}{
	(MsgType)(0),        // 0: natinfo.pb.MsgType
	(NATType)(0),        // 1: natinfo.pb.NATType
	(PortChangeType)(0), // 2: natinfo.pb.PortChangeType
	(*Message)(nil),     // 3: natinfo.pb.Message
	(*NATTypeInfo)(nil), // 4: natinfo.pb.NATTypeInfo
}
var file_pb_natinfo_proto_depIdxs = []int32{
	0, // 0: natinfo.pb.Message.type:type_name -> natinfo.pb.MsgType
	1, // 1: natinfo.pb.NATTypeInfo.nat_type:type_name -> natinfo.pb.NATType
	2, // 2: natinfo.pb.NATTypeInfo.udp_port_change_rule:type_name -> natinfo.pb.PortChangeType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pb_natinfo_proto_init() }
func file_pb_natinfo_proto_init() {
	if File_pb_natinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_natinfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_pb_natinfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NATTypeInfo); i {
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
			RawDescriptor: file_pb_natinfo_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_natinfo_proto_goTypes,
		DependencyIndexes: file_pb_natinfo_proto_depIdxs,
		EnumInfos:         file_pb_natinfo_proto_enumTypes,
		MessageInfos:      file_pb_natinfo_proto_msgTypes,
	}.Build()
	File_pb_natinfo_proto = out.File
	file_pb_natinfo_proto_rawDesc = nil
	file_pb_natinfo_proto_goTypes = nil
	file_pb_natinfo_proto_depIdxs = nil
}
