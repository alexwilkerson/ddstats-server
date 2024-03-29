// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gamesubmission

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

// GameRecorderClient is the client API for GameRecorder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameRecorderClient interface {
	SubmitGame(ctx context.Context, in *SubmitGameRequest, opts ...grpc.CallOption) (*SubmitGameReply, error)
	ClientStart(ctx context.Context, in *ClientStartRequest, opts ...grpc.CallOption) (*ClientStartReply, error)
}

type gameRecorderClient struct {
	cc grpc.ClientConnInterface
}

func NewGameRecorderClient(cc grpc.ClientConnInterface) GameRecorderClient {
	return &gameRecorderClient{cc}
}

func (c *gameRecorderClient) SubmitGame(ctx context.Context, in *SubmitGameRequest, opts ...grpc.CallOption) (*SubmitGameReply, error) {
	out := new(SubmitGameReply)
	err := c.cc.Invoke(ctx, "/gamesubmission.GameRecorder/SubmitGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameRecorderClient) ClientStart(ctx context.Context, in *ClientStartRequest, opts ...grpc.CallOption) (*ClientStartReply, error) {
	out := new(ClientStartReply)
	err := c.cc.Invoke(ctx, "/gamesubmission.GameRecorder/ClientStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameRecorderServer is the server API for GameRecorder service.
// All implementations must embed UnimplementedGameRecorderServer
// for forward compatibility
type GameRecorderServer interface {
	SubmitGame(context.Context, *SubmitGameRequest) (*SubmitGameReply, error)
	ClientStart(context.Context, *ClientStartRequest) (*ClientStartReply, error)
	mustEmbedUnimplementedGameRecorderServer()
}

// UnimplementedGameRecorderServer must be embedded to have forward compatible implementations.
type UnimplementedGameRecorderServer struct {
}

func (UnimplementedGameRecorderServer) SubmitGame(context.Context, *SubmitGameRequest) (*SubmitGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitGame not implemented")
}
func (UnimplementedGameRecorderServer) ClientStart(context.Context, *ClientStartRequest) (*ClientStartReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientStart not implemented")
}
func (UnimplementedGameRecorderServer) mustEmbedUnimplementedGameRecorderServer() {}

// UnsafeGameRecorderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameRecorderServer will
// result in compilation errors.
type UnsafeGameRecorderServer interface {
	mustEmbedUnimplementedGameRecorderServer()
}

func RegisterGameRecorderServer(s grpc.ServiceRegistrar, srv GameRecorderServer) {
	s.RegisterService(&GameRecorder_ServiceDesc, srv)
}

func _GameRecorder_SubmitGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameRecorderServer).SubmitGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gamesubmission.GameRecorder/SubmitGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameRecorderServer).SubmitGame(ctx, req.(*SubmitGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameRecorder_ClientStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameRecorderServer).ClientStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gamesubmission.GameRecorder/ClientStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameRecorderServer).ClientStart(ctx, req.(*ClientStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GameRecorder_ServiceDesc is the grpc.ServiceDesc for GameRecorder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameRecorder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gamesubmission.GameRecorder",
	HandlerType: (*GameRecorderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitGame",
			Handler:    _GameRecorder_SubmitGame_Handler,
		},
		{
			MethodName: "ClientStart",
			Handler:    _GameRecorder_ClientStart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gamesubmission/gamesubmission.proto",
}
