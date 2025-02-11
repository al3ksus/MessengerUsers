package app

import (
	"log"

	"github.com/al3ksus/messengerusers/internal/app/grpcapp"
)

type App struct {
	GRPCServer *grpcapp.GRPCServer
}

func New(log *log.Logger, port int) *App {
	grpcApp := grpcapp.New(log, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
