package users

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/lib/crypt"
	loggermocks "github.com/al3ksus/messengerusers/internal/logger/mocks"
	userspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/users"
	"github.com/al3ksus/messengerusers/internal/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var (
	TestUsername            = "user1"
	TestPass                = "qwerty"
	TestWrongPassword       = "wrongpass"
	TestPassHash            = generateTestPassHash(TestPass)
	TestUserId        int64 = 1
	EmptyUserId       int64 = 0
)

var (
	TestUser = models.User{
		Id:           TestUserId,
		Username:     TestUsername,
		PasswordHash: TestPassHash,
		IsActive:     true,
	}
	EmptyUser = models.User{}
)

var (
	ErrInvalidCredentialsTest = fmt.Errorf("%s, %w", "users.Login", ErrInvalidCredentials)
)

func TestUsers_Login(t *testing.T) {
	type mockBehavior func(
		log *loggermocks.Logger,
		userSaver *mocks.UserSaver,
		userProvider *mocks.UserProvider,
		crypter *mocks.Crypter,
		ctx context.Context,
		username string,
		password string,
	)

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int64
		wantErr      error
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				userProvider.On("GetUser", ctx, username).Return(TestUser, nil)
				crypter.On("CompareHashAndPassword", TestUser.PasswordHash, []byte(TestPass)).Return(nil)
			},
			want: TestUserId,
		},
		{
			name: "WrongUsername",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				userProvider.On("GetUser", ctx, username).Return(EmptyUser, userspsql.ErrUserNotFound)
				log.On("Warnf", mock.Anything, mock.Anything)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "WrongPassword",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestWrongPassword,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				userProvider.On("GetUser", ctx, username).Return(TestUser, nil)
				crypter.On("CompareHashAndPassword", TestUser.PasswordHash, []byte(TestWrongPassword)).Return(errors.New(""))
				log.On("Warnf", mock.Anything, mock.Anything)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "InternalError",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				userProvider.On("GetUser", ctx, username).Return(EmptyUser, errors.New(""))
				log.On("Errorf", mock.Anything, mock.Anything)
			},
			wantErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSaver := mocks.NewUserSaver(t)
			userProvider := mocks.NewUserProvider(t)
			log := loggermocks.NewLogger(t)
			crypter := mocks.NewCrypter(t)

			tt.mockBehavior(log, userSaver, userProvider, crypter, tt.args.ctx, tt.args.username, tt.args.password)
			u := &Users{
				userSaver:    userSaver,
				log:          log,
				userProvider: userProvider,
				crypter:      crypter,
			}
			got, err := u.Login(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("users.Login() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, "users.Login, "+tt.wantErr.Error(), fmt.Sprintf("users.Login() error = %v, wantErr %v", err, tt.wantErr))
			}
			assert.Equal(t, got, tt.want, fmt.Sprintf("users.Login() = %v, want %v", got, tt.want))
		})
	}
}

