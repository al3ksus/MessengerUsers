package psql

import (
	"database/sql"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/logger"
)

type PSQLConn struct {
	log logger.Logger
	db  *sql.DB
}

func New(log logger.Logger, host, user, password, dbname string, port int) *PSQLConn {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error connect to database, source=%s. %w", psqlInfo, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error ping database, source=%s. %w", psqlInfo, err)
	}

	log.Infof("successfully connected to database")

	return &PSQLConn{
		log: log,
		db:  db,
	}
}

func (c *PSQLConn) Close() {
	if err := c.db.Close(); err != nil {
		c.log.Errorf("error closing database connection. %w", err)
	} else {
		c.log.Infof("database connection closed")
	}
}
