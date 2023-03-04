package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bplaat/bassiemusic/commands"
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/utils/dotenv"
	_ "github.com/go-sql-driver/mysql"
)

func createDirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Fatalln(err)
		}
	}
}

func createStorageDirs() {
	createDirIfNotExists("storage")
	createDirIfNotExists("storage/avatars")
	createDirIfNotExists("storage/playlists")
	createDirIfNotExists("storage/artists")
	createDirIfNotExists("storage/artists/small")
	createDirIfNotExists("storage/artists/medium")
	createDirIfNotExists("storage/artists/large")
	createDirIfNotExists("storage/albums")
	createDirIfNotExists("storage/albums/small")
	createDirIfNotExists("storage/albums/medium")
	createDirIfNotExists("storage/albums/large")
	createDirIfNotExists("storage/genres")
	createDirIfNotExists("storage/genres/small")
	createDirIfNotExists("storage/genres/medium")
	createDirIfNotExists("storage/genres/large")
	createDirIfNotExists("storage/tracks")
}

func main() {
	// Set log settings
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load .env file
	if err := dotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	// Create missing storage dirs
	createStorageDirs()

	// Connect to the database
	database.Connect()

	// Check arguments for subcommand
	if len(os.Args) > 1 {
		if os.Args[1] == "serve" {
			commands.Serve()
			return
		}

		if os.Args[1] == "restore" {
			commands.Restore()
			return
		}

		if os.Args[1] == "clean" {
			commands.Clean()
			return
		}

		if os.Args[1] == "fix" {
			commands.Fix()
			return
		}
	}

	// When no command is given print help text
	fmt.Print("BassieMusic server executable\n\n" +
		"Usage:\n" +
		"\t./bassiemusic <command>\n\n" +
		"The commands are:\n" +
		"\tserve\t\tStart the BassieMusic server and serve content\n" +
		"\trestore\t\tRedownload the storage/ folder with your filled database\n" +
		"\tclean\t\tClean up storage files that are not referenced in the database\n" +
		"\tfix\t\tCreate missing track entries and try to find and download track YouTube videos\n")
}
