package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SearchIndex(c *fiber.Ctx) error {
	query, _, _ := utils.ParseIndexVars(c)

	// Get artists
	artistsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `artists` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", 0, 10)
	defer artistsQuery.Close()
	artists := models.ArtistsScan(c, artistsQuery, false, false)

	// Get albums
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums` WHERE `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", 0, 10)
	defer albumsQuery.Close()
	albums := models.AlbumsScan(c, albumsQuery, true, true, false)

	// Get tracks
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks` WHERE `title` LIKE ? ORDER BY `plays` DESC, LOWER(`title`) LIMIT ?, ?", "%"+query+"%", 0, 5)
	defer tracksQuery.Close()
	tracks := models.TracksScan(c, tracksQuery, true, true)

	// Get Genres
	genresQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `genres` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", 0, 10)
	genres := models.GenresScan(c, genresQuery, false)

	// Return all values
	return c.JSON(fiber.Map{
		"success": true,
		"artists": artists,
		"albums":  albums,
		"tracks":  tracks,
		"genres":  genres,
	})
}
