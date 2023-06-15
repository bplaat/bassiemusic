package commands

import (
	"log"
	"os"

	"github.com/bplaat/bassiemusic/routes"
	"github.com/bplaat/bassiemusic/tasks"
)

func Serve() {
	// Start download task
	go tasks.DownloadTask()

	// Start server with api routes
	api := routes.Api()
	log.Fatalln(api.Listen(":" + os.Getenv("SERVER_PORT")))
}
