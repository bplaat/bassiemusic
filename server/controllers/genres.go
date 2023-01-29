package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total genres
	total := database.Count("SELECT COUNT(`id`) FROM `genres` WHERE `name` LIKE ?", "%"+query+"%")

	// Get genres
	genresQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `genres` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer genresQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.GenresScan(c, genresQuery, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GenresShow(c *fiber.Ctx) error {
	// Check if genre exists
	genreQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `genres` WHERE `id` = UUID_TO_BIN(?)", c.Params("genreID"))
	defer genreQuery.Close()
	if !genreQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.GenreScan(c, genreQuery, true))
}
