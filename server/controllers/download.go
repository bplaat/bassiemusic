package controllers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/bplaat/bassiemusic/controllers/websocket"
	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

type DownloadArtistBody struct {
	DeezerID    string `form:"deezer_id" validate:"required|integer"`
	DisplayName string `form:"display_name" validate:"required"`
}

type Data struct {
	DeezerID  int64  `json:"deezer_id"`
	YoutubeID string `json:"youtube_id"`
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
	data := Data{
		DeezerID: deezerID,
	}
	jsonData, _ := json.Marshal(data)
	downloadTask := models.DownloadTaskModel.Create(database.Map{
		"type":         models.DownloadTaskTypeDeezerArtist,
		"data":         jsonData,
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
	data := Data{
		DeezerID: deezerID,
	}
	jsonData, _ := json.Marshal(data)
	downloadTask := models.DownloadTaskModel.Create(database.Map{
		"type":         models.DownloadTaskTypeDeezerAlbum,
		"data":         jsonData,
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

func DownloadDelete(c *fiber.Ctx) error {
	// Parse download task id uuid
	downloadTaskID, err := uuid.Parse(c.Params("downloadTaskID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if task exists
	downloadTask := models.DownloadTaskModel.Find(downloadTaskID)
	if downloadTask == nil {
		return fiber.ErrNotFound
	}

	// Check if task is pending
	if downloadTask.Status != models.DownloadTaskStatusPending {
		return fiber.ErrBadRequest
	}

	// Broadcast delete task message to all admins
	if err := websocket.BroadcastAdmin("download_tasks.delete", downloadTask); err != nil {
		log.Println(err)
	}

	// Delete download task
	models.DownloadTaskModel.Where("id", downloadTaskID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
