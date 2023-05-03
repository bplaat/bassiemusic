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

var Api *fiber.App

func init() {
	Api = fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})
	Api.Use(compress.New())
	Api.Use(cors.New())
	Api.Use(favicon.New())
	Api.Use(logger.New())

	Api.Get("/", controllers.Home)

	// Host storage folder when in dev env
	if os.Getenv("APP_ENV") == "dev" {
		Api.Use("/storage", func(c *fiber.Ctx) error {
			// Fix bug in Chrome where you can't seek in media files when HTTP range request
			// won't work so fake it by saying we support it but we really dont ;)
			c.Response().Header.Add("Accept-Ranges", "bytes")
			return c.Next()
		})
		Api.Use("/storage", filesystem.New(filesystem.Config{
			Root:   http.Dir("./storage"),
			MaxAge: 30 * 24 * 60 * 60,
		}))
	}

	// Apps
	Api.Get("/apps/macos/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0")
	})
	Api.Get("/apps/macos/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})
	Api.Get("/apps/windows/version", func(c *fiber.Ctx) error {
		return c.SendString("0.1.0.0")
	})
	Api.Get("/apps/windows/download", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/bplaat/bassiemusic/releases")
	})

	// Websocket
	Api.Get("/ws", websocket.ServerHandle)

	// Get agent information
	Api.Get("/agent", func(c *fiber.Ctx) error {
		return c.JSON(utils.ParseUserAgent(c))
	})

	Api.Post("/auth/login", controllers.AuthLogin)

	// Authed routes
	Api.Use(middlewares.IsAuthed)

	Api.Get("/auth/validate", controllers.AuthValidate)
	Api.Put("/auth/logout", controllers.AuthLogout)

	Api.Get("/search", controllers.SearchIndex)
	Api.Get("/deezer_search", controllers.DeezerSearchIndex)

	Api.Get("/artists", controllers.ArtistsIndex)
	Api.Get("/artists/:artistID", controllers.ArtistsShow)
	Api.Get("/artists/:artistID/tracks", controllers.ArtistsTracks)
	Api.Put("/artists/:artistID/like", controllers.ArtistsLike)
	Api.Delete("/artists/:artistID/like", controllers.ArtistsLikeDelete)

	Api.Get("/genres", controllers.GenresIndex)
	Api.Get("/genres/:genreID", controllers.GenresShow)
	Api.Get("/genres/:genreID/albums", controllers.GenresAlbums)
	Api.Put("/genres/:genreID/like", controllers.GenresLike)
	Api.Delete("/genres/:genreID/like", controllers.GenresLikeDelete)

	Api.Get("/albums", controllers.AlbumsIndex)
	Api.Get("/albums/:albumID", controllers.AlbumsShow)
	Api.Put("/albums/:albumID/like", controllers.AlbumsLike)
	Api.Delete("/albums/:albumID/like", controllers.AlbumsLikeDelete)

	Api.Get("/tracks", controllers.TracksIndex)
	Api.Get("/tracks/:trackID", controllers.TracksShow)
	Api.Put("/tracks/:trackID/play", controllers.TracksPlay)
	Api.Put("/tracks/:trackID/like", controllers.TracksLike)
	Api.Delete("/tracks/:trackID/like", controllers.TracksLikeDelete)

	Api.Get("/playlists", controllers.PlaylistsIndex)
	Api.Post("/playlists", controllers.PlaylistsCreate)
	Api.Get("/playlists/:playlistID", controllers.PlaylistsShow)
	Api.Put("/playlists/:playlistID", controllers.PlaylistsUpdate)
	Api.Delete("/playlists/:playlistID", controllers.PlaylistsDelete)
	Api.Post("/playlists/:playlistID/tracks", controllers.PlaylistsAppendTrack)
	Api.Put("/playlists/:playlistID/tracks/:position", controllers.PlaylistsInsertTrack)
	Api.Delete("/playlists/:playlistID/tracks/:position", controllers.PlaylistsRemoveTrack)
	Api.Put("/playlists/:playlistID/like", controllers.PlaylistsLike)
	Api.Delete("/playlists/:playlistID/like", controllers.PlaylistsLikeDelete)

	Api.Get("/users/:userID", controllers.UsersShow)
	Api.Put("/users/:userID", controllers.UsersUpdate)
	Api.Delete("/users/:userID", controllers.UsersDelete)
	Api.Get("/users/:userID/liked_artists", controllers.UsersLikedArtists)
	Api.Get("/users/:userID/liked_genres", controllers.UsersLikedGenres)
	Api.Get("/users/:userID/liked_albums", controllers.UsersLikedAlbums)
	Api.Get("/users/:userID/liked_tracks", controllers.UsersLikedTracks)
	Api.Get("/users/:userID/liked_playlists", controllers.UsersLikedPlaylists)
	Api.Get("/users/:userID/played_tracks", controllers.UsersPlayedTracks)
	Api.Get("/users/:userID/sessions", controllers.UsersSessions)
	Api.Get("/users/:userID/active_sessions", controllers.UsersActiveSessions)
	Api.Get("/users/:userID/playlists", controllers.UsersPlaylists)

	// Admin routes
	Api.Use(middlewares.IsAdmin)

	Api.Get("/storage_size", func(c *fiber.Ctx) error {
		storageSize, err := utils.DirSize("storage/")
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(fiber.Map{
			"used": storageSize,
			"max":  utils.ParseBytes(os.Getenv("STORAGE_MAX_SIZE")),
		})
	})

	Api.Post("/download/artist", controllers.DownloadArtist)
	Api.Post("/download/album", controllers.DownloadAlbum)

	Api.Put("/artists/:artistID", controllers.ArtistsUpdate)
	Api.Delete("/artists/:artistID", controllers.ArtistsDelete)

	Api.Put("/genres/:genreID", controllers.GenresUpdate)
	Api.Delete("/genres/:genreID", controllers.GenresDelete)

	Api.Put("/albums/:albumID", controllers.AlbumsUpdate)
	Api.Delete("/albums/:albumID", controllers.AlbumsDelete)

	Api.Put("/tracks/:trackID", controllers.TracksUpdate)
	Api.Delete("/tracks/:trackID", controllers.TracksDelete)

	Api.Get("/users", controllers.UsersIndex)
	Api.Post("/users", controllers.UsersCreate)

	Api.Get("/sessions", controllers.SessionsIndex)
	Api.Get("/sessions/:sessionID", controllers.SessionsShow)
	Api.Put("/sessions/:sessionID/revoke", controllers.SessionsRevoke)
}
