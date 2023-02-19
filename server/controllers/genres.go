package controllers

import (
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.GenreModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Paginate(page, limit))
}

func GenresShow(c *fiber.Ctx) error {
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(genre)
}

func GenresAlbums(c *fiber.Ctx) error {
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}

	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.AlbumModel(c).With("artists", "genres").WhereIn("album_genre", "album_id", "genre_id", genre.ID).
		WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByDesc("released_at").Paginate(page, limit))
}
