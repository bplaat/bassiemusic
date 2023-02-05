package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

func DownloadArtist(c *fiber.Ctx) error {
	database.Exec("INSERT INTO `download_tasks` (`id`, `type`, `deezer_id`, `singles`) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?)",
		models.DownloadTaskTypeDeezerArtist, c.Query("deezer_id"), c.Query("singles"))

	return c.JSON(fiber.Map{"success": true})
}

func DownloadAlbum(c *fiber.Ctx) error {
	database.Exec("INSERT INTO `download_tasks` (`id`, `type`, `deezer_id`) VALUES (UUID_TO_BIN(UUID()), ?, ?)",
		models.DownloadTaskTypeDeezerAlbum, c.Query("deezer_id"))

	return c.JSON(fiber.Map{"success": true})
}
