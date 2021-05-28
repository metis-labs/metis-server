package rpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server/database"
	"oss.navercorp.com/metis/metis-server/server/types"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// TODO(hackerwins): do authenticate only against authMethods
	ctx, err := authenticate(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := handler(ctx, req)
	if err == nil {
		log.Logger.Infof("RPC : %q %s", info.FullMethod, time.Since(start))
	} else {
		err = toStatusError(err)
		log.Logger.Warnf("RPC : %q %s: %q => %q", info.FullMethod, time.Since(start), req, err)
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
		log.Logger.Infof("stream %q => ok", info.FullMethod)
	} else {
		log.Logger.Warn("stream %q => %s", info.FullMethod, err.Error())
	}

	return err
}

func authenticate(ctx context.Context) (context.Context, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := data["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	userID := values[0]
	if len(userID) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	return types.CtxWithUserID(ctx, userID), nil
}

// toStatusError returns a status.Error from the given logic error. If an error
// occurs while executing logic in API handler, gRPC status.error should be
// returned so that the client can know more about the status of the request.
func toStatusError(err error) error {
	if errors.Is(err, database.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, database.ErrInvalidID) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
