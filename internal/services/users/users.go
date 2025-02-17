package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/logger"
	repository "github.com/al3ksus/messengerusers/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	log          logger.Logger
	userSaver    UserSaver
	userProvider UserProvider
	crypter      Crypter
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UserSaver
type UserSaver interface {
	SaveUser(ctx context.Context, username string, password []byte) (int64, error)
	SetInactive(ctx context.Context, userId int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UserProvider
type UserProvider interface {
	GetUser(ctx context.Context, username string) (models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Crypter
type Crypter interface {
	GenerateFromPassword(pass []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserAlreadyInactive = errors.New("user already inactive")
)

func New(log logger.Logger, userSaver UserSaver, userProvider UserProvider, crypter Crypter) *Users {
	return &Users{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		crypter:      crypter,
	}
}

func (u *Users) Login(ctx context.Context, username, password string) (int64, error) {
	const op = "users.Login"

	user, err := u.userProvider.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			u.log.Warnf("user not found. %w", err)
			return 0, fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}

		u.log.Errorf("error getting user. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	if err = u.crypter.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		u.log.Warnf("invalid credentials. %w", err)
		return 0, fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
	}

	return user.Id, nil
}

func (u *Users) RegisterNewUser(ctx context.Context, username, password string) (int64, error) {
	const op = "users.RegisterNewUser"

	passHash, err := u.crypter.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Errorf("error generating hash from password. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := u.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlredyExists) {
			u.log.Warnf("user already exists. %w", err)
			return 0, fmt.Errorf("%s, %w", op, ErrUserAlreadyExists)
		}

		u.log.Errorf("error saving user. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

func (u *Users) MakeUserInactive(ctx context.Context, userId int64) error {
	const op = "users.MakeUserInactive"

	if err := u.userSaver.SetInactive(ctx, userId); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			u.log.Warnf("user not found. %w", err)
			return fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}
		if errors.Is(err, repository.ErrUserAlreadyInactive) {
			u.log.Warnf("user already inactive. %w", err)
			return fmt.Errorf("%s, %w", op, ErrUserAlreadyInactive)
		}

		u.log.Errorf("error making user inactive. %w", err)
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
