package app

import (
	"database/sql"

	"github.com/al3ksus/messengerusers/internal/app/grpcapp"
	"github.com/al3ksus/messengerusers/internal/lib/crypt"
	"github.com/al3ksus/messengerusers/internal/logger"
	friendrequestspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/friendrequests"
	userspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/users"
	"github.com/al3ksus/messengerusers/internal/services/friendrequests"
	"github.com/al3ksus/messengerusers/internal/services/users"
)

type App struct {
	GRPCServer *grpcapp.GRPCServer
}

func New(log logger.Logger, gRPCPort int, db *sql.DB) *App {
	//Репозиторий (DAO)
	usersrep := userspsql.NewUsersRepository(db)
	frrep := friendrequestspsql.NewFriendRequestsRepository(db)
	crypter := &crypt.Crypter{}

	//Сервис
	users := users.NewUsers(log, usersrep, usersrep, crypter)
	friendRequests := friendrequests.NewFriendRequests(log, frrep, frrep, frrep)

	//обертка grpc сервера
	grpcApp := grpcapp.New(log, gRPCPort, users, friendRequests)

	return &App{
		GRPCServer: grpcApp,
	}
}
