package app

import (
	"log"
	"usr/internal/app/grpcapp"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(log *log.Logger, port int) *App {
	grpcApp := grpcapp.New(log, port)

	return &App{
		GrpcServer: grpcApp,
	}
}
