package usersgrpc

import (
	"context"
	"errors"

	messengerv1 "github.com/al3ksus/messengerprotos/gen/go"
	"github.com/al3ksus/messengerusers/internal/services/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	messengerv1.UnimplementedUsersServer
	users Users
}

type Users interface {
	Login(ctx context.Context, username string, password string) (id int64, err error)
	RegisterNewUser(ctx context.Context, username string, password string) (id int64, err error)
	MakeUserInactive(ctx context.Context, userId int64) error
}

func Register(gRPCServer *grpc.Server, users Users) {
	messengerv1.RegisterUsersServer(gRPCServer, &serverAPI{users: users})
}

func (s *serverAPI) Login(ctx context.Context, in *messengerv1.LoginRequest) (*messengerv1.LoginResponse, error) {
	if err := validate(in.Password, in.Username); err != nil {
		return nil, err
	}

	id, err := s.users.Login(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.LoginResponse{
		UserId: id,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *messengerv1.RegisterRequest) (*messengerv1.RegisterResponse, error) {
	if err := validate(in.Password, in.Username); err != nil {
		return nil, err
	}

	id, err := s.users.RegisterNewUser(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "username already taken")

		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.RegisterResponse{
		UserId: id,
	}, nil
}

func (s *serverAPI) ToInactive(ctx context.Context, in *messengerv1.ToInactiveRequest) (*messengerv1.Empty, error) {
	if err := s.users.MakeUserInactive(ctx, in.GetUserId()); err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "user not found")
		}
		if errors.Is(err, users.ErrUserAlreadyInactive) {
			return nil, status.Error(codes.AlreadyExists, "user already inactive")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.Empty{}, nil
}

func validate(password, username string) error {
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if username == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	return nil
}
