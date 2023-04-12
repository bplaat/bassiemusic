package controllers

import (
	"fmt"
	"os"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.GenreModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "created_at" {
		q = q.OrderBy("created_at")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByDesc("created_at")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`name`) DESC")
	} else {
		q = q.OrderByRaw("LOWER(`name`)")
	}
	return c.JSON(q.Paginate(page, limit))
}

func GenresShow(c *fiber.Ctx) error {
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(genre)
}

func GenresDelete(c *fiber.Ctx) error {
	// Check if genre exists
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Delete genre image if exists
	if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); err == nil {
		_ = os.Remove(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
		_ = os.Remove(fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
		_ = os.Remove(fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
	}

	// Delete genre
	models.GenreModel(c).Where("id", genre.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
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