func TestUsers_RegisterNewUser(t *testing.T) {
	type mockBehavior func(
		log *loggermocks.Logger,
		userSaver *mocks.UserSaver,
		userProvider *mocks.UserProvider,
		crypter *mocks.Crypter,
		ctx context.Context,
		username string,
		password string,
	)

	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int64
		wantErr      error
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				crypter.On("GenerateFromPassword", []byte(TestPass), bcrypt.DefaultCost).Return(TestPassHash, nil)
				userSaver.On("SaveUser", ctx, username, TestPassHash).Return(TestUserId, nil)
			},
			want: TestUserId,
		},
		{
			name: "UsernameTaken",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				crypter.On("GenerateFromPassword", []byte(TestPass), bcrypt.DefaultCost).Return(TestPassHash, nil)
				userSaver.On("SaveUser", ctx, username, TestPassHash).Return(EmptyUserId, userspsql.ErrUserAlredyExists)
				log.On("Warnf", mock.Anything, mock.Anything)
			},
			wantErr: ErrUserAlreadyExists,
		},
		{
			name: "GenPassError",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestWrongPassword,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				crypter.On("GenerateFromPassword", []byte(TestWrongPassword), bcrypt.DefaultCost).Return(nil, errors.New(""))
				log.On("Errorf", mock.Anything, mock.Anything)
			},
			wantErr: errors.New(""),
		},
		{
			name: "InternalError",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				username string,
				password string,
			) {
				crypter.On("GenerateFromPassword", []byte(TestPass), bcrypt.DefaultCost).Return(TestPassHash, nil)
				userSaver.On("SaveUser", ctx, username, TestPassHash).Return(EmptyUserId, errors.New(""))
				log.On("Errorf", mock.Anything, mock.Anything)
			},
			wantErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSaver := mocks.NewUserSaver(t)
			userProvider := mocks.NewUserProvider(t)
			log := loggermocks.NewLogger(t)
			crypter := mocks.NewCrypter(t)

			tt.mockBehavior(log, userSaver, userProvider, crypter, tt.args.ctx, tt.args.username, tt.args.password)
			u := &Users{
				userSaver:    userSaver,
				log:          log,
				userProvider: userProvider,
				crypter:      crypter,
			}
			got, err := u.RegisterNewUser(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("users.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				assert.EqualError(
					t,
					err,
					"users.RegisterNewUser, "+tt.wantErr.Error(),
					fmt.Sprintf("users.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr),
				)
			}
			assert.Equal(t, got, tt.want, fmt.Sprintf("users.RegisterNewUser() = %v, want %v", got, tt.want))
		})
	}
}

func TestUsers_MakeUserInactive(t *testing.T) {
	type mockBehavior func(
		log *loggermocks.Logger,
		userSaver *mocks.UserSaver,
		userProvider *mocks.UserProvider,
		crypter *mocks.Crypter,
		ctx context.Context,
		userId int64,
	)

	type args struct {
		ctx    context.Context
		userId int64
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				userId int64,
			) {
				userSaver.On("SetInactive", ctx, userId).Return(nil)
			},
		},
		{
			name: "WrongUserId",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				userId int64,
			) {
				userSaver.On("SetInactive", ctx, userId).Return(userspsql.ErrUserNotFound)
				log.On("Warnf", mock.Anything, mock.Anything)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "AlreadyInactive",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				userId int64,
			) {
				userSaver.On("SetInactive", ctx, userId).Return(userspsql.ErrUserAlreadyInactive)
				log.On("Warnf", mock.Anything, mock.Anything)
			},
			wantErr: ErrUserAlreadyInactive,
		},
		{
			name: "InternalError",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(
				log *loggermocks.Logger,
				userSaver *mocks.UserSaver,
				userProvider *mocks.UserProvider,
				crypter *mocks.Crypter,
				ctx context.Context,
				userId int64,
			) {
				userSaver.On("SetInactive", ctx, userId).Return(errors.New(""))
				log.On("Errorf", mock.Anything, mock.Anything)
			},
			wantErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSaver := mocks.NewUserSaver(t)
			userProvider := mocks.NewUserProvider(t)
			log := loggermocks.NewLogger(t)
			crypter := mocks.NewCrypter(t)

			tt.mockBehavior(log, userSaver, userProvider, crypter, tt.args.ctx, tt.args.userId)
			u := &Users{
				userSaver:    userSaver,
				log:          log,
				userProvider: userProvider,
				crypter:      crypter,
			}
			err := u.MakeUserInactive(tt.args.ctx, tt.args.userId)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("users.MakeUserInactive() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				assert.EqualError(
					t,
					err,
					"users.MakeUserInactive, "+tt.wantErr.Error(),
					fmt.Sprintf("users.MakeUserInactive() error = %v, wantErr %v", err, tt.wantErr),
				)
			}
		})
	}
}

func generateTestPassHash(pass string) []byte {
	crypter := crypt.Crypter{}
	passHash, err := crypter.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return passHash
}
