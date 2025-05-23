// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.3
// source: cart.proto

package cart

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

// CartServiceClient is the client API for CartService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CartServiceClient interface {
	DeleteCartByUserID(ctx context.Context, in *DeleteCartByUserIdRequest, opts ...grpc.CallOption) (*DeleteCartByUserIdResponse, error)
}

type cartServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCartServiceClient(cc grpc.ClientConnInterface) CartServiceClient {
	return &cartServiceClient{cc}
}

func (c *cartServiceClient) DeleteCartByUserID(ctx context.Context, in *DeleteCartByUserIdRequest, opts ...grpc.CallOption) (*DeleteCartByUserIdResponse, error) {
	out := new(DeleteCartByUserIdResponse)
	err := c.cc.Invoke(ctx, "/cart.CartService/DeleteCartByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CartServiceServer is the server API for CartService service.
// All implementations must embed UnimplementedCartServiceServer
// for forward compatibility
type CartServiceServer interface {
	DeleteCartByUserID(context.Context, *DeleteCartByUserIdRequest) (*DeleteCartByUserIdResponse, error)
	mustEmbedUnimplementedCartServiceServer()
}

// UnimplementedCartServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCartServiceServer struct {
}

func (UnimplementedCartServiceServer) DeleteCartByUserID(context.Context, *DeleteCartByUserIdRequest) (*DeleteCartByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCartByUserID not implemented")
}
func (UnimplementedCartServiceServer) mustEmbedUnimplementedCartServiceServer() {}

// UnsafeCartServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CartServiceServer will
// result in compilation errors.
type UnsafeCartServiceServer interface {
	mustEmbedUnimplementedCartServiceServer()
}

func RegisterCartServiceServer(s grpc.ServiceRegistrar, srv CartServiceServer) {
	s.RegisterService(&CartService_ServiceDesc, srv)
}

func _CartService_DeleteCartByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCartByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).DeleteCartByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cart.CartService/DeleteCartByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).DeleteCartByUserID(ctx, req.(*DeleteCartByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CartService_ServiceDesc is the grpc.ServiceDesc for CartService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CartService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cart.CartService",
	HandlerType: (*CartServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteCartByUserID",
			Handler:    _CartService_DeleteCartByUserID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cart.proto",
}
