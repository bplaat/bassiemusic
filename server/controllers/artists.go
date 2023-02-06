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
	artistsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `artists` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer artistsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.ArtistsScan(c, artistsQuery, false, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func ArtistsShow(c *fiber.Ctx) error {
	// Check if artist exists
	artistQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `artists` WHERE `id` = UUID_TO_BIN(?)", c.Params("artistID"))
	defer artistQuery.Close()
	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.ArtistScan(c, artistQuery, true, true))
}

func ArtistsLike(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if artist exists
	artistQuery := database.Query("SELECT `id` FROM `artists` WHERE `id` = UUID_TO_BIN(?)", c.Params("artistID"))
	defer artistQuery.Close()
	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if artist_likes binding exists
	artistLikeQuery := database.Query("SELECT `id` FROM `artist_likes` WHERE `artist_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("artistID"), authUser.ID)
	defer artistLikeQuery.Close()
	if artistLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Create artist_likes binding
	database.Exec("INSERT INTO `artist_likes` (`id`, `artist_id`, `user_id`) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), UUID_TO_BIN(?))", c.Params("artistID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLikeDelete(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if artist exists
	artistQuery := database.Query("SELECT `id` FROM `artists` WHERE `id` = UUID_TO_BIN(?)", c.Params("artistID"))
	defer artistQuery.Close()
	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if artist_likes binding doesn't exists
	artistLikeQuery := database.Query("SELECT `id` FROM `artist_likes` WHERE `artist_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("artistID"), authUser.ID)
	defer artistLikeQuery.Close()
	if !artistLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete artist_likes binding
	database.Exec("DELETE FROM `artist_likes` WHERE `artist_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("artistID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}
