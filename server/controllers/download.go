package controllers

import (
	"encoding/json"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type DownloadArtistParams struct {
	DeezerID    string `form:"deezer_id" validate:"required,numeric"`
	DisplayName string `form:"display_name" validate:"required"`
}

func DownloadArtist(c *fiber.Ctx) error {
	// Parse body
	var params DownloadArtistParams
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate values
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Create download task
	newTask := database.Map{
		"type":         models.DownloadTaskTypeDeezerArtist,
		"deezer_id":    params.DeezerID,
		"display_name": params.DisplayName,
	}

	task := models.DownloadTaskModel().Create(newTask)

	// Send new task to all admins
	data, _ := json.Marshal(fiber.Map{
		"success": true,
		"type":    "newTask",
		"data":    task,
	})

	SendMessageToAll(data)

	return c.JSON(fiber.Map{"success": true})
}

type DownloadAlbumParams struct {
	DeezerID    string `form:"deezer_id" validate:"required,numeric"`
	DisplayName string `form:"display_name" validate:"required"`
}

func DownloadAlbum(c *fiber.Ctx) error {
	// Parse body
	var params DownloadAlbumParams
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate values
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Create download task
	newTask := database.Map{
		"type":         models.DownloadTaskTypeDeezerAlbum,
		"deezer_id":    params.DeezerID,
		"display_name": params.DisplayName,
	}

	task := models.DownloadTaskModel().Create(newTask)

	// Send tasks to all admins
	data, _ := json.Marshal(fiber.Map{
		"success": true,
		"type":    "newTask",
		"data":    task,
	})

	SendMessageToAll(data)

	return c.JSON(fiber.Map{"success": true})
}
