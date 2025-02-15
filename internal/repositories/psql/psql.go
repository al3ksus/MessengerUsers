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

type Repository struct {
	db *sql.DB
}

func Connect(conn string) (*sql.DB, error) {
	const op = "repositories.psql.New"

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

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveUser(ctx context.Context, username string, password []byte) (int64, error) {
	const op = "repositories.psql.SaveUser"

	res, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password, is_active) VALUES ($1, $2, true)", username, password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == repository.CodeConstraintUnique {
			return 0, fmt.Errorf("%s, %w", op, repository.ErrUserAlredyExists)
		}

		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

func (r *Repository) GetUser(ctx context.Context, username string) (models.User, error) {
	const op = "repositories.psql.GetUser"

	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1 AND is_active = true", username)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s, %w", op, repository.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s, %w", op, err)
	}

	return user, nil
}
