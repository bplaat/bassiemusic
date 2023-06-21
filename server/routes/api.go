package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/bplaat/bassiemusic/controllers"
	"github.com/bplaat/bassiemusic/controllers/websocket"
	"github.com/bplaat/bassiemusic/middlewares"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Api() *fiber.App {
	api := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})
	api.Use(compress.New())
	api.Use(cors.New(cors.Config{
		MaxAge: 60 * 60,
	}))
	api.Use(favicon.New())
	api.Use(logger.New())

	api.Get("/", controllers.Home)

	// Host storage folder when in dev env
	if os.Getenv("APP_ENV") == "dev" {
		api.Use("/storage", func(c *fiber.Ctx) error {
			// Fix bug in Chrome where you can't seek in media files when HTTP range request
			// won't work so fake it by saying we support it but we really dont ;)
			c.Response().Header.Add("Accept-Ranges", "bytes")
			return c.Next()
		})
		api.Use("/storage", filesystem.New(filesystem.Config{
			Root:   http.Dir("./storage"),
			MaxAge: 30 * 24 * 60 * 60,
		}))
	}

	// Apps
	api.Get("/apps/macos/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0")
	})
	api.Get("/apps/macos/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})
	api.Get("/apps/windows/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0.0")
	})
	api.Get("/apps/windows/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})

	// Websocket
	api.Get("/ws", websocket.ServerHandle)

	// Get agent information
	api.Get("/agent", func(c *fiber.Ctx) error {
		return c.JSON(utils.ParseUserAgent(c))
	})

	api.Post("/auth/login", controllers.AuthLogin)

	// Authed routes
	api.Use(middlewares.IsAuthed)

	api.Get("/auth/validate", controllers.AuthValidate)
	api.Put("/auth/logout", controllers.AuthLogout)

	api.Get("/search", controllers.SearchIndex)
	api.Get("/deezer_search", controllers.DeezerSearchIndex)

	api.Get("/artists", controllers.ArtistsIndex)
	api.Get("/artists/:artistID", controllers.ArtistsShow)
	api.Get("/artists/:artistID/tracks", controllers.ArtistsTracks)
	api.Put("/artists/:artistID/like", controllers.ArtistsLike)
	api.Delete("/artists/:artistID/like", controllers.ArtistsLikeDelete)

	api.Get("/genres", controllers.GenresIndex)
	api.Get("/genres/:genreID", controllers.GenresShow)
	api.Get("/genres/:genreID/albums", controllers.GenresAlbums)
	api.Put("/genres/:genreID/like", controllers.GenresLike)
	api.Delete("/genres/:genreID/like", controllers.GenresLikeDelete)

	api.Get("/albums", controllers.AlbumsIndex)
	api.Get("/albums/:albumID", controllers.AlbumsShow)
	api.Put("/albums/:albumID/like", controllers.AlbumsLike)
	api.Delete("/albums/:albumID/like", controllers.AlbumsLikeDelete)

	api.Get("/tracks", controllers.TracksIndex)
	api.Get("/tracks/:trackID", controllers.TracksShow)
	api.Put("/tracks/:trackID/play", controllers.TracksPlay)
	api.Put("/tracks/:trackID/like", controllers.TracksLike)
	api.Delete("/tracks/:trackID/like", controllers.TracksLikeDelete)

	api.Get("/playlists", controllers.PlaylistsIndex)
	api.Post("/playlists", controllers.PlaylistsCreate)
	api.Get("/playlists/:playlistID", controllers.PlaylistsShow)
	api.Put("/playlists/:playlistID", controllers.PlaylistsUpdate)
	api.Delete("/playlists/:playlistID", controllers.PlaylistsDelete)
	api.Post("/playlists/:playlistID/tracks", controllers.PlaylistsAppendTrack)
	api.Put("/playlists/:playlistID/tracks/:position", controllers.PlaylistsInsertTrack)
	api.Delete("/playlists/:playlistID/tracks/:position", controllers.PlaylistsRemoveTrack)
	api.Put("/playlists/:playlistID/like", controllers.PlaylistsLike)
	api.Delete("/playlists/:playlistID/like", controllers.PlaylistsLikeDelete)

	api.Get("/users/:userID", controllers.UsersShow)
	api.Put("/users/:userID", controllers.UsersUpdate)
	api.Delete("/users/:userID", controllers.UsersDelete)
	api.Get("/users/:userID/liked_artists", controllers.UsersLikedArtists)
	api.Get("/users/:userID/liked_genres", controllers.UsersLikedGenres)
	api.Get("/users/:userID/liked_albums", controllers.UsersLikedAlbums)
	api.Get("/users/:userID/liked_tracks", controllers.UsersLikedTracks)
	api.Get("/users/:userID/liked_playlists", controllers.UsersLikedPlaylists)
	api.Get("/users/:userID/played_tracks", controllers.UsersPlayedTracks)
	api.Get("/users/:userID/sessions", controllers.UsersSessions)
	api.Get("/users/:userID/active_sessions", controllers.UsersActiveSessions)
	api.Get("/users/:userID/playlists", controllers.UsersPlaylists)

	// Admin routes
	api.Use(middlewares.IsAdmin)

	api.Get("/storage_size", func(c *fiber.Ctx) error {
		storageSize, err := utils.DirSize("storage/")
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(fiber.Map{
			"used": storageSize,
			"max":  utils.ParseBytes(os.Getenv("STORAGE_MAX_SIZE")),
		})
	})

	api.Post("/download/artist", controllers.DownloadArtist)
	api.Post("/download/album", controllers.DownloadAlbum)
	api.Delete("/download/:downloadTaskID", controllers.DownloadDelete)

	api.Post("/artists", controllers.ArtistsCreate)
	api.Put("/artists/:artistID", controllers.ArtistsUpdate)
	api.Delete("/artists/:artistID", controllers.ArtistsDelete)

	api.Post("/genres", controllers.GenresCreate)
	api.Put("/genres/:genreID", controllers.GenresUpdate)
	api.Delete("/genres/:genreID", controllers.GenresDelete)

	api.Post("/albums", controllers.AlbumsCreate)
	api.Put("/albums/:albumID", controllers.AlbumsUpdate)
	api.Delete("/albums/:albumID", controllers.AlbumsDelete)

	api.Put("/tracks/:trackID", controllers.TracksUpdate)
	api.Delete("/tracks/:trackID", controllers.TracksDelete)

	api.Get("/users", controllers.UsersIndex)
	api.Post("/users", controllers.UsersCreate)

	api.Get("/sessions", controllers.SessionsIndex)
	api.Get("/sessions/:sessionID", controllers.SessionsShow)
	api.Put("/sessions/:sessionID/revoke", controllers.SessionsRevoke)

	return api
}
