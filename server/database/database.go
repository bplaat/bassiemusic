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

func Count(query string, args ...any) int64 {
	countQuery := Query(query, args...)
	defer countQuery.Close()
	countQuery.Next()
	var count int64
	countQuery.Scan(&count)
	return count
}
