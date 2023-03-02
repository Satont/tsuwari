// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: tokens.proto

package tokens

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

// TokensClient is the client API for Tokens service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokensClient interface {
	RequestAppToken(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Token, error)
	RequestUserToken(ctx context.Context, in *GetUserTokenRequest, opts ...grpc.CallOption) (*Token, error)
	RequestBotToken(ctx context.Context, in *GetBotTokenRequest, opts ...grpc.CallOption) (*Token, error)
	UpdateBotToken(ctx context.Context, in *UpdateToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateUserToken(ctx context.Context, in *UpdateToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type tokensClient struct {
	cc grpc.ClientConnInterface
}

func NewTokensClient(cc grpc.ClientConnInterface) TokensClient {
	return &tokensClient{cc}
}

func (c *tokensClient) RequestAppToken(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/tokens.Tokens/RequestAppToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokensClient) RequestUserToken(ctx context.Context, in *GetUserTokenRequest, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/tokens.Tokens/RequestUserToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokensClient) RequestBotToken(ctx context.Context, in *GetBotTokenRequest, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/tokens.Tokens/RequestBotToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokensClient) UpdateBotToken(ctx context.Context, in *UpdateToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/tokens.Tokens/UpdateBotToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokensClient) UpdateUserToken(ctx context.Context, in *UpdateToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/tokens.Tokens/UpdateUserToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokensServer is the server API for Tokens service.
// All implementations must embed UnimplementedTokensServer
// for forward compatibility
type TokensServer interface {
	RequestAppToken(context.Context, *emptypb.Empty) (*Token, error)
	RequestUserToken(context.Context, *GetUserTokenRequest) (*Token, error)
	RequestBotToken(context.Context, *GetBotTokenRequest) (*Token, error)
	UpdateBotToken(context.Context, *UpdateToken) (*emptypb.Empty, error)
	UpdateUserToken(context.Context, *UpdateToken) (*emptypb.Empty, error)
	mustEmbedUnimplementedTokensServer()
}

// UnimplementedTokensServer must be embedded to have forward compatible implementations.
type UnimplementedTokensServer struct {
}

func (UnimplementedTokensServer) RequestAppToken(context.Context, *emptypb.Empty) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAppToken not implemented")
}
func (UnimplementedTokensServer) RequestUserToken(context.Context, *GetUserTokenRequest) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestUserToken not implemented")
}
func (UnimplementedTokensServer) RequestBotToken(context.Context, *GetBotTokenRequest) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestBotToken not implemented")
}
func (UnimplementedTokensServer) UpdateBotToken(context.Context, *UpdateToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBotToken not implemented")
}
func (UnimplementedTokensServer) UpdateUserToken(context.Context, *UpdateToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserToken not implemented")
}
func (UnimplementedTokensServer) mustEmbedUnimplementedTokensServer() {}

// UnsafeTokensServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokensServer will
// result in compilation errors.
type UnsafeTokensServer interface {
	mustEmbedUnimplementedTokensServer()
}

func RegisterTokensServer(s grpc.ServiceRegistrar, srv TokensServer) {
	s.RegisterService(&Tokens_ServiceDesc, srv)
}

func _Tokens_RequestAppToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServer).RequestAppToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokens.Tokens/RequestAppToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServer).RequestAppToken(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tokens_RequestUserToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServer).RequestUserToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokens.Tokens/RequestUserToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServer).RequestUserToken(ctx, req.(*GetUserTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tokens_RequestBotToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBotTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServer).RequestBotToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokens.Tokens/RequestBotToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServer).RequestBotToken(ctx, req.(*GetBotTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tokens_UpdateBotToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServer).UpdateBotToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokens.Tokens/UpdateBotToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServer).UpdateBotToken(ctx, req.(*UpdateToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tokens_UpdateUserToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokensServer).UpdateUserToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokens.Tokens/UpdateUserToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokensServer).UpdateUserToken(ctx, req.(*UpdateToken))
	}
	return interceptor(ctx, in, info, handler)
}

// Tokens_ServiceDesc is the grpc.ServiceDesc for Tokens service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tokens_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tokens.Tokens",
	HandlerType: (*TokensServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestAppToken",
			Handler:    _Tokens_RequestAppToken_Handler,
		},
		{
			MethodName: "RequestUserToken",
			Handler:    _Tokens_RequestUserToken_Handler,
		},
		{
			MethodName: "RequestBotToken",
			Handler:    _Tokens_RequestBotToken_Handler,
		},
		{
			MethodName: "UpdateBotToken",
			Handler:    _Tokens_UpdateBotToken_Handler,
		},
		{
			MethodName: "UpdateUserToken",
			Handler:    _Tokens_UpdateUserToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokens.proto",
}
