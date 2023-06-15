package controllers

import (
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SessionsIndex(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.SessionModel.With("user").OrderByDesc("created_at").Paginate(page, limit))
}

func SessionsShow(c *fiber.Ctx) error {
	// Parse session id uuid
	sessionID, err := uuid.Parse(c.Params("sessionID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if session exists
	session := models.SessionModel.With("user").Find(sessionID)
	if session == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(session)
}

func SessionsRevoke(c *fiber.Ctx) error {
	// Parse session id uuid
	sessionID, err := uuid.Parse(c.Params("sessionID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if session exists
	session := models.SessionModel.Find(sessionID)
	if session == nil {
		return fiber.ErrNotFound
	}

	// Revoke session
	models.SessionModel.Where("id", session.ID).Update(database.Map{
		"expires_at": time.Now(),
	})
	return c.JSON(fiber.Map{"success": true})
}
