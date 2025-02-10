package users

import (
	"context"

	messengerv1 "github.com/al3ksus/messengerprotos/gen/go"
	"google.golang.org/grpc"
)

type serverAPI struct {
	messengerv1.UnimplementedUsersServer
	users Users
}

// Тот самый интерфейс, котрый мы передавали в grpcApp
type Users interface {
	Login(

		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, users Users) {
	messengerv1.RegisterUsersServer(gRPCServer, &serverAPI{users: users})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *messengerv1.LoginRequest,
) (*messengerv1.LoginResponse, error) {
	// TODO
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *messengerv1.RegisterRequest,
) (*messengerv1.RegisterResponse, error) {
	// TODO
}
