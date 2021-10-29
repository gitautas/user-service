package api

import (
	"user-service/src/api/generated"
	"user-service/src/models"
	"user-service/src/storage"

	"google.golang.org/grpc"
)

type RpcServer struct {
	Server *grpc.Server
	mysql storage.Database
}

func NewRpcServer(mysql storage.Database) *RpcServer {
	server := grpc.NewServer()

	return &RpcServer{
		Server: server,
		mysql:  mysql,
	}
}

func (rs *RpcServer) CreateUser(req *generated.CreateUserReq) (*generated.CreateUserResp, error) {
	user, err := CreateUser(models.FromRPC(req.User), rs.mysql)

	return &generated.CreateUserResp{
		User: user.ToRPC(),
	}, err
}

func (rs *RpcServer) UpdateUser(req *generated.UpdateUserReq) (*generated.UpdateUserResp, error) {
	user, err := UpdateUser(models.FromRPC(req.User), rs.mysql)

	return &generated.UpdateUserResp{
		User: user.ToRPC(),
	}, err
}

func (rs *RpcServer) DeleteUser(req *generated.DeleteUserReq) (*generated.DeleteUserResp, error) {
	err := RemoveUser(req.UserId, rs.mysql)
	if err != nil {
		return nil, err
	}

	return &generated.DeleteUserResp{}, nil
}

func (rs *RpcServer) GetUser(req *generated.GetUserReq) (*generated.GetUserResp, error) {
	user, err := GetUser(req.UserId, rs.mysql)
	if err != nil {
		return nil, err
	}

	return &generated.GetUserResp{ // I could simplify these models, but this makes for more readable proto files.
		User: user.ToRPC(),
	}, nil
}

func (rs *RpcServer) GetUserList(req *generated.GetUserListReq, stream generated.UserService_GetUserListServer) (error) {
	users, err := GetUsers(int(req.PageSize), int(req.PageOffset), rs.mysql)
	if err != nil {
		return err
	}

	for _, user := range(users) {
		err = stream.Send(&generated.GetUserListResp{
			User: user.ToRPC(),
		})
	}

	return nil

}
