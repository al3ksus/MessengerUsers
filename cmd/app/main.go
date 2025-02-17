package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/al3ksus/messengerusers/internal/config"
	"github.com/al3ksus/messengerusers/internal/repositories/psql"
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

	db, err := psql.Connect(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.DBPort, cfg.User, cfg.Password, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	application := app.New(logger, cfg.GRPCPort, db)
	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	logger.Info("app stopped")
}

// setupLogger получает объект zap.logger и возвращает zap.SugaredLogger.
// Вызывыет панику в случае ошибки
func setupLogger() *zap.SugaredLogger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic("error setup looger. " + err.Error())
	}

	return l.Sugar()
}
