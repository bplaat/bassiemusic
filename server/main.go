package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Connect to the database
	var err error
	db, err = sql.Open("mysql", "bassiemusic:bassiemusic@tcp(127.0.0.1:3306)/bassiemusic?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(32)
	db.SetConnMaxLifetime(time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// Run subcommand
	if len(os.Args) >= 2 {
		if os.Args[1] == "add" {
			startAdd()
			return
		}

		if os.Args[1] == "remove" {
			startRemove()
			return
		}

		if os.Args[1] == "server" {
			startServer()
			return
		}
	}

	// Print error when no or unkown subcommand is given
	log.Fatalln("Expected add, remove or server command")
}
