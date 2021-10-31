package api

import (
	"context"
	"fmt"
	"net"
	"user-service/src/api/health"

	"google.golang.org/grpc"
)

type HealthServer struct {
	Server      *grpc.Server
	HealthChan  chan health.HealthCheckResponse_ServingStatus
	httpRouter  *HttpRouter
	grpcServer  *RpcServer
	healthState health.HealthCheckResponse_ServingStatus
	health.UnimplementedHealthServer
}

func NewHealthServer() *HealthServer {
	server := grpc.NewServer()

	healthServer := &HealthServer{
		Server:     server,
		HealthChan: make(chan health.HealthCheckResponse_ServingStatus),
	}

	server.RegisterService(&health.Health_ServiceDesc, healthServer)
	return healthServer
}

func (hs *HealthServer) Connect(addr string) {
	fmt.Printf("[Health] Listening and serving on %s\n", addr)
	rpcListener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Could not start health listener.")
	}
	hs.Server.Serve(rpcListener)
}

func (hs *HealthServer) StateServing() {
	hs.healthState = health.HealthCheckResponse_SERVING
}

func (hs *HealthServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	switch req.Service {
	case "user":
		return &health.HealthCheckResponse{
			Status: hs.healthState,
		}, nil

	default:
		return &health.HealthCheckResponse{
			Status: health.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}
}

func (hs *HealthServer) Watch(req *health.HealthCheckRequest, stream health.Health_WatchServer) error {
	go func(service string, stream health.Health_WatchServer) {
		stream.Send(&health.HealthCheckResponse{
			Status: hs.healthState,
		})

		for {
			state := <-hs.HealthChan
			if state != hs.healthState {
				stream.Send(&health.HealthCheckResponse{
					Status: state,
				})
			}
		}
	}(req.Service, stream)
	return nil
}
