package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func ArtistsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total artists
	total := database.Count("SELECT COUNT(`id`) FROM `artists` WHERE `name` LIKE ?", "%"+query+"%")

	// Get artists
	artistsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer artistsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.ArtistsScan(c, artistsQuery, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func ArtistsShow(c *fiber.Ctx) error {
	// Check if artist exists
	artistQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `id` = UUID_TO_BIN(?)", c.Params("artistID"))
	defer artistQuery.Close()
	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.ArtistScan(c, artistQuery, true))
}
