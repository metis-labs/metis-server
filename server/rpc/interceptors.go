package rpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	if err == nil {
		fmt.Printf("RPC : %q %s", info.FullMethod, time.Since(start))
	} else {
		fmt.Printf("RPC : %q %s: %q => %q", info.FullMethod, time.Since(start), req, err)
	}

	return resp, err
}

func streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	err := handler(srv, ss)
	if err == nil {
		fmt.Printf("stream %q => ok", info.FullMethod)
	} else {
		fmt.Printf("stream %q => %s", info.FullMethod, err.Error())
	}

	return err
}
