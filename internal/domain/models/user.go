package models

// Users модель данных пользователя.
type User struct {
	Id           int64
	Username     string
	PasswordHash []byte
	IsActive     bool
}
