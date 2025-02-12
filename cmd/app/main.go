package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/al3ksus/messengerusers/internal/config"
	"go.uber.org/zap"

	"github.com/al3ksus/messengerusers/internal/app"

	_ "github.com/lib/pq"
)

type User struct {
	Id       int
	Password string
	Username string
}

func main() {
	cfg := config.MustLoad()
	logger := setupLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Print(err.Error())
		}
	}()

	application := app.New(logger, cfg.GRPCPort, cfg.DBPort, cfg.Host, cfg.User, cfg.Password, cfg.DBName)
	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	application.PSQLConn.Close()

	logger.Info("app stopped")
}

func setupLogger() *zap.SugaredLogger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic("error setup looger. " + err.Error())
	}

	return l.Sugar()
}
