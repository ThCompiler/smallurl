// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// ShortcutClient is the client API for Shortcut service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortcutClient interface {
	GetShortURL(ctx context.Context, in *LongUrl, opts ...grpc.CallOption) (*ShortUrl, error)
	GetLongURL(ctx context.Context, in *ShortUrl, opts ...grpc.CallOption) (*LongUrl, error)
}

type shortcutClient struct {
	cc grpc.ClientConnInterface
}

func NewShortcutClient(cc grpc.ClientConnInterface) ShortcutClient {
	return &shortcutClient{cc}
}

func (c *shortcutClient) GetShortURL(ctx context.Context, in *LongUrl, opts ...grpc.CallOption) (*ShortUrl, error) {
	out := new(ShortUrl)
	err := c.cc.Invoke(ctx, "/Shortcut/GetShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortcutClient) GetLongURL(ctx context.Context, in *ShortUrl, opts ...grpc.CallOption) (*LongUrl, error) {
	out := new(LongUrl)
	err := c.cc.Invoke(ctx, "/Shortcut/GetLongURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortcutServer is the server API for Shortcut service.
// All implementations must embed UnimplementedShortcutServer
// for forward compatibility
type ShortcutServer interface {
	GetShortURL(context.Context, *LongUrl) (*ShortUrl, error)
	GetLongURL(context.Context, *ShortUrl) (*LongUrl, error)
	mustEmbedUnimplementedShortcutServer()
}

// UnimplementedShortcutServer must be embedded to have forward compatible implementations.
type UnimplementedShortcutServer struct {
}

func (UnimplementedShortcutServer) GetShortURL(context.Context, *LongUrl) (*ShortUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortURL not implemented")
}
func (UnimplementedShortcutServer) GetLongURL(context.Context, *ShortUrl) (*LongUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLongURL not implemented")
}
func (UnimplementedShortcutServer) mustEmbedUnimplementedShortcutServer() {}

// UnsafeShortcutServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortcutServer will
// result in compilation errors.
type UnsafeShortcutServer interface {
	mustEmbedUnimplementedShortcutServer()
}

func RegisterShortcutServer(s grpc.ServiceRegistrar, srv ShortcutServer) {
	s.RegisterService(&Shortcut_ServiceDesc, srv)
}

func _Shortcut_GetShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LongUrl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortcutServer).GetShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Shortcut/GetShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortcutServer).GetShortURL(ctx, req.(*LongUrl))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortcut_GetLongURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortUrl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortcutServer).GetLongURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Shortcut/GetLongURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortcutServer).GetLongURL(ctx, req.(*ShortUrl))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortcut_ServiceDesc is the grpc.ServiceDesc for Shortcut service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortcut_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Shortcut",
	HandlerType: (*ShortcutServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetShortURL",
			Handler:    _Shortcut_GetShortURL_Handler,
		},
		{
			MethodName: "GetLongURL",
			Handler:    _Shortcut_GetLongURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/grpc/v1/api.proto",
}