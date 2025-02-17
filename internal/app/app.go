package app

import (
	"database/sql"

	"github.com/al3ksus/messengerusers/internal/app/grpcapp"
	"github.com/al3ksus/messengerusers/internal/lib/crypt"
	"github.com/al3ksus/messengerusers/internal/logger"
	"github.com/al3ksus/messengerusers/internal/repositories/psql"
	"github.com/al3ksus/messengerusers/internal/services/users"
)

type App struct {
	GRPCServer *grpcapp.GRPCServer
}

func New(log logger.Logger, gRPCPort int, db *sql.DB) *App {
	rep := psql.New(db)
	crypter := &crypt.Crypter{}

	users := users.New(log, rep, rep, crypter)
	grpcApp := grpcapp.New(log, gRPCPort, users)

	return &App{
		GRPCServer: grpcApp,
	}
}
