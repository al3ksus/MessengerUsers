package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	repository "github.com/al3ksus/messengerusers/internal/repositories"
	"github.com/lib/pq"
)

// Repository - объект репозитория
type Repository struct {
	db *sql.DB
}

// Connect создает подключение к базе данных PostgresSQL, принимает на вход строку подключения.
func Connect(conn string) (*sql.DB, error) {
	const op = "psql.New"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("%s. %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s. %w", op, err)
	}

	return db, nil
}

// New возвращает новый объект *Repository
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// SaveUser сохраняет нового пользователя в базу данных, возвращает id нового пользователя.
// В случае нарушения constraint unique, возвращает ошибку repository.ErrUserAlredyExists.
func (r *Repository) SaveUser(ctx context.Context, username string, password []byte) (int64, error) {
	const op = "psql.SaveUser"

	var id int64
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO users (
			username, 
			pass_hash, 
			is_active
		) VALUES ($1, $2, true) RETURNING id`,
		username, password)
	if err := row.Scan(&id); err != nil {
		//Ошибка нарушения constraint unique
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == repository.CodeConstraintUnique {
			return 0, fmt.Errorf("%s, %w", op, repository.ErrUserAlredyExists)
		}

		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

// GetUser получает пользователя по username. Если пользователь не найден, возвращает ошибку repository.ErrUserNotFound.
func (r *Repository) GetUser(ctx context.Context, username string) (models.User, error) {
	const op = "psql.GetUser"

	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1 AND is_active = true", username)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s, %w", op, repository.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s, %w", op, err)
	}

	return user, nil
}

// SetInactive устанавливает пользователю с указанным id значение is_active = false.
// Если пользователь с таким id не найден, возвращает ошибку repository.ErrUserNotFound.
// Если пользователь уже неактивен, возвращает ошибку repository.ErrUserAlreadyInactive.
func (r *Repository) SetInactive(ctx context.Context, userId int64) error {
	const op = "psql.SetInactive"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	row := tx.QueryRowContext(ctx, "SELECT is_active FROM users WHERE id = $1", userId)

	var isActive bool
	err = row.Scan(&isActive)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s, %w", op, repository.ErrUserNotFound)
		}

		return fmt.Errorf("%s, %w", op, err)
	}

	if !isActive {
		return fmt.Errorf("%s, %w", op, repository.ErrUserAlreadyInactive)
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET is_active = FALSE WHERE id = $1", userId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s, %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
