package controllers

import (
	"log"
	"strconv"

	"github.com/bplaat/bassiemusic/controllers/websocket"
	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

type DownloadArtistBody struct {
	DeezerID    string `form:"deezer_id" validate:"required|integer"`
	DisplayName string `form:"display_name" validate:"required"`
}

func DownloadArtist(c *fiber.Ctx) error {
	// Parse body
	var body DownloadArtistBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create download task
	deezerID, _ := strconv.ParseInt(body.DeezerID, 10, 64)
	downloadTask := models.DownloadTaskModel.Create(database.Map{
		"type":         models.DownloadTaskTypeDeezerArtist,
		"deezer_id":    deezerID,
		"display_name": body.DisplayName,
		"status":       models.DownloadTaskStatusPending,
		"progress":     0,
	})

	// Broadcast new task message to all listening admins
	if err := websocket.BroadcastAdmin("download_tasks.create", downloadTask); err != nil {
		log.Println(err)
	}

	// Return success response
	return c.JSON(fiber.Map{"success": true})
}

type DownloadAlbumBody struct {
	DeezerID    string `form:"deezer_id" validate:"required|integer"`
	DisplayName string `form:"display_name" validate:"required"`
}

func DownloadAlbum(c *fiber.Ctx) error {
	// Parse body
	var body DownloadAlbumBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create download task
	deezerID, _ := strconv.ParseInt(body.DeezerID, 10, 64)
	downloadTask := models.DownloadTaskModel.Create(database.Map{
		"type":         models.DownloadTaskTypeDeezerAlbum,
		"deezer_id":    deezerID,
		"display_name": body.DisplayName,
		"status":       models.DownloadTaskStatusPending,
		"progress":     0,
	})

	// Broadcast create task message to all admins
	if err := websocket.BroadcastAdmin("download_tasks.create", downloadTask); err != nil {
		log.Println(err)
	}

	// Return success response
	return c.JSON(fiber.Map{"success": true})
}
