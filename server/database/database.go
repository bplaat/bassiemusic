package database

import (
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

func Connect() {
	// Connect to the database
	var err error
	db, err = sql.Open("mysql", "bassiemusic:bassiemusic@tcp(127.0.0.1:3306)/bassiemusic?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(16)
	db.SetConnMaxLifetime(time.Minute)

	// Ping the database
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func Exec(query string, args ...any) (sql.Result, error) {
	return db.Exec(query, args...)
}
