package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	response, err := handler(ctx, req)
	duration := time.Since(startTime)
	statusCode := codes.Unknown
	if stat, ok := status.FromError(err); ok {
		statusCode = stat.Code()
	}
	grpcRequestLogger := log.Info()
	if err != nil {
		grpcRequestLogger = log.Error().Err(err)
	}
	grpcRequestLogger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("statusCcode", int(statusCode)).
		Str("statusText", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")
	return response, err
}
