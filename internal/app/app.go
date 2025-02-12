package app

import (
	"github.com/al3ksus/messengerusers/internal/app/grpcapp"
	"github.com/al3ksus/messengerusers/internal/app/psql"
	"github.com/al3ksus/messengerusers/internal/logger"
)

type App struct {
	GRPCServer *grpcapp.GRPCServer
	PSQLConn   *psql.PSQLConn
}

func New(log logger.Logger, gRPCPort, dbPort int, dbHost, dbUser, dbPassword, dbName string) *App {
	grpcApp := grpcapp.New(log, gRPCPort)
	psqlConn := psql.New(log, dbHost, dbUser, dbPassword, dbName, dbPort)

	return &App{
		GRPCServer: grpcApp,
		PSQLConn:   psqlConn,
	}
}
