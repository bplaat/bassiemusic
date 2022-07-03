package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
)

var db *sql.DB

func main() {
	// Parse flags
	downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)

	// Connect to the database
	var err error
	db, err = sql.Open("mysql", "bassiemusic:bassiemusic@tcp(127.0.0.1:3306)/bassiemusic?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Run subcommand
	if len(os.Args) >= 2 {
		if os.Args[1] == "download" {
			downloadCommand.Parse(os.Args[2:])
			downloadTracks()
			return
		}

		if os.Args[1] == "server" {
			serverCommand.Parse(os.Args[2:])
			startServer()
			return
		}
	}

	// Print error when no or unkown subcommand is given
	log.Fatalln("Expected download or serve command")
}
