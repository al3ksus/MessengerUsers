package service

type UsersService interface {
	Login(username string, password string) error
	Register(username string, password string) error
}
