package commands

import (
	"log"
	"os"

	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/routes"
)

func Serve() {
	// Start download task
	go tasks.DownloadTask()

	// Start server with api routes
	log.Fatal(routes.Api.Listen(":" + os.Getenv("SERVER_PORT")))
}
