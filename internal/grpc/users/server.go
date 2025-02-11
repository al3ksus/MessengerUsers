package usersgrpc

import (
	"context"

	messengerv1 "github.com/al3ksus/messengerprotos/gen/go"
	"google.golang.org/grpc"
)

type serverAPI struct {
	messengerv1.UnimplementedUsersServer
	// users Users
}

// Тот самый интерфейс, котрый мы передавали в grpcApp
type Users interface {
	Login(
		ctx context.Context,
		username string,
		password string,
	) (id int64, err error)
	RegisterNewUser(
		ctx context.Context,
		username string,
		password string,
	) (id int64, err error)
}

func Register(gRPCServer *grpc.Server) {
	messengerv1.RegisterUsersServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *messengerv1.LoginRequest,
) (*messengerv1.LoginResponse, error) {
	// resp, err := s.users.Login(ctx, in.Username, in.Password)
	panic("prrr")
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *messengerv1.RegisterRequest,
) (*messengerv1.RegisterResponse, error) {
	// TODO
	panic("brrr")
}
