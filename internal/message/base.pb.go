// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        (unknown)
// source: base.proto

package message

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// 战斗相关属性
type FIGHT_ATTR int32

const (
	FIGHT_ATTR_UNKNOWN FIGHT_ATTR = 0
	FIGHT_ATTR_HP      FIGHT_ATTR = 1 // 血
	FIGHT_ATTR_ATT     FIGHT_ATTR = 2 // 攻
)

// Enum value maps for FIGHT_ATTR.
var (
	FIGHT_ATTR_name = map[int32]string{
		0: "UNKNOWN",
		1: "HP",
		2: "ATT",
	}
	FIGHT_ATTR_value = map[string]int32{
		"UNKNOWN": 0,
		"HP":      1,
		"ATT":     2,
	}
)

func (x FIGHT_ATTR) Enum() *FIGHT_ATTR {
	p := new(FIGHT_ATTR)
	*p = x
	return p
}

func (x FIGHT_ATTR) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FIGHT_ATTR) Descriptor() protoreflect.EnumDescriptor {
	return file_base_proto_enumTypes[0].Descriptor()
}

func (FIGHT_ATTR) Type() protoreflect.EnumType {
	return &file_base_proto_enumTypes[0]
}

func (x FIGHT_ATTR) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FIGHT_ATTR.Descriptor instead.
func (FIGHT_ATTR) EnumDescriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

type Vec3F struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float64 `protobuf:"fixed64,1,opt,name=X,json=x,proto3" json:"X"`
	Y float64 `protobuf:"fixed64,2,opt,name=Y,json=y,proto3" json:"Y"`
	Z float64 `protobuf:"fixed64,3,opt,name=Z,json=z,proto3" json:"Z"`
}

func (x *Vec3F) Reset() {
	*x = Vec3F{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vec3F) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vec3F) ProtoMessage() {}

func (x *Vec3F) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vec3F.ProtoReflect.Descriptor instead.
func (*Vec3F) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

func (x *Vec3F) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Vec3F) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Vec3F) GetZ() float64 {
	if x != nil {
		return x.Z
	}
	return 0
}

var File_base_proto protoreflect.FileDescriptor

var file_base_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x31, 0x0a, 0x05, 0x56, 0x65, 0x63, 0x33, 0x46, 0x12, 0x0c,
	0x0a, 0x01, 0x58, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01,
	0x59, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x5a, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x7a, 0x2a, 0x2a, 0x0a, 0x0a, 0x46, 0x49, 0x47, 0x48,
	0x54, 0x5f, 0x41, 0x54, 0x54, 0x52, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x48, 0x50, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x41,
	0x54, 0x54, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_base_proto_rawDescOnce sync.Once
	file_base_proto_rawDescData = file_base_proto_rawDesc
)

func file_base_proto_rawDescGZIP() []byte {
	file_base_proto_rawDescOnce.Do(func() {
		file_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_base_proto_rawDescData)
	})
	return file_base_proto_rawDescData
}

var file_base_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_base_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_base_proto_goTypes = []interface{}{
	(FIGHT_ATTR)(0), // 0: message.FIGHT_ATTR
	(*Vec3F)(nil),   // 1: message.Vec3F
}
var file_base_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_base_proto_init() }
func file_base_proto_init() {
	if File_base_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vec3F); i {
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
			RawDescriptor: file_base_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_base_proto_goTypes,
		DependencyIndexes: file_base_proto_depIdxs,
		EnumInfos:         file_base_proto_enumTypes,
		MessageInfos:      file_base_proto_msgTypes,
	}.Build()
	File_base_proto = out.File
	file_base_proto_rawDesc = nil
	file_base_proto_goTypes = nil
	file_base_proto_depIdxs = nil
}
