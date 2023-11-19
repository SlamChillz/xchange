package gapi

import (
	"context"
	"net/http"
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

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body	   []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}
		next.ServeHTTP(recorder, r)
		duration := time.Since(startTime)
		httpRequestLogger := log.Info()
		if recorder.StatusCode >= http.StatusBadRequest {
			httpRequestLogger = log.Error().Bytes("body", recorder.Body)
		}
		httpRequestLogger.Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("statusCode", recorder.StatusCode).
			Str("statusText", http.StatusText(recorder.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
