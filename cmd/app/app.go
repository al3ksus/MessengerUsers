package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Id       int
	Password string
	Username string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "7554", "messenger")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select * from users")
	if err != nil {
		log.Fatal(err)
	}

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Password, &user.Username); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	log.Print(users)

	fmt.Println("Successfully connected!")
}
