package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/bplaat/bassiemusic/controllers"
	"github.com/bplaat/bassiemusic/middlewares"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func serve() {
	// Start download task
	go tasks.DownloadTask()

	// Start server
	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(os.Getenv("APP_NAME") + " API v" + os.Getenv("APP_VERSION"))
	})

	app.Use("/storage", func(c *fiber.Ctx) error {
		// Fix bug in Chrome where you can't seek in media files when HTTP range request won't work
		// so fake it by saying we support it but we really dont ;)
		c.Response().Header.Add("Accept-Ranges", "bytes")
		return c.Next()
	})
	app.Use("/storage", filesystem.New(filesystem.Config{
		Root:   http.Dir("./storage"),
		MaxAge: 24 * 60 * 60,
	}))

	// Get agent information
	app.Get("/agent", func(c *fiber.Ctx) error {
		return c.JSON(utils.ParseUserAgent(c))
	})

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

	app.Use(middlewares.IsAuthed)

	app.Get("/auth/validate", controllers.AuthValidate)
	app.Get("/auth/logout", controllers.AuthLogout)

	app.Get("/search", controllers.SearchIndex)

	app.Get("/artists", controllers.ArtistsIndex)
	app.Get("/artists/:artistID", controllers.ArtistsShow)
	app.Get("/artists/:artistID/like", controllers.ArtistsLike)
	app.Get("/artists/:artistID/like/delete", controllers.ArtistsLikeDelete)

	app.Get("/albums", controllers.AlbumsIndex)
	app.Get("/albums/:albumID", controllers.AlbumsShow)
	app.Get("/albums/:albumID/like", controllers.AlbumsLike)
	app.Get("/albums/:albumID/like/delete", controllers.AlbumsLikeDelete)

	app.Get("/genres", controllers.GenresIndex)
	app.Get("/genres/:genreID", controllers.GenresShow)

	app.Get("/tracks", controllers.TracksIndex)
	app.Get("/tracks/:trackID", controllers.TracksShow)
	app.Get("/tracks/:trackID/like", controllers.TracksLike)
	app.Get("/tracks/:trackID/like/delete", controllers.TracksLikeDelete)
	app.Get("/tracks/:trackID/play", controllers.TracksPlay)

	app.Get("/users/:userID", controllers.UsersShow)
	app.Get("/users/:userID/liked_artists", controllers.UsersLikedArtists)
	app.Get("/users/:userID/liked_albums", controllers.UsersLikedAlbums)
	app.Get("/users/:userID/liked_tracks", controllers.UsersLikedTracks)
	app.Get("/users/:userID/played_tracks", controllers.UsersPlayedTracks)
	app.Get("/users/:userID/sessions", controllers.UsersSessions)
	app.Post("/users/:userID", controllers.UsersEdit)
	app.Post("/users/:userID/avatar", controllers.UsersAvatar)
	app.Get("/users/:userID/avatar/delete", controllers.UsersAvatarDelete)

	app.Use(middlewares.IsAdmin)

	app.Get("/download/artist", controllers.DownloadArtist)
	app.Get("/download/album", controllers.DownloadAlbum)

	app.Get("/users", controllers.UsersIndex)
	app.Post("/users", controllers.UsersCreate)
	app.Get("/users/:userID/delete", controllers.UsersDelete)

	app.Get("/sessions", controllers.SessionsIndex)
	app.Get("/sessions/:sessionID", controllers.SessionsShow)
	app.Get("/sessions/:sessionID/revoke", controllers.SessionsRevoke)

	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
