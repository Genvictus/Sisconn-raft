// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: raft_service.proto

package raftpc

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
	RaftService_Ping_FullMethodName       = "/RaftService/Ping"
	RaftService_Get_FullMethodName        = "/RaftService/Get"
	RaftService_Set_FullMethodName        = "/RaftService/Set"
	RaftService_Strln_FullMethodName      = "/RaftService/Strln"
	RaftService_Del_FullMethodName        = "/RaftService/Del"
	RaftService_Append_FullMethodName     = "/RaftService/Append"
	RaftService_ReqLog_FullMethodName     = "/RaftService/ReqLog"
	RaftService_AddNode_FullMethodName    = "/RaftService/AddNode"
	RaftService_RemoveNode_FullMethodName = "/RaftService/RemoveNode"
)

// RaftServiceClient is the client API for RaftService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	Get(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error)
	Set(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	Strln(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error)
	Del(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error)
	Append(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	ReqLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error)
	AddNode(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	RemoveNode(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type raftServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftServiceClient(cc grpc.ClientConnInterface) RaftServiceClient {
	return &raftServiceClient{cc}
}

func (c *raftServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, RaftService_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Get(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, RaftService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Set(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, RaftService_Set_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Strln(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, RaftService_Strln_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Del(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, RaftService_Del_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Append(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, RaftService_Append_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) ReqLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error) {
	out := new(LogResponse)
	err := c.cc.Invoke(ctx, RaftService_ReqLog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) AddNode(ctx context.Context, in *KeyValuedRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, RaftService_AddNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) RemoveNode(ctx context.Context, in *KeyedRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, RaftService_RemoveNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServiceServer is the server API for RaftService service.
// All implementations must embed UnimplementedRaftServiceServer
// for forward compatibility
type RaftServiceServer interface {
	Ping(context.Context, *PingRequest) (*MessageResponse, error)
	Get(context.Context, *KeyedRequest) (*ValueResponse, error)
	Set(context.Context, *KeyValuedRequest) (*MessageResponse, error)
	Strln(context.Context, *KeyedRequest) (*ValueResponse, error)
	Del(context.Context, *KeyedRequest) (*ValueResponse, error)
	Append(context.Context, *KeyValuedRequest) (*MessageResponse, error)
	ReqLog(context.Context, *LogRequest) (*LogResponse, error)
	AddNode(context.Context, *KeyValuedRequest) (*MessageResponse, error)
	RemoveNode(context.Context, *KeyedRequest) (*MessageResponse, error)
	mustEmbedUnimplementedRaftServiceServer()
}

// UnimplementedRaftServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRaftServiceServer struct {
}

func (UnimplementedRaftServiceServer) Ping(context.Context, *PingRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedRaftServiceServer) Get(context.Context, *KeyedRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedRaftServiceServer) Set(context.Context, *KeyValuedRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedRaftServiceServer) Strln(context.Context, *KeyedRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Strln not implemented")
}
func (UnimplementedRaftServiceServer) Del(context.Context, *KeyedRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedRaftServiceServer) Append(context.Context, *KeyValuedRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedRaftServiceServer) ReqLog(context.Context, *LogRequest) (*LogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReqLog not implemented")
}
func (UnimplementedRaftServiceServer) AddNode(context.Context, *KeyValuedRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNode not implemented")
}
func (UnimplementedRaftServiceServer) RemoveNode(context.Context, *KeyedRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveNode not implemented")
}
func (UnimplementedRaftServiceServer) mustEmbedUnimplementedRaftServiceServer() {}

// UnsafeRaftServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftServiceServer will
// result in compilation errors.
type UnsafeRaftServiceServer interface {
	mustEmbedUnimplementedRaftServiceServer()
}

func RegisterRaftServiceServer(s grpc.ServiceRegistrar, srv RaftServiceServer) {
	s.RegisterService(&RaftService_ServiceDesc, srv)
}

func _RaftService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Get(ctx, req.(*KeyedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValuedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Set(ctx, req.(*KeyValuedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Strln_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Strln(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Strln_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Strln(ctx, req.(*KeyedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Del_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Del(ctx, req.(*KeyedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValuedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_Append_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Append(ctx, req.(*KeyValuedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_ReqLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).ReqLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_ReqLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).ReqLog(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValuedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_AddNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).AddNode(ctx, req.(*KeyValuedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_RemoveNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).RemoveNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftService_RemoveNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).RemoveNode(ctx, req.(*KeyedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RaftService_ServiceDesc is the grpc.ServiceDesc for RaftService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RaftService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RaftService",
	HandlerType: (*RaftServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _RaftService_Ping_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _RaftService_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _RaftService_Set_Handler,
		},
		{
			MethodName: "Strln",
			Handler:    _RaftService_Strln_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _RaftService_Del_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _RaftService_Append_Handler,
		},
		{
			MethodName: "ReqLog",
			Handler:    _RaftService_ReqLog_Handler,
		},
		{
			MethodName: "AddNode",
			Handler:    _RaftService_AddNode_Handler,
		},
		{
			MethodName: "RemoveNode",
			Handler:    _RaftService_RemoveNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft_service.proto",
}