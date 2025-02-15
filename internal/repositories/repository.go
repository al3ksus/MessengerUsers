package repository

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlredyExists = errors.New("user already exists")
)

var (
	CodeConstraintUnique = "23505"
)
