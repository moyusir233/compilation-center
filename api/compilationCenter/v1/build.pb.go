// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: api/compilationCenter/v1/build.proto

package v1

import (
	v1 "gitee.com/moyusir/util/api/util/v1"
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

// 编译请求和响应
type BuildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 服务的对象，会作为编译缓存的标识
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// 代码生成所需的用户的注册信息
	DeviceStateRegisterInfos  []*v1.DeviceStateRegisterInfo  `protobuf:"bytes,2,rep,name=device_state_register_infos,json=deviceStateRegisterInfos,proto3" json:"device_state_register_infos,omitempty"`
	DeviceConfigRegisterInfos []*v1.DeviceConfigRegisterInfo `protobuf:"bytes,3,rep,name=device_config_register_infos,json=deviceConfigRegisterInfos,proto3" json:"device_config_register_infos,omitempty"`
}

func (x *BuildRequest) Reset() {
	*x = BuildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_compilationCenter_v1_build_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildRequest) ProtoMessage() {}

func (x *BuildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_compilationCenter_v1_build_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildRequest.ProtoReflect.Descriptor instead.
func (*BuildRequest) Descriptor() ([]byte, []int) {
	return file_api_compilationCenter_v1_build_proto_rawDescGZIP(), []int{0}
}

func (x *BuildRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BuildRequest) GetDeviceStateRegisterInfos() []*v1.DeviceStateRegisterInfo {
	if x != nil {
		return x.DeviceStateRegisterInfos
	}
	return nil
}

func (x *BuildRequest) GetDeviceConfigRegisterInfos() []*v1.DeviceConfigRegisterInfo {
	if x != nil {
		return x.DeviceConfigRegisterInfos
	}
	return nil
}

type BuildReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exe []byte `protobuf:"bytes,1,opt,name=exe,proto3" json:"exe,omitempty"`
}

func (x *BuildReply) Reset() {
	*x = BuildReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_compilationCenter_v1_build_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildReply) ProtoMessage() {}

func (x *BuildReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_compilationCenter_v1_build_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildReply.ProtoReflect.Descriptor instead.
func (*BuildReply) Descriptor() ([]byte, []int) {
	return file_api_compilationCenter_v1_build_proto_rawDescGZIP(), []int{1}
}

func (x *BuildReply) GetExe() []byte {
	if x != nil {
		return x.Exe
	}
	return nil
}

var File_api_compilationCenter_v1_build_proto protoreflect.FileDescriptor

var file_api_compilationCenter_v1_build_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x1a, 0x1e, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x2f,
	0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xf7, 0x01, 0x0a, 0x0c, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x63, 0x0a,
	0x1b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x18, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x73, 0x12, 0x66, 0x0a, 0x1c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x75,
	0x74, 0x69, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x19, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x22, 0x1e, 0x0a, 0x0a, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x65, 0x78, 0x65, 0x32, 0xed, 0x01, 0x0a, 0x05, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x12, 0x71, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f,
	0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x30, 0x01, 0x12, 0x71, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x26, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x65, 0x6e, 0x74,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x30, 0x01, 0x42, 0x71, 0x0a, 0x2b, 0x61, 0x70,
	0x69, 0x2e, 0x67, 0x69, 0x74, 0x65, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x79, 0x75,
	0x73, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d,
	0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x40, 0x67, 0x69, 0x74,
	0x65, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x79, 0x75, 0x73, 0x69, 0x72, 0x2f, 0x63,
	0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x63, 0x65, 0x6e, 0x74, 0x65,
	0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_compilationCenter_v1_build_proto_rawDescOnce sync.Once
	file_api_compilationCenter_v1_build_proto_rawDescData = file_api_compilationCenter_v1_build_proto_rawDesc
)

func file_api_compilationCenter_v1_build_proto_rawDescGZIP() []byte {
	file_api_compilationCenter_v1_build_proto_rawDescOnce.Do(func() {
		file_api_compilationCenter_v1_build_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_compilationCenter_v1_build_proto_rawDescData)
	})
	return file_api_compilationCenter_v1_build_proto_rawDescData
}

var file_api_compilationCenter_v1_build_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_compilationCenter_v1_build_proto_goTypes = []interface{}{
	(*BuildRequest)(nil),                // 0: api.compilationCenter.v1.BuildRequest
	(*BuildReply)(nil),                  // 1: api.compilationCenter.v1.BuildReply
	(*v1.DeviceStateRegisterInfo)(nil),  // 2: api.util.v1.DeviceStateRegisterInfo
	(*v1.DeviceConfigRegisterInfo)(nil), // 3: api.util.v1.DeviceConfigRegisterInfo
}
var file_api_compilationCenter_v1_build_proto_depIdxs = []int32{
	2, // 0: api.compilationCenter.v1.BuildRequest.device_state_register_infos:type_name -> api.util.v1.DeviceStateRegisterInfo
	3, // 1: api.compilationCenter.v1.BuildRequest.device_config_register_infos:type_name -> api.util.v1.DeviceConfigRegisterInfo
	0, // 2: api.compilationCenter.v1.Build.GetDataCollectionServiceProgram:input_type -> api.compilationCenter.v1.BuildRequest
	0, // 3: api.compilationCenter.v1.Build.GetDataProcessingServiceProgram:input_type -> api.compilationCenter.v1.BuildRequest
	1, // 4: api.compilationCenter.v1.Build.GetDataCollectionServiceProgram:output_type -> api.compilationCenter.v1.BuildReply
	1, // 5: api.compilationCenter.v1.Build.GetDataProcessingServiceProgram:output_type -> api.compilationCenter.v1.BuildReply
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_compilationCenter_v1_build_proto_init() }
func file_api_compilationCenter_v1_build_proto_init() {
	if File_api_compilationCenter_v1_build_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_compilationCenter_v1_build_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildRequest); i {
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
		file_api_compilationCenter_v1_build_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildReply); i {
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
			RawDescriptor: file_api_compilationCenter_v1_build_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_compilationCenter_v1_build_proto_goTypes,
		DependencyIndexes: file_api_compilationCenter_v1_build_proto_depIdxs,
		MessageInfos:      file_api_compilationCenter_v1_build_proto_msgTypes,
	}.Build()
	File_api_compilationCenter_v1_build_proto = out.File
	file_api_compilationCenter_v1_build_proto_rawDesc = nil
	file_api_compilationCenter_v1_build_proto_goTypes = nil
	file_api_compilationCenter_v1_build_proto_depIdxs = nil
}
