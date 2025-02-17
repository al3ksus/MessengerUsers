package usersgrpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	messengerv1 "github.com/al3ksus/messengerprotos/gen/go"
	"github.com/al3ksus/messengerusers/internal/grpc/users/mocks"
	usersservice "github.com/al3ksus/messengerusers/internal/services/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	TestUserId int64 = 1
	// EmptyUserId   int64 = 0
	TestUsername = "user1"
	// EmptyUsername       = ""
	TestPassword = "qwerty"
	// EmptyPassword       = ""
)

var (
	TestErrInvalidCredentials = status.Error(codes.InvalidArgument, "invalid credentials")
	TestErrEmptyUsername      = status.Error(codes.InvalidArgument, "username is required")
	TestErrEmptyPassword      = status.Error(codes.InvalidArgument, "password is required")
	TestErrEmptyUserId        = status.Error(codes.InvalidArgument, "user_id is required")
	TestErrUsernameTaken      = status.Error(codes.AlreadyExists, "username already taken")
	TestErrInternal           = status.Error(codes.Internal, "internal error")
	TestErrAlreadyInactive    = status.Error(codes.AlreadyExists, "user already inactive")
	TestErrUserNotFound       = status.Error(codes.InvalidArgument, "user not found")
)

func Test_serverAPI_Login(t *testing.T) {
	type mockBehavior func(s *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest)
	type args struct {
		ctx context.Context
		in  *messengerv1.LoginRequest
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *messengerv1.LoginResponse
		wantErr      error
	}{
		{
			name: "WrongCredentials",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.LoginRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest) {
				users.On("Login", ctx, in.Username, in.Password).Return(EmptyUserId, usersservice.ErrInvalidCredentials)
			},
			wantErr: TestErrInvalidCredentials,
		},
		{
			name: "InternalError",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.LoginRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest) {
				users.On("Login", ctx, in.Username, in.Password).Return(EmptyUserId, errors.New(""))
			},
			wantErr: TestErrInternal,
		},
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.LoginRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest) {
				users.On("Login", ctx, in.Username, in.Password).Return(TestUserId, nil)
			},
			want: &messengerv1.LoginResponse{
				UserId: TestUserId,
			},
		},
		{
			name: "EmptyPassword",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.LoginRequest{
					Username: TestUsername,
					Password: EmptyPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest) {},
			wantErr:      TestErrEmptyPassword,
		},
		{
			name: "EmptyUsername",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.LoginRequest{
					Username: EmptyUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.LoginRequest) {},
			wantErr:      TestErrEmptyUsername,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := mocks.NewUsers(t)

			tt.mockBehavior(users, tt.args.ctx, tt.args.in)
			s := &serverAPI{
				users: users,
			}
			got, err := s.Login(tt.args.ctx, tt.args.in)
			assert.Equal(t, err, tt.wantErr, fmt.Sprintf("serverAPI.Login() error = %v, wantErr %v", err, tt.wantErr))
			assert.Equal(t, got, tt.want, fmt.Sprintf("serverAPI.Login() = %v, want %v", got, tt.want))
		})
	}
}

func Test_serverAPI_Register(t *testing.T) {
	type mockBehavior func(s *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest)
	type args struct {
		ctx context.Context
		in  *messengerv1.RegisterRequest
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *messengerv1.RegisterResponse
		wantErr      error
	}{
		{
			name: "UsernameTaken",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.RegisterRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest) {
				users.On("RegisterNewUser", ctx, in.Username, in.Password).Return(EmptyUserId, usersservice.ErrUserAlreadyExists)
			},
			wantErr: TestErrUsernameTaken,
		},
		{
			name: "InternalError",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.RegisterRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest) {
				users.On("RegisterNewUser", ctx, in.Username, in.Password).Return(EmptyUserId, errors.New(""))
			},
			wantErr: TestErrInternal,
		},
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.RegisterRequest{
					Username: TestUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest) {
				users.On("RegisterNewUser", ctx, in.Username, in.Password).Return(TestUserId, nil)
			},
			want: &messengerv1.RegisterResponse{
				UserId: TestUserId,
			},
		},
		{
			name: "EmptyPassword",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.RegisterRequest{
					Username: TestUsername,
					Password: EmptyPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest) {},
			wantErr:      TestErrEmptyPassword,
		},
		{
			name: "EmptyUsername",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.RegisterRequest{
					Username: EmptyUsername,
					Password: TestPassword,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.RegisterRequest) {},
			wantErr:      TestErrEmptyUsername,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := mocks.NewUsers(t)

			tt.mockBehavior(users, tt.args.ctx, tt.args.in)
			s := &serverAPI{
				users: users,
			}
			got, err := s.Register(tt.args.ctx, tt.args.in)
			assert.Equal(t, err, tt.wantErr, fmt.Sprintf("serverAPI.Register() error = %v, wantErr %v", err, tt.wantErr))
			assert.Equal(t, got, tt.want, fmt.Sprintf("serverAPI.Register() = %v, want %v", got, tt.want))
		})
	}
}

func Test_serverAPI_ToInactive(t *testing.T) {
	type mockBehavior func(u *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest)
	type args struct {
		ctx context.Context
		in  *messengerv1.ToInactiveRequest
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *messengerv1.Empty
		wantErr      error
	}{
		{
			name: "AlreadyInactive",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.ToInactiveRequest{
					UserId: TestUserId,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest) {
				users.On("MakeUserInactive", ctx, in.UserId).Return(usersservice.ErrUserAlreadyInactive)
			},
			wantErr: TestErrAlreadyInactive,
		},
		{
			name: "WrongId",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.ToInactiveRequest{
					UserId: TestUserId,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest) {
				users.On("MakeUserInactive", ctx, in.UserId).Return(usersservice.ErrInvalidCredentials)
			},
			wantErr: TestErrUserNotFound,
		},
		{
			name: "InternalError",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.ToInactiveRequest{
					UserId: TestUserId,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest) {
				users.On("MakeUserInactive", ctx, in.UserId).Return(errors.New(""))
			},
			wantErr: TestErrInternal,
		},
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.ToInactiveRequest{
					UserId: TestUserId,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest) {
				users.On("MakeUserInactive", ctx, in.UserId).Return(nil)
			},
			want: &messengerv1.Empty{},
		},
		{
			name: "EmptyId",
			args: args{
				ctx: context.Background(),
				in: &messengerv1.ToInactiveRequest{
					UserId: EmptyUserId,
				},
			},
			mockBehavior: func(users *mocks.Users, ctx context.Context, in *messengerv1.ToInactiveRequest) {},
			wantErr:      TestErrEmptyUserId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := mocks.NewUsers(t)

			tt.mockBehavior(users, tt.args.ctx, tt.args.in)
			s := &serverAPI{
				users: users,
			}
			got, err := s.ToInactive(tt.args.ctx, tt.args.in)
			assert.Equal(t, err, tt.wantErr, fmt.Sprintf("serverAPI.Register() error = %v, wantErr %v", err, tt.wantErr))
			assert.Equal(t, got, tt.want, fmt.Sprintf("serverAPI.Register() = %v, want %v", got, tt.want))
		})
	}
}
