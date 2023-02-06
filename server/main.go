package main

import (
	"log"
	"net/url"
	"os"

	"github.com/bplaat/bassiemusic/controllers"
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/middlewares"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func createStorageDirs() {
	if _, err := os.Stat("storage/"); os.IsNotExist(err) {
		os.Mkdir("storage", 0755)
	}
	if _, err := os.Stat("storage/artists"); os.IsNotExist(err) {
		os.Mkdir("storage/artists", 0755)
	}
	if _, err := os.Stat("storage/albums"); os.IsNotExist(err) {
		os.Mkdir("storage/albums", 0755)
	}
	if _, err := os.Stat("storage/genres"); os.IsNotExist(err) {
		os.Mkdir("storage/genres", 0755)
	}
	if _, err := os.Stat("storage/tracks"); os.IsNotExist(err) {
		os.Mkdir("storage/tracks", 0755)
	}
}

func main() {
	// Set log settings
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Create missing storage dirs
	createStorageDirs()

	// Connect to the database
	database.Connect()

	// Start download task
	go tasks.DownloadTask()

	// Start server
	app := fiber.New(fiber.Config{
		AppName: "BassieMusic",
	})
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("BassieMusic API")
	})

	app.Static("/storage", "./storage")

	app.Post("/auth/login", controllers.AuthLogin)

	// Deezer API proxies
	app.Get("/deezer/artists", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
		_, err := c.Write(utils.Fetch("https://api.deezer.com/search/artist?q=" + url.QueryEscape(c.Query("q"))))
		return err
	})
	app.Get("/deezer/albums", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
		_, err := c.Write(utils.Fetch("https://api.deezer.com/search/album?q=" + url.QueryEscape(c.Query("q"))))
		return err
	})

	// app.Use(middlewares.IsAuthed)

	app.Get("/auth/validate", controllers.AuthValidate)
	app.Get("/auth/logout", controllers.AuthLogout)

	app.Get("/artists", controllers.ArtistsIndex)
	app.Get("/artists/:artistID", controllers.ArtistsShow)

	app.Get("/search", controllers.SearchIndex)

	app.Get("/albums", controllers.AlbumsIndex)
	app.Get("/albums/:albumID", controllers.AlbumsShow)

	app.Get("/genres", controllers.GenresIndex)
	app.Get("/genres/:genreID", controllers.GenresShow)

	app.Get("/tracks", controllers.TracksIndex)
	app.Get("/tracks/:trackID", controllers.TracksShow)
	app.Get("/tracks/:trackID/play", controllers.TracksPlay)

	app.Get("/users/:userID", controllers.UsersShow)
	app.Get("/users/:userID/sessions", controllers.UsersSessions)
	app.Post("/users/:userID", controllers.UsersEdit)

	app.Use(middlewares.IsAdmin)

	app.Get("/download/artist", controllers.DownloadArtist)
	app.Get("/download/album", controllers.DownloadAlbum)

	app.Get("/users", controllers.UsersIndex)
	app.Post("/users", controllers.UsersCreate)
	app.Get("/users/:userID/delete", controllers.UsersDelete)

	app.Get("/sessions", controllers.SessionsIndex)
	app.Get("/sessions/:sessionID", controllers.SessionsShow)
	app.Get("/sessions/:sessionID/revoke", controllers.SessionsRevoke)

	log.Fatal(app.Listen(":8080"))
}
