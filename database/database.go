package database

import (
	"database/sql"
	"log"
	"time"
)

var DB *sql.DB

func IntDB() {
	var error error
	connStr := "postgres://postgres:dev@localhost/inventorydb?sslmode=disable"
	DB, error = sql.Open("postgres", connStr)
	if error != nil {
		log.Fatal(error)
	}
	DB.SetConnMaxLifetime(60 * time.Second)
	DB.SetMaxOpenConns(4)
	DB.SetMaxIdleConns(4)
}
