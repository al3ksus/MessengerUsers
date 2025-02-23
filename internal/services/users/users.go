package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/logger"
	userspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/users"
	"golang.org/x/crypto/bcrypt"
)

// Users - объект сервиса, реализует логику работы с данными пользователя.
type Users struct {
	log          logger.Logger
	userSaver    UserSaver
	userProvider UserProvider
	crypter      Crypter
}

// UserSaver предоставляет методы создания новых пользователей и изменения существующих.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UserSaver
type UserSaver interface {
	// SaveUser сохраняет нового пользователя в базу данных, возвращает id нового пользователя.
	// В случае нарушения constraint unique, возвращает ошибку repository.ErrUserAlredyExists.
	SaveUser(ctx context.Context, username string, password []byte) (int64, error)

	// SetInactive устанавливает пользователю с указанным id значение is_active = false.
	// Если пользователь с таким id не найден, возвращает ошибку repository.ErrUserNotFound.
	// Если пользователь уже неактивен, возвращает ошибку repository.ErrUserAlreadyInactive.
	SetInactive(ctx context.Context, userId int64) error
}

// UserProvider предоставляет методы получения пользователей.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UserProvider
type UserProvider interface {
	// GetUser получает пользователя по username. Если пользователь не найден, возвращает ошибку repository.ErrUserNotFound.
	GetUser(ctx context.Context, username string) (models.User, error)
}

// Crypter - интерфейс для работы с хэшами.
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Crypter
type Crypter interface {
	// GenerateFromPassword возвращает хэш указанного пароля с заданной стоимостью.
	GenerateFromPassword(pass []byte, cost int) ([]byte, error)

	// CompareHashAndPassword сравнивает захэшированный пароль с исходным.
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserAlreadyInactive = errors.New("user already inactive")
)

// New - конструктор для типа Users.
func NewUsers(log logger.Logger, userSaver UserSaver, userProvider UserProvider, crypter Crypter) *Users {
	return &Users{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		crypter:      crypter,
	}
}

// Login реализует логику авторизации пользователя по логину и паролю.
// Если логин или пароль неверные, возвращает users.ErrInvalidCredentials.
func (u *Users) Login(ctx context.Context, username, password string) (int64, error) {
	const op = "users.Login"

	user, err := u.userProvider.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, userspsql.ErrUserNotFound) {
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

// RegisterNewUser реализует логику регистрации нового пользователя.
// Если заданный username уже занят, возвращает users.ErrUserAlreadyExists.
func (u *Users) RegisterNewUser(ctx context.Context, username, password string) (int64, error) {
	const op = "users.RegisterNewUser"

	passHash, err := u.crypter.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Errorf("error generating hash from password. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := u.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, userspsql.ErrUserAlredyExists) {
			u.log.Warnf("user already exists. %w", err)
			return 0, fmt.Errorf("%s, %w", op, ErrUserAlreadyExists)
		}

		u.log.Errorf("error saving user. %w", err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

// MakeUserInactive реализует логику переведения пользователя в статус 'неактивен'.
// Если пользователь с заданным userId не найден, возвращает users.ErrInvalidCredentials.
// Если найденный пользователь уже имеет статус 'неактивен', возвращает ошибку repository.ErrUserAlreadyInactive.
func (u *Users) MakeUserInactive(ctx context.Context, userId int64) error {
	const op = "users.MakeUserInactive"

	if err := u.userSaver.SetInactive(ctx, userId); err != nil {
		if errors.Is(err, userspsql.ErrUserNotFound) {
			u.log.Warnf("user not found. %w", err)
			return fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}
		if errors.Is(err, userspsql.ErrUserAlreadyInactive) {
			u.log.Warnf("user already inactive. %w", err)
			return fmt.Errorf("%s, %w", op, ErrUserAlreadyInactive)
		}

		u.log.Errorf("error making user inactive. %w", err)
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
