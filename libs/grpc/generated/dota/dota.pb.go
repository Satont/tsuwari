// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: dota.proto

package dota

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

type GetPlayerCardRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId int64 `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
}

func (x *GetPlayerCardRequest) Reset() {
	*x = GetPlayerCardRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dota_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPlayerCardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPlayerCardRequest) ProtoMessage() {}

func (x *GetPlayerCardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dota_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPlayerCardRequest.ProtoReflect.Descriptor instead.
func (*GetPlayerCardRequest) Descriptor() ([]byte, []int) {
	return file_dota_proto_rawDescGZIP(), []int{0}
}

func (x *GetPlayerCardRequest) GetAccountId() int64 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

type GetPlayerCardResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId       string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	RankTier        *int64 `protobuf:"varint,2,opt,name=rank_tier,json=rankTier,proto3,oneof" json:"rank_tier,omitempty"`
	LeaderboardRank *int64 `protobuf:"varint,3,opt,name=leaderboard_rank,json=leaderboardRank,proto3,oneof" json:"leaderboard_rank,omitempty"`
}

func (x *GetPlayerCardResponse) Reset() {
	*x = GetPlayerCardResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dota_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPlayerCardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPlayerCardResponse) ProtoMessage() {}

func (x *GetPlayerCardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dota_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPlayerCardResponse.ProtoReflect.Descriptor instead.
func (*GetPlayerCardResponse) Descriptor() ([]byte, []int) {
	return file_dota_proto_rawDescGZIP(), []int{1}
}

func (x *GetPlayerCardResponse) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *GetPlayerCardResponse) GetRankTier() int64 {
	if x != nil && x.RankTier != nil {
		return *x.RankTier
	}
	return 0
}

func (x *GetPlayerCardResponse) GetLeaderboardRank() int64 {
	if x != nil && x.LeaderboardRank != nil {
		return *x.LeaderboardRank
	}
	return 0
}

var File_dota_proto protoreflect.FileDescriptor

var file_dota_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x6f, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x6f,
	0x74, 0x61, 0x22, 0x34, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43,
	0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x22, 0xab, 0x01, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x09, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x74, 0x69, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x72, 0x61, 0x6e, 0x6b, 0x54, 0x69, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x10, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52,
	0x0f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x61, 0x6e, 0x6b,
	0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x74, 0x69, 0x65,
	0x72, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x32, 0x52, 0x0a, 0x04, 0x44, 0x6f, 0x74, 0x61, 0x12, 0x4a,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x61, 0x72, 0x64, 0x12,
	0x1a, 0x2e, 0x64, 0x6f, 0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x43, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x6f,
	0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x61, 0x72, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x13, 0x5a, 0x11, 0x74, 0x73,
	0x75, 0x77, 0x61, 0x72, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x64, 0x6f, 0x74, 0x61, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dota_proto_rawDescOnce sync.Once
	file_dota_proto_rawDescData = file_dota_proto_rawDesc
)

func file_dota_proto_rawDescGZIP() []byte {
	file_dota_proto_rawDescOnce.Do(func() {
		file_dota_proto_rawDescData = protoimpl.X.CompressGZIP(file_dota_proto_rawDescData)
	})
	return file_dota_proto_rawDescData
}

var file_dota_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_dota_proto_goTypes = []interface{}{
	(*GetPlayerCardRequest)(nil),  // 0: dota.GetPlayerCardRequest
	(*GetPlayerCardResponse)(nil), // 1: dota.GetPlayerCardResponse
}
var file_dota_proto_depIdxs = []int32{
	0, // 0: dota.Dota.GetPlayerCard:input_type -> dota.GetPlayerCardRequest
	1, // 1: dota.Dota.GetPlayerCard:output_type -> dota.GetPlayerCardResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dota_proto_init() }
func file_dota_proto_init() {
	if File_dota_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dota_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPlayerCardRequest); i {
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
		file_dota_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPlayerCardResponse); i {
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
	file_dota_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dota_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dota_proto_goTypes,
		DependencyIndexes: file_dota_proto_depIdxs,
		MessageInfos:      file_dota_proto_msgTypes,
	}.Build()
	File_dota_proto = out.File
	file_dota_proto_rawDesc = nil
	file_dota_proto_goTypes = nil
	file_dota_proto_depIdxs = nil
}
