// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: pkg/records_grpc/v1/records.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RecordsService_GetRandomRecords_FullMethodName = "/records.RecordsService/GetRandomRecords"
)

// RecordsServiceClient is the client API for RecordsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecordsServiceClient interface {
	GetRandomRecords(ctx context.Context, in *GetRandomRecordsRequest, opts ...grpc.CallOption) (*GetRandomRecordsResponse, error)
}

type recordsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecordsServiceClient(cc grpc.ClientConnInterface) RecordsServiceClient {
	return &recordsServiceClient{cc}
}

func (c *recordsServiceClient) GetRandomRecords(ctx context.Context, in *GetRandomRecordsRequest, opts ...grpc.CallOption) (*GetRandomRecordsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRandomRecordsResponse)
	err := c.cc.Invoke(ctx, RecordsService_GetRandomRecords_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecordsServiceServer is the server API for RecordsService service.
// All implementations must embed UnimplementedRecordsServiceServer
// for forward compatibility.
type RecordsServiceServer interface {
	GetRandomRecords(context.Context, *GetRandomRecordsRequest) (*GetRandomRecordsResponse, error)
	mustEmbedUnimplementedRecordsServiceServer()
}

// UnimplementedRecordsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRecordsServiceServer struct{}

func (UnimplementedRecordsServiceServer) GetRandomRecords(context.Context, *GetRandomRecordsRequest) (*GetRandomRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRandomRecords not implemented")
}
func (UnimplementedRecordsServiceServer) mustEmbedUnimplementedRecordsServiceServer() {}
func (UnimplementedRecordsServiceServer) testEmbeddedByValue()                        {}

// UnsafeRecordsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecordsServiceServer will
// result in compilation errors.
type UnsafeRecordsServiceServer interface {
	mustEmbedUnimplementedRecordsServiceServer()
}

func RegisterRecordsServiceServer(s grpc.ServiceRegistrar, srv RecordsServiceServer) {
	// If the following call pancis, it indicates UnimplementedRecordsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RecordsService_ServiceDesc, srv)
}

func _RecordsService_GetRandomRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRandomRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordsServiceServer).GetRandomRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordsService_GetRandomRecords_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordsServiceServer).GetRandomRecords(ctx, req.(*GetRandomRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecordsService_ServiceDesc is the grpc.ServiceDesc for RecordsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecordsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "records.RecordsService",
	HandlerType: (*RecordsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRandomRecords",
			Handler:    _RecordsService_GetRandomRecords_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/records_grpc/v1/records.proto",
}
