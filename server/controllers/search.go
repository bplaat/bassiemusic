package controllers

import (
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SearchIndex(c *fiber.Ctx) error {
	query, _, _ := utils.ParseIndexVars(c)

	// Get tracks
	tracks := models.TrackModel(c).With("like", "artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("`plays` DESC, `updated_at` DESC").Limit("5").Get()

	// Get albums
	albums := models.AlbumModel(c).With("artists", "genres").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`title`)").Limit("10").Get()

	// Get artists
	artists := models.ArtistModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit("10").Get()

	// Get Genres
	genres := models.GenreModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Limit("10").Get()

	// Return all values
	return c.JSON(fiber.Map{
		"success": true,
		"tracks":  tracks,
		"albums":  albums,
		"artists": artists,
		"genres":  genres,
	})
}
