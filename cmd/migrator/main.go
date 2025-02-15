package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var database, source, upOrDown string

	flag.StringVar(&database, "database", "", "url to database")
	flag.StringVar(&source, "source", "", "path to migrations")
	flag.Parse()

	upOrDown = flag.Args()[0]

	if database == "" {
		panic("database is required")
	}

	if source == "" {
		panic("source is required")
	}

	m, err := migrate.New(source, database)
	if err != nil {
		panic(err)
	}

	if upOrDown == "up" {
		err = m.Up()
	} else if upOrDown == "down" {
		err = m.Down()
	} else {
		panic("unknown key " + upOrDown)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("migration successfull")
}
