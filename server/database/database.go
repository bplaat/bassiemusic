package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var db *sql.DB

func Connect() {
	// Connect to the database
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME")))
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(32)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
}

func Query(query string, args ...any) *sql.Rows {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalln(err)
	}
	return rows
}

func Exec(query string, args ...any) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalln(err)
	}
	return result
}
