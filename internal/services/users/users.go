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
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, password []byte) (int64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, username string) (models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

func New(log logger.Logger, userSaver UserSaver, userProvider UserProvider) *Users {
	return &Users{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
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

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		u.log.Errorf("invalid credentials. %w", err)
		return 0, fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
	}

	return user.Id, nil
}

func (u *Users) RegisterNewUser(ctx context.Context, username, password string) (int64, error) {
	const op = "users.RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Errorf("error generating hash from password. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := u.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlredyExists) {
			u.log.Errorf("user already exists. %w", err)
			return 0, fmt.Errorf("%s, %w", op, ErrUserAlreadyExists)
		}

		u.log.Errorf("error saving user. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}
