// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: service_xchange.proto

package pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_xchange_proto protoreflect.FileDescriptor

var file_service_xchange_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x70, 0x63, 0x5f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x72, 0x70, 0x63, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f,
	0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x72, 0x70, 0x63, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x77, 0x61, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65,
	0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x8a, 0x02, 0x0a, 0x07, 0x58, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x12, 0x64, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2f, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x60, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x43, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x37, 0x0a, 0x08, 0x43, 0x6f, 0x69,
	0x6e, 0x53, 0x77, 0x61, 0x70, 0x12, 0x13, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x53,
	0x77, 0x61, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x62, 0x2e,
	0x43, 0x6f, 0x69, 0x6e, 0x53, 0x77, 0x61, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x83, 0x01, 0x92, 0x41, 0x5e, 0x12, 0x5c, 0x0a, 0x0b, 0x58, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x20, 0x41, 0x50, 0x49, 0x22, 0x48, 0x0a, 0x0a, 0x4d, 0x65, 0x6e, 0x64, 0x79,
	0x20, 0x53, 0x6c, 0x61, 0x6d, 0x12, 0x25, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6c, 0x61, 0x6d, 0x63, 0x68,
	0x69, 0x6c, 0x6c, 0x7a, 0x2f, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x1a, 0x13, 0x6d, 0x65,
	0x6e, 0x64, 0x79, 0x73, 0x6c, 0x61, 0x6d, 0x40, 0x67, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63, 0x6f,
	0x6d, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x6c, 0x61, 0x6d, 0x63, 0x68, 0x69, 0x6c, 0x6c, 0x7a, 0x2f, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_xchange_proto_goTypes = []interface{}{
	(*CreateCustomerRequest)(nil),  // 0: pb.CreateCustomerRequest
	(*LoginCustomerRequest)(nil),   // 1: pb.LoginCustomerRequest
	(*CoinSwapRequest)(nil),        // 2: pb.CoinSwapRequest
	(*CreateCustomerResponse)(nil), // 3: pb.CreateCustomerResponse
	(*LoginCustomerResponse)(nil),  // 4: pb.LoginCustomerResponse
	(*CoinSwapResponse)(nil),       // 5: pb.CoinSwapResponse
}
var file_service_xchange_proto_depIdxs = []int32{
	0, // 0: pb.Xchange.CreateCustomer:input_type -> pb.CreateCustomerRequest
	1, // 1: pb.Xchange.LoginCustomer:input_type -> pb.LoginCustomerRequest
	2, // 2: pb.Xchange.CoinSwap:input_type -> pb.CoinSwapRequest
	3, // 3: pb.Xchange.CreateCustomer:output_type -> pb.CreateCustomerResponse
	4, // 4: pb.Xchange.LoginCustomer:output_type -> pb.LoginCustomerResponse
	5, // 5: pb.Xchange.CoinSwap:output_type -> pb.CoinSwapResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_xchange_proto_init() }
func file_service_xchange_proto_init() {
	if File_service_xchange_proto != nil {
		return
	}
	file_rpc_create_customer_proto_init()
	file_rpc_login_customer_proto_init()
	file_rpc_create_swap_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_xchange_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_xchange_proto_goTypes,
		DependencyIndexes: file_service_xchange_proto_depIdxs,
	}.Build()
	File_service_xchange_proto = out.File
	file_service_xchange_proto_rawDesc = nil
	file_service_xchange_proto_goTypes = nil
	file_service_xchange_proto_depIdxs = nil
}
