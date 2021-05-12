package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO(youngteac.hong): We need to change below userID with accessToken.
func unaryInterceptor(
	userID string,
) func(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(attachToken(ctx, userID), method, req, reply, cc, opts...)
	}
}

func attachToken(ctx context.Context, userID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", userID)
}
