package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.GenreModel.WithArgs("liked", c.Locals("authUser")).WhereRaw("`name` LIKE ?", "%"+query+"%")
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

type GenresCreateBody struct {
	Name     *string `form:"name" validate:"min:2"`
	DeezerID *string `form:"deezer_id" validate:"integer"`
	Image    *string `form:"image" validate:"nullable"`
}

func GenresCreate(c *fiber.Ctx) error {
	// Parse body
	var body GenresUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, nil, &body); err != nil {
		return nil
	}

	// Create genre
	genreID := uuid.New()
	models.GenreModel.Create(database.Map{
		"id":        genreID,
		"name":      body.Name,
		"deezer_id": body.DeezerID,
	})

	// Store new genre image
	if imageFile, err := c.FormFile("image"); err == nil {
		if err := utils.StoreUploadedImage(c, "genres", genreID, imageFile, true); err != nil {
			return err
		}
	}

	// Get new genre
	return c.JSON(models.GenreModel.Find(genreID))
}

func GenresShow(c *fiber.Ctx) error {
	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.WithArgs("liked", c.Locals("authUser")).Find(genreID)
	if genre == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(genre)
}

type GenresUpdateBody struct {
	Name     *string `form:"name" validate:"min:2"`
	DeezerID *string `form:"deezer_id" validate:"integer"`
	Image    *string `form:"image" validate:"nullable"`
}

func GenresUpdate(c *fiber.Ctx) error {
	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.Find(genreID)
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Parse body
	var body GenresUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, genre, &body); err != nil {
		return nil
	}

	// Run updates
	updates := database.Map{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.DeezerID != nil {
		deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
		updates["deezer_id"] = deezerID
	}
	if imageFile, err := c.FormFile("image"); err == nil {
		// Remove old image file
		if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); err == nil {
			_ = os.Remove(fmt.Sprintf("storage/genres/original/%s", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
		}

		// Store new image
		if err := utils.StoreUploadedImage(c, "genres", genre.ID, imageFile, true); err != nil {
			return err
		}
	}
	if body.Image != nil && *body.Image == "" {
		// Remove old image file
		if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); err == nil {
			_ = os.Remove(fmt.Sprintf("storage/genres/original/%s", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
			_ = os.Remove(fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
		}
	}
	models.GenreModel.Where("id", genre.ID).Update(updates)

	// Get updated genre
	return c.JSON(models.GenreModel.Find(genre.ID))
}

func GenresDelete(c *fiber.Ctx) error {
	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.Find(genreID)
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
	models.GenreModel.Where("id", genre.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func GenresAlbums(c *fiber.Ctx) error {
	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.Find(genreID)
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Get genre albums
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.AlbumModel.With("artists", "genres").WhereInQuery("id", models.AlbumGenreModel.Select("album_id").Where("genre_id", genre.ID)).
		WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByDesc("released_at").Paginate(page, limit))
}

func GenresLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.Find(genreID)
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Check if genre already liked
	genreLike := models.GenreLikeModel.Where("genre_id", genre.ID).Where("user_id", authUser.ID).First()
	if genreLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like genre
	models.GenreLikeModel.Create(database.Map{
		"genre_id": genre.ID,
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func GenresLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse genre id uuid
	genreID, err := uuid.Parse(c.Params("genreID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if genre exists
	genre := models.GenreModel.Find(genreID)
	if genre == nil {
		return fiber.ErrNotFound
	}

	// Check if genre not liked
	genreLike := models.GenreLikeModel.Where("genre_id", genre.ID).Where("user_id", authUser.ID).First()
	if genreLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.GenreLikeModel.Where("genre_id", genre.ID).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
