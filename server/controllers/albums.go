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
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums` WHERE `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
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
	albumQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums` WHERE `id` = UUID_TO_BIN(?)", c.Params("albumID"))
	defer albumQuery.Close()
	if !albumQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.AlbumScan(c, albumQuery, true, true, true))
}

func AlbumsLike(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if album exists
	albumQuery := database.Query("SELECT `id` FROM `albums` WHERE `id` = UUID_TO_BIN(?)", c.Params("albumID"))
	defer albumQuery.Close()
	if !albumQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if album_likes binding exists
	albumLikeQuery := database.Query("SELECT `id` FROM `album_likes` WHERE `album_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("albumID"), authUser.ID)
	defer albumLikeQuery.Close()
	if albumLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Create album_likes binding
	database.Exec("INSERT INTO `album_likes` (`id`, `album_id`, `user_id`) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), UUID_TO_BIN(?))", c.Params("albumID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}

func AlbumsLikeDelete(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if album exists
	albumQuery := database.Query("SELECT `id` FROM `albums` WHERE `id` = UUID_TO_BIN(?)", c.Params("albumID"))
	defer albumQuery.Close()
	if !albumQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if album_likes binding doesn't exists
	albumLikeQuery := database.Query("SELECT `id` FROM `album_likes` WHERE `album_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("albumID"), authUser.ID)
	defer albumLikeQuery.Close()
	if !albumLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete album_likes binding
	database.Exec("DELETE FROM `album_likes` WHERE `album_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("albumID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}
