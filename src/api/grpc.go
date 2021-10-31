package api

import (
	"context"
	"fmt"
	"net"
	"user-service/src/api/generated"
	"user-service/src/api/health"
	"user-service/src/models"

	"google.golang.org/grpc"
)

type RpcServer struct {
	Server *grpc.Server
	us *UserService
	generated.UnimplementedUserServiceServer
	healthChannel chan health.HealthCheckResponse_ServingStatus
}

func NewRpcServer(userService *UserService, healthChan chan health.HealthCheckResponse_ServingStatus) *RpcServer {
	server := grpc.NewServer()

	rpcServer := &RpcServer{
		Server: server,
		us: userService,
		healthChannel: healthChan,
	}

	server.RegisterService(&generated.UserService_ServiceDesc, rpcServer)
	return rpcServer
}

func (rs *RpcServer) Connect(addr string) {
	fmt.Printf("[gRPC] Listening and serving on %s\n", addr)
	rpcListener, err := net.Listen("tcp", addr)
	if err != nil {
		rs.healthChannel <- health.HealthCheckResponse_NOT_SERVING
		panic("Could not start gRPC listener.")
	}
	err = rs.Server.Serve(rpcListener)
	if err != nil {
		rs.healthChannel <- health.HealthCheckResponse_NOT_SERVING
		panic("Could not start gRPC listener.")
	}
}

func (rs *RpcServer) CreateUser(ctx context.Context, req *generated.CreateUserReq) (*generated.CreateUserResp, error) {
	user, err := rs.us.CreateUser((*models.User)(req.User))
	if err != nil {
		return nil, err
	}

	return &generated.CreateUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) UpdateUser(ctx context.Context, req *generated.UpdateUserReq) (*generated.UpdateUserResp, error) {
	user, err := rs.us.UpdateUser((*models.User)(req.User))
	if err != nil {
		return nil, err
	}

	return &generated.UpdateUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) DeleteUser(ctx context.Context, req *generated.DeleteUserReq) (*generated.DeleteUserResp, error) {
	err := rs.us.RemoveUser(req.UserId)
	if err != nil {
		return nil, err
	}

	return &generated.DeleteUserResp{}, nil
}

func (rs *RpcServer) GetUser(ctx context.Context, req *generated.GetUserReq) (*generated.GetUserResp, error) {
	user, err := rs.us.GetUser(req.UserId)
	if err != nil {
		return nil, err
	}

	return &generated.GetUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) GetUserList(req *generated.GetUserListReq, stream generated.UserService_GetUserListServer) (error) {
	users, err := rs.us.GetUserList(int(req.Limit), int(req.Skip), req.Filter)
	if err != nil {
		return err
	}

	for _, user := range(users) {
		err = stream.Send(&generated.GetUserListResp{
			User: (*generated.User)(user),
		})
	}

	return nil
}
