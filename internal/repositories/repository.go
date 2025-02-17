package repository

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlredyExists    = errors.New("user already exists")
	ErrUserAlreadyInactive = errors.New("user already inactive")
)

//Код ошибки PostgreSQL
const (
	CodeConstraintUnique = "unique_violation"
)
