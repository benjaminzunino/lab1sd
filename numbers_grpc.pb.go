// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: numbers.proto

package notificaciones

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	NumberService_SendNumber_FullMethodName = "/NumberService/SendNumber"
)

// NumberServiceClient is the client API for NumberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NumberServiceClient interface {
	SendNumber(ctx context.Context, in *NumberRequest, opts ...grpc.CallOption) (*NumberResponse, error)
}

type numberServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNumberServiceClient(cc grpc.ClientConnInterface) NumberServiceClient {
	return &numberServiceClient{cc}
}

func (c *numberServiceClient) SendNumber(ctx context.Context, in *NumberRequest, opts ...grpc.CallOption) (*NumberResponse, error) {
	out := new(NumberResponse)
	err := c.cc.Invoke(ctx, NumberService_SendNumber_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NumberServiceServer is the server API for NumberService service.
// All implementations must embed UnimplementedNumberServiceServer
// for forward compatibility
type NumberServiceServer interface {
	SendNumber(context.Context, *NumberRequest) (*NumberResponse, error)
	mustEmbedUnimplementedNumberServiceServer()
}

// UnimplementedNumberServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNumberServiceServer struct {
}

func (UnimplementedNumberServiceServer) SendNumber(context.Context, *NumberRequest) (*NumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNumber not implemented")
}
func (UnimplementedNumberServiceServer) mustEmbedUnimplementedNumberServiceServer() {}

// UnsafeNumberServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NumberServiceServer will
// result in compilation errors.
type UnsafeNumberServiceServer interface {
	mustEmbedUnimplementedNumberServiceServer()
}

func RegisterNumberServiceServer(s grpc.ServiceRegistrar, srv NumberServiceServer) {
	s.RegisterService(&NumberService_ServiceDesc, srv)
}

func _NumberService_SendNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NumberServiceServer).SendNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NumberService_SendNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NumberServiceServer).SendNumber(ctx, req.(*NumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NumberService_ServiceDesc is the grpc.ServiceDesc for NumberService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NumberService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NumberService",
	HandlerType: (*NumberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendNumber",
			Handler:    _NumberService_SendNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "numbers.proto",
}
