package controllers

import (
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

func DownloadArtist(c *fiber.Ctx) error {
	deezerID := c.Query("deezer_id")
	Singles := c.Query("singles")

	_, err := database.Exec("INSERT INTO `download_tasks` (`id`, `type`, `deezer_id`, `singles`) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?)",
		models.DownloadTaskTypeDeezerArtist, deezerID, Singles)
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(&fiber.Map{
		"success": true,
	})
}

func DownloadAlbum(c *fiber.Ctx) error {
	deezerID := c.Query("deezer_id")

	_, err := database.Exec("INSERT INTO `download_tasks` (`id`, `type`, `deezer_id`) VALUES (UUID_TO_BIN(UUID()), ?, ?)",
		models.DownloadTaskTypeDeezerAlbum, deezerID)
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(&fiber.Map{
		"success": true,
	})
}
