// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package DistributedMutualExclusion

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

// CriticalSectionServiceClient is the client API for CriticalSectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CriticalSectionServiceClient interface {
	Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionReply, error)
	LeaderDeclaration(ctx context.Context, in *LeaderRequest, opts ...grpc.CallOption) (*LeaderReply, error)
	QueueUp(ctx context.Context, in *CriticalSectionRequest, opts ...grpc.CallOption) (*CriticalSectionReply, error)
	GrantAccess(ctx context.Context, in *GrantAccessRequest, opts ...grpc.CallOption) (*GrantAccessReply, error)
	LeaveCriticalSection(ctx context.Context, in *LeaveCriticalSectionRequest, opts ...grpc.CallOption) (*LeaveCriticalSectionReply, error)
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error)
	UpdatePorts(ctx context.Context, in *UpdatePortsRequest, opts ...grpc.CallOption) (*UpdatePortsReply, error)
}

type criticalSectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCriticalSectionServiceClient(cc grpc.ClientConnInterface) CriticalSectionServiceClient {
	return &criticalSectionServiceClient{cc}
}

func (c *criticalSectionServiceClient) Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionReply, error) {
	out := new(ElectionReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/Election", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) LeaderDeclaration(ctx context.Context, in *LeaderRequest, opts ...grpc.CallOption) (*LeaderReply, error) {
	out := new(LeaderReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/LeaderDeclaration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) QueueUp(ctx context.Context, in *CriticalSectionRequest, opts ...grpc.CallOption) (*CriticalSectionReply, error) {
	out := new(CriticalSectionReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/QueueUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) GrantAccess(ctx context.Context, in *GrantAccessRequest, opts ...grpc.CallOption) (*GrantAccessReply, error) {
	out := new(GrantAccessReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/GrantAccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) LeaveCriticalSection(ctx context.Context, in *LeaveCriticalSectionRequest, opts ...grpc.CallOption) (*LeaveCriticalSectionReply, error) {
	out := new(LeaveCriticalSectionReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/LeaveCriticalSection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error) {
	out := new(JoinReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error) {
	out := new(LeaveReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *criticalSectionServiceClient) UpdatePorts(ctx context.Context, in *UpdatePortsRequest, opts ...grpc.CallOption) (*UpdatePortsReply, error) {
	out := new(UpdatePortsReply)
	err := c.cc.Invoke(ctx, "/ElectionOrganizer.CriticalSectionService/UpdatePorts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CriticalSectionServiceServer is the server API for CriticalSectionService service.
// All implementations must embed UnimplementedCriticalSectionServiceServer
// for forward compatibility
type CriticalSectionServiceServer interface {
	Election(context.Context, *ElectionRequest) (*ElectionReply, error)
	LeaderDeclaration(context.Context, *LeaderRequest) (*LeaderReply, error)
	QueueUp(context.Context, *CriticalSectionRequest) (*CriticalSectionReply, error)
	GrantAccess(context.Context, *GrantAccessRequest) (*GrantAccessReply, error)
	LeaveCriticalSection(context.Context, *LeaveCriticalSectionRequest) (*LeaveCriticalSectionReply, error)
	Join(context.Context, *JoinRequest) (*JoinReply, error)
	Leave(context.Context, *LeaveRequest) (*LeaveReply, error)
	UpdatePorts(context.Context, *UpdatePortsRequest) (*UpdatePortsReply, error)
	mustEmbedUnimplementedCriticalSectionServiceServer()
}

// UnimplementedCriticalSectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCriticalSectionServiceServer struct {
}

func (UnimplementedCriticalSectionServiceServer) Election(context.Context, *ElectionRequest) (*ElectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Election not implemented")
}
func (UnimplementedCriticalSectionServiceServer) LeaderDeclaration(context.Context, *LeaderRequest) (*LeaderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaderDeclaration not implemented")
}
func (UnimplementedCriticalSectionServiceServer) QueueUp(context.Context, *CriticalSectionRequest) (*CriticalSectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueueUp not implemented")
}
func (UnimplementedCriticalSectionServiceServer) GrantAccess(context.Context, *GrantAccessRequest) (*GrantAccessReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantAccess not implemented")
}
func (UnimplementedCriticalSectionServiceServer) LeaveCriticalSection(context.Context, *LeaveCriticalSectionRequest) (*LeaveCriticalSectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCriticalSection not implemented")
}
func (UnimplementedCriticalSectionServiceServer) Join(context.Context, *JoinRequest) (*JoinReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedCriticalSectionServiceServer) Leave(context.Context, *LeaveRequest) (*LeaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedCriticalSectionServiceServer) UpdatePorts(context.Context, *UpdatePortsRequest) (*UpdatePortsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePorts not implemented")
}
func (UnimplementedCriticalSectionServiceServer) mustEmbedUnimplementedCriticalSectionServiceServer() {
}

// UnsafeCriticalSectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CriticalSectionServiceServer will
// result in compilation errors.
type UnsafeCriticalSectionServiceServer interface {
	mustEmbedUnimplementedCriticalSectionServiceServer()
}

func RegisterCriticalSectionServiceServer(s grpc.ServiceRegistrar, srv CriticalSectionServiceServer) {
	s.RegisterService(&CriticalSectionService_ServiceDesc, srv)
}

func _CriticalSectionService_Election_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ElectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).Election(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/Election",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).Election(ctx, req.(*ElectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_LeaderDeclaration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).LeaderDeclaration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/LeaderDeclaration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).LeaderDeclaration(ctx, req.(*LeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_QueueUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CriticalSectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).QueueUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/QueueUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).QueueUp(ctx, req.(*CriticalSectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_GrantAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).GrantAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/GrantAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).GrantAccess(ctx, req.(*GrantAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_LeaveCriticalSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveCriticalSectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).LeaveCriticalSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/LeaveCriticalSection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).LeaveCriticalSection(ctx, req.(*LeaveCriticalSectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CriticalSectionService_UpdatePorts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePortsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CriticalSectionServiceServer).UpdatePorts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElectionOrganizer.CriticalSectionService/UpdatePorts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CriticalSectionServiceServer).UpdatePorts(ctx, req.(*UpdatePortsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CriticalSectionService_ServiceDesc is the grpc.ServiceDesc for CriticalSectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CriticalSectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ElectionOrganizer.CriticalSectionService",
	HandlerType: (*CriticalSectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Election",
			Handler:    _CriticalSectionService_Election_Handler,
		},
		{
			MethodName: "LeaderDeclaration",
			Handler:    _CriticalSectionService_LeaderDeclaration_Handler,
		},
		{
			MethodName: "QueueUp",
			Handler:    _CriticalSectionService_QueueUp_Handler,
		},
		{
			MethodName: "GrantAccess",
			Handler:    _CriticalSectionService_GrantAccess_Handler,
		},
		{
			MethodName: "LeaveCriticalSection",
			Handler:    _CriticalSectionService_LeaveCriticalSection_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _CriticalSectionService_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _CriticalSectionService_Leave_Handler,
		},
		{
			MethodName: "UpdatePorts",
			Handler:    _CriticalSectionService_UpdatePorts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "DistributedMutualExclusion/electionOrganizer.proto",
}
