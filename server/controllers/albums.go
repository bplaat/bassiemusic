package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func AlbumsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total albums
	total := database.Count("SELECT COUNT(`id`) FROM `albums` WHERE `title` LIKE ?", "%"+query+"%")

	// Get albums
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer albumsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.AlbumsScan(c, albumsQuery, true, true, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func AlbumsShow(c *fiber.Ctx) error {
	// Check if album exists
	albumQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `id` = UUID_TO_BIN(?)", c.Params("albumID"))
	defer albumQuery.Close()
	if !albumQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.AlbumScan(c, albumQuery, true, true, true))
}
