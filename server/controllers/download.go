package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

func DownloadArtist(c *fiber.Ctx) error {
	models.DownloadTaskModel().Create(database.Map{
		"type":      models.DownloadTaskTypeDeezerArtist,
		"deezer_id": c.Query("deezer_id"),
	})
	return c.JSON(fiber.Map{"success": true})
}

func DownloadAlbum(c *fiber.Ctx) error {
	models.DownloadTaskModel().Create(database.Map{
		"type":      models.DownloadTaskTypeDeezerAlbum,
		"deezer_id": c.Query("deezer_id"),
	})
	return c.JSON(fiber.Map{"success": true})
}
