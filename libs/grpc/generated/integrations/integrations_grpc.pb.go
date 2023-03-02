// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: integrations.proto

package integrations

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IntegrationsClient is the client API for Integrations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IntegrationsClient interface {
	AddIntegration(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RemoveIntegration(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type integrationsClient struct {
	cc grpc.ClientConnInterface
}

func NewIntegrationsClient(cc grpc.ClientConnInterface) IntegrationsClient {
	return &integrationsClient{cc}
}

func (c *integrationsClient) AddIntegration(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/integrations.Integrations/AddIntegration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *integrationsClient) RemoveIntegration(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/integrations.Integrations/RemoveIntegration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IntegrationsServer is the server API for Integrations service.
// All implementations must embed UnimplementedIntegrationsServer
// for forward compatibility
type IntegrationsServer interface {
	AddIntegration(context.Context, *Request) (*emptypb.Empty, error)
	RemoveIntegration(context.Context, *Request) (*emptypb.Empty, error)
	mustEmbedUnimplementedIntegrationsServer()
}

// UnimplementedIntegrationsServer must be embedded to have forward compatible implementations.
type UnimplementedIntegrationsServer struct {
}

func (UnimplementedIntegrationsServer) AddIntegration(context.Context, *Request) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIntegration not implemented")
}
func (UnimplementedIntegrationsServer) RemoveIntegration(context.Context, *Request) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveIntegration not implemented")
}
func (UnimplementedIntegrationsServer) mustEmbedUnimplementedIntegrationsServer() {}

// UnsafeIntegrationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IntegrationsServer will
// result in compilation errors.
type UnsafeIntegrationsServer interface {
	mustEmbedUnimplementedIntegrationsServer()
}

func RegisterIntegrationsServer(s grpc.ServiceRegistrar, srv IntegrationsServer) {
	s.RegisterService(&Integrations_ServiceDesc, srv)
}

func _Integrations_AddIntegration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntegrationsServer).AddIntegration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integrations.Integrations/AddIntegration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntegrationsServer).AddIntegration(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Integrations_RemoveIntegration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntegrationsServer).RemoveIntegration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/integrations.Integrations/RemoveIntegration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntegrationsServer).RemoveIntegration(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Integrations_ServiceDesc is the grpc.ServiceDesc for Integrations service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Integrations_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "integrations.Integrations",
	HandlerType: (*IntegrationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddIntegration",
			Handler:    _Integrations_AddIntegration_Handler,
		},
		{
			MethodName: "RemoveIntegration",
			Handler:    _Integrations_RemoveIntegration_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "integrations.proto",
}
