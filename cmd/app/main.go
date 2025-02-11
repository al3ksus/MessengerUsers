package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"usr/internal/app"

	_ "github.com/lib/pq"
)

type User struct {
	Id       int
	Password string
	Username string
}

func main() {
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

	application := app.New(log.Default(), 8080)
	go application.GrpcServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GrpcServer.Stop()

	log.Print("app stopped")
	// fmt.Println("Successfully connected!")
}
