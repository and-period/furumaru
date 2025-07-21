package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type healthServer struct{}

func RegisterHealthServer(s *grpc.Server) {
	health.RegisterHealthServer(s, &healthServer{})
}

func (s *healthServer) Check(
	_ context.Context,
	_ *health.HealthCheckRequest,
) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthServer) List(
	_ context.Context,
	_ *health.HealthListRequest,
) (*health.HealthListResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Listing is not supported")
}

func (s *healthServer) Watch(_ *health.HealthCheckRequest, _ health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
