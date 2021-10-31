package api

import (
	"context"
	"user-service/src/api/health"

	"google.golang.org/grpc"
)

type HealthServer struct {
	Server *grpc.Server
	health.UnimplementedHealthServer
}

func NewHealthServer() *HealthServer {
	server := grpc.NewServer()

	healthServer := &HealthServer{
		Server: server,
	}

	server.RegisterService(&health.Health_ServiceDesc, healthServer)
	return healthServer
}

func (hs *HealthServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	if req.Service != "user" {
		return &health.HealthCheckResponse{
			Status: health.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}


	return nil, nil //TODO
}

func (hs *HealthServer) Watch(*health.HealthCheckRequest, health.Health_WatchServer) error {
	return  nil //TODO
}
