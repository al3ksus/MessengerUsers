package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/al3ksus/messengerusers/internal/config"

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
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	"localhost", 5432, "postgres", "7554", "messenger")
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	application := app.New(log.Default(), cfg.GRPCPort)
	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	log.Print("app stopped")
	// fmt.Println("Successfully connected!")
}
