// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: create_server.proto

package session

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

// The request message containing the book info
type TicketReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TicketNo string `protobuf:"bytes,1,opt,name=ticketNo,proto3" json:"ticketNo,omitempty"`
}

func (x *TicketReq) Reset() {
	*x = TicketReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_create_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TicketReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TicketReq) ProtoMessage() {}

func (x *TicketReq) ProtoReflect() protoreflect.Message {
	mi := &file_create_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TicketReq.ProtoReflect.Descriptor instead.
func (*TicketReq) Descriptor() ([]byte, []int) {
	return file_create_server_proto_rawDescGZIP(), []int{0}
}

func (x *TicketReq) GetTicketNo() string {
	if x != nil {
		return x.TicketNo
	}
	return ""
}

// The response message containing the info about booking
type TicketInfoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FlightDate    string `protobuf:"bytes,1,opt,name=flightDate,proto3" json:"flightDate,omitempty"`
	FlightId      string `protobuf:"bytes,2,opt,name=flightId,proto3" json:"flightId,omitempty"`
	PassengerName string `protobuf:"bytes,3,opt,name=passengerName,proto3" json:"passengerName,omitempty"`
	FlightFrom    string `protobuf:"bytes,4,opt,name=flightFrom,proto3" json:"flightFrom,omitempty"`
	FlightTo      string `protobuf:"bytes,5,opt,name=flightTo,proto3" json:"flightTo,omitempty"`
}

func (x *TicketInfoReply) Reset() {
	*x = TicketInfoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_create_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TicketInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TicketInfoReply) ProtoMessage() {}

func (x *TicketInfoReply) ProtoReflect() protoreflect.Message {
	mi := &file_create_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TicketInfoReply.ProtoReflect.Descriptor instead.
func (*TicketInfoReply) Descriptor() ([]byte, []int) {
	return file_create_server_proto_rawDescGZIP(), []int{1}
}

func (x *TicketInfoReply) GetFlightDate() string {
	if x != nil {
		return x.FlightDate
	}
	return ""
}

func (x *TicketInfoReply) GetFlightId() string {
	if x != nil {
		return x.FlightId
	}
	return ""
}

func (x *TicketInfoReply) GetPassengerName() string {
	if x != nil {
		return x.PassengerName
	}
	return ""
}

func (x *TicketInfoReply) GetFlightFrom() string {
	if x != nil {
		return x.FlightFrom
	}
	return ""
}

func (x *TicketInfoReply) GetFlightTo() string {
	if x != nil {
		return x.FlightTo
	}
	return ""
}

var File_create_server_proto protoreflect.FileDescriptor

var file_create_server_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x27,
	0x0a, 0x09, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x74,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x6f, 0x22, 0xaf, 0x01, 0x0a, 0x0f, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x66,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x44, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x70, 0x61, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1a, 0x0a,
	0x08, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x54, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x54, 0x6f, 0x32, 0x51, 0x0a, 0x0e, 0x41, 0x69, 0x72,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x0d, 0x47,
	0x65, 0x74, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x2e, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x1a, 0x18, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x14, 0x5a, 0x12,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x52,
	0x50, 0x43, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_create_server_proto_rawDescOnce sync.Once
	file_create_server_proto_rawDescData = file_create_server_proto_rawDesc
)

func file_create_server_proto_rawDescGZIP() []byte {
	file_create_server_proto_rawDescOnce.Do(func() {
		file_create_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_create_server_proto_rawDescData)
	})
	return file_create_server_proto_rawDescData
}

var file_create_server_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_create_server_proto_goTypes = []interface{}{
	(*TicketReq)(nil),       // 0: session.TicketReq
	(*TicketInfoReply)(nil), // 1: session.TicketInfoReply
}
var file_create_server_proto_depIdxs = []int32{
	0, // 0: session.AirplaneServer.GetTicketInfo:input_type -> session.TicketReq
	1, // 1: session.AirplaneServer.GetTicketInfo:output_type -> session.TicketInfoReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_create_server_proto_init() }
func file_create_server_proto_init() {
	if File_create_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_create_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TicketReq); i {
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
		file_create_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TicketInfoReply); i {
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
			RawDescriptor: file_create_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_create_server_proto_goTypes,
		DependencyIndexes: file_create_server_proto_depIdxs,
		MessageInfos:      file_create_server_proto_msgTypes,
	}.Build()
	File_create_server_proto = out.File
	file_create_server_proto_rawDesc = nil
	file_create_server_proto_goTypes = nil
	file_create_server_proto_depIdxs = nil
}
