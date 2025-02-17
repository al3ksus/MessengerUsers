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

// serverAPI реализует хэндлеры
type serverAPI struct {
	messengerv1.UnimplementedUsersServer
	users Users
}

var (
	EmptyPassword       = ""
	EmptyUsername       = ""
	EmptyUserId   int64 = 0
)

// Users предоставляет методы для работы с сервисным слоем приложения.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Users
type Users interface {
	// Login - авторизация пользователя по логину и паролю.
	// Если логин или пароль неверные, возвращает users.ErrInvalidCredentials.
	Login(ctx context.Context, username string, password string) (id int64, err error)

	// RegisterNewUser - регистрация нового пользователя.
	// Если заданный username уже занят, возвращает users.ErrUserAlreadyExists.
	RegisterNewUser(ctx context.Context, username string, password string) (id int64, err error)

	// MakeUserInactive переводит пользователя в статус 'неактивен'.
	// Если пользователь с заданным userId не найден, возвращает users.ErrInvalidCredentials.
	// Если найденный пользователь уже имеет статус 'неактивен', возвращает ошибку repository.ErrUserAlreadyInactive.
	MakeUserInactive(ctx context.Context, userId int64) error
}

// Register регистрирует grpc сервер
func Register(gRPCServer *grpc.Server, users Users) {
	messengerv1.RegisterUsersServer(gRPCServer, &serverAPI{users: users})
}

// Хэндлер Login отвечает за авторизацию пользователей по логину и паролю.
// Если логин или пароль неверные, возвращает ошибку InvalidArguments.
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

// Хэндлер Register отвечает за регистрацию новых пользователей.
// Если логин уже занят, возвращает ошибку AlreadyExists.
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

// Хэндлер ToInactive отвечает за перевод пользователей в состояние 'неаткивен'.
// Если пользователь не найден, возвращает ошибку InvalidArgument.
// Если пользователь уже неактивен, возвращает ошибку AlreadyExists.
func (s *serverAPI) ToInactive(ctx context.Context, in *messengerv1.ToInactiveRequest) (*messengerv1.Empty, error) {
	if err := validateId(in.UserId); err != nil {
		return nil, err
	}

	if err := s.users.MakeUserInactive(ctx, in.GetUserId()); err != nil {
		//Пользователь с таким id не существует или неактивен
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

// validate валидирует пароль и логин.
// Проверка на пустоту.
func validate(password, username string) error {
	if username == EmptyUsername {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	if password == EmptyPassword {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// validateId проверяет id пользователя на пустоту.
func validateId(userId int64) error {
	if userId == EmptyUserId {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
