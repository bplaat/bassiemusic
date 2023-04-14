package controllers

import (
	"strconv"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

type DownloadArtistBody struct {
	DeezerID string `form:"deezer_id" validate:"required|integer"`
}

func DownloadArtist(c *fiber.Ctx) error {
	// Parse body
	var body DownloadArtistBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.Validate(c, &body); err != nil {
		return err
	}

	// Create download task
	deezerID, _ := strconv.ParseInt(body.DeezerID, 10, 64)
	models.DownloadTaskModel.Create(database.Map{
		"type":      models.DownloadTaskTypeDeezerArtist,
		"deezer_id": deezerID,
	})
	return c.JSON(fiber.Map{"success": true})
}

type DownloadAlbumBody struct {
	DeezerID string `form:"deezer_id" validate:"required|integer"`
}

func DownloadAlbum(c *fiber.Ctx) error {
	// Parse body
	var body DownloadAlbumBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.Validate(c, &body); err != nil {
		return err
	}

	// Create download task
	deezerID, _ := strconv.ParseInt(body.DeezerID, 10, 64)
	models.DownloadTaskModel.Create(database.Map{
		"type":      models.DownloadTaskTypeDeezerAlbum,
		"deezer_id": deezerID,
	})
	return c.JSON(fiber.Map{"success": true})
}
