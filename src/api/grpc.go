package api

import (
	"context"
	"user-service/src/api/generated"
	"user-service/src/models"
	"user-service/src/storage"

	"google.golang.org/grpc"
)

type RpcServer struct {
	Server *grpc.Server
	db storage.Database
	generated.UnimplementedUserServiceServer
}

func NewRpcServer(db storage.Database) *RpcServer {
	server := grpc.NewServer()

	rpcServer := &RpcServer{
		Server: server,
		db:  db,
	}

	server.RegisterService(&generated.UserService_ServiceDesc, rpcServer)
	return rpcServer
}

func (rs *RpcServer) CreateUser(ctx context.Context, req *generated.CreateUserReq) (*generated.CreateUserResp, error) {
	user, err := CreateUser((*models.User)(req.User), rs.db)
	if err != nil {
		return nil, err
	}

	return &generated.CreateUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) UpdateUser(ctx context.Context, req *generated.UpdateUserReq) (*generated.UpdateUserResp, error) {
	user, err := UpdateUser((*models.User)(req.User), rs.db)
	if err != nil {
		return nil, err
	}

	return &generated.UpdateUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) DeleteUser(ctx context.Context, req *generated.DeleteUserReq) (*generated.DeleteUserResp, error) {
	err := RemoveUser(req.UserId, rs.db)
	if err != nil {
		return nil, err
	}

	return &generated.DeleteUserResp{}, nil
}

func (rs *RpcServer) GetUser(ctx context.Context, req *generated.GetUserReq) (*generated.GetUserResp, error) {
	user, err := GetUser(req.UserId, rs.db)
	if err != nil {
		return nil, err
	}

	return &generated.GetUserResp{
		User: (*generated.User)(user),
	}, nil
}

func (rs *RpcServer) GetUserList(req *generated.GetUserListReq, stream generated.UserService_GetUserListServer) (error) {
	users, err := GetUsers(int(req.Limit), int(req.Skip), req.Filter, rs.db)
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
