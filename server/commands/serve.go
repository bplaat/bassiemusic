package commands

import (
	"fmt"
	"log"
	"net/http"
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

func Serve() {
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
		return c.SendString(fmt.Sprintf("%s API v%s", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION")))
	})

	// Host storage folder when in dev env
	if os.Getenv("APP_ENV") == "dev" {
		app.Use("/storage", func(c *fiber.Ctx) error {
			// Fix bug in Chrome where you can't seek in media files when HTTP range request
			// won't work so fake it by saying we support it but we really dont ;)
			c.Response().Header.Add("Accept-Ranges", "bytes")
			return c.Next()
		})
		app.Use("/storage", filesystem.New(filesystem.Config{
			Root:   http.Dir("./storage"),
			MaxAge: 30 * 24 * 60 * 60,
		}))
	}

	// Apps
	app.Get("/apps/macos/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0")
	})
	app.Get("/apps/macos/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})
	app.Get("/apps/windows/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0.0")
	})
	app.Get("/apps/windows/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})

	// Websocket
	app.Get("/ws", controllers.Websocket)

	// Get agent information
	app.Get("/agent", func(c *fiber.Ctx) error {
		return c.JSON(utils.ParseUserAgent(c))
	})

	app.Post("/auth/login", controllers.AuthLogin)

	app.Use(middlewares.IsAuthed)

	app.Get("/auth/validate", controllers.AuthValidate)
	app.Put("/auth/logout", controllers.AuthLogout)

	app.Get("/search", controllers.SearchIndex)
	app.Get("/deezer_search", controllers.DeezerSearchIndex)

	app.Get("/artists", controllers.ArtistsIndex)
	app.Get("/artists/:artistID", controllers.ArtistsShow)
	app.Put("/artists/:artistID/like", controllers.ArtistsLike)
	app.Delete("/artists/:artistID/like", controllers.ArtistsLikeDelete)

	app.Get("/albums", controllers.AlbumsIndex)
	app.Get("/albums/:albumID", controllers.AlbumsShow)
	app.Put("/albums/:albumID/like", controllers.AlbumsLike)
	app.Delete("/albums/:albumID/like", controllers.AlbumsLikeDelete)

	app.Get("/genres", controllers.GenresIndex)
	app.Get("/genres/:genreID", controllers.GenresShow)
	app.Get("/genres/:genreID/albums", controllers.GenresAlbums)

	app.Get("/tracks", controllers.TracksIndex)
	app.Get("/tracks/:trackID", controllers.TracksShow)
	app.Put("/tracks/:trackID/like", controllers.TracksLike)
	app.Delete("/tracks/:trackID/like", controllers.TracksLikeDelete)
	app.Put("/tracks/:trackID/play", controllers.TracksPlay)

	app.Get("/playlists", controllers.PlaylistsIndex)
	app.Post("/playlists", controllers.PlaylistsCreate)
	app.Get("/playlists/:playlistID", controllers.PlaylistsShow)
	app.Put("/playlists/:playlistID", controllers.PlaylistsEdit)
	app.Delete("/playlists/:playlistID", controllers.PlaylistsDelete)
	app.Post("/playlists/:playlistID/image", controllers.PlaylistsImage)
	app.Delete("/playlists/:playlistID/image", controllers.PlaylistsImageDelete)
	app.Post("/playlists/:playlistID/tracks", controllers.PlaylistsAppendTrack)
	app.Put("/playlists/:playlistID/tracks/:position", controllers.PlaylistsInsertTrack)
	app.Delete("/playlists/:playlistID/tracks/:position", controllers.PlaylistsRemoveTrack)
	app.Put("/playlists/:playlistID/like", controllers.PlaylistsLike)
	app.Delete("/playlists/:playlistID/like", controllers.PlaylistsLikeDelete)

	app.Get("/users/:userID", controllers.UsersShow)
	app.Put("/users/:userID", controllers.UsersEdit)
	app.Delete("/users/:userID", controllers.UsersDelete)
	app.Post("/users/:userID/avatar", controllers.UsersAvatar)
	app.Delete("/users/:userID/avatar", controllers.UsersAvatarDelete)
	app.Get("/users/:userID/liked_artists", controllers.UsersLikedArtists)
	app.Get("/users/:userID/liked_albums", controllers.UsersLikedAlbums)
	app.Get("/users/:userID/liked_tracks", controllers.UsersLikedTracks)
	app.Get("/users/:userID/liked_playlists", controllers.UsersLikedPlaylists)
	app.Get("/users/:userID/played_tracks", controllers.UsersPlayedTracks)
	app.Get("/users/:userID/sessions", controllers.UsersSessions)
	app.Get("/users/:userID/active_sessions", controllers.UsersActiveSessions)
	app.Get("/users/:userID/playlists", controllers.UsersPlaylists)

	app.Use(middlewares.IsAdmin)

	app.Get("/storage_size", func(c *fiber.Ctx) error {
		storageSize, err := utils.DirSize("storage/")
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(fiber.Map{
			"used": storageSize,
			"max":  utils.ParseBytes(os.Getenv("STORAGE_MAX_SIZE")),
		})
	})

	app.Post("/download/artist", controllers.DownloadArtist)
	app.Post("/download/album", controllers.DownloadAlbum)

	app.Get("/users", controllers.UsersIndex)
	app.Post("/users", controllers.UsersCreate)

	app.Get("/sessions", controllers.SessionsIndex)
	app.Get("/sessions/:sessionID", controllers.SessionsShow)
	app.Put("/sessions/:sessionID/revoke", controllers.SessionsRevoke)

	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
