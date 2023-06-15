package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Database connection
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

// Database queries
var logFile *os.File

type QueryLogLine struct {
	Query string `json:"query"`
	Args  []any  `json:"args"`
	Time  int64  `json:"time"`
}

func logQuery(query string, args []any, start time.Time) {
	if logFile == nil {
		var err error
		logFile, err = os.OpenFile("queries.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalln(err)
		}
	}
	bytes, _ := json.Marshal(QueryLogLine{query, args, time.Since(start).Microseconds()})
	if _, err := logFile.Write(bytes); err != nil {
		log.Fatalln(err)
	}
	if _, err := logFile.WriteString("\n"); err != nil {
		log.Fatalln(err)
	}
}

func Query(query string, args ...any) *sql.Rows {
	start := time.Now()
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalln(err)
	}
	if os.Getenv("DATABASE_LOG") == "YES" {
		logQuery(query, args, start)
	}
	return rows
}

func Exec(query string, args ...any) sql.Result {
	start := time.Now()
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalln(err)
	}
	if os.Getenv("DATABASE_LOG") == "YES" {
		logQuery(query, args, start)
	}
	return result
}
