package psql

import (
	"database/sql"
	"fmt"
)

const (
	CodeConstraintUnique     = "unique_violation"
	CodeConstraintForeignKey = "foreign_key_violation"
)

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
