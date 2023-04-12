package controllers

import (
	"fmt"
	"os"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.GenreModel(c).With("liked").WhereRaw("`name` LIKE ?", "%"+query+"%")
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
	genre := models.GenreModel(c).With("liked").Find(c.Params("genreID"))
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

func GenresLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if genre exists
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Check if genre already liked
	genreLike := models.GenreLikeModel().Where("genre_id", c.Params("genreID")).Where("user_id", authUser.ID).First()
	if genreLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like genre
	models.GenreLikeModel().Create(database.Map{
		"genre_id": c.Params("genreID"),
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func GenresLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if genre exists
	genre := models.GenreModel(c).Find(c.Params("genreID"))
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Check if genre not liked
	genreLike := models.GenreLikeModel().Where("genre_id", c.Params("genreID")).Where("user_id", authUser.ID).First()
	if genreLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.GenreLikeModel().Where("genre_id", c.Params("genreID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
