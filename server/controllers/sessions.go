package controllers

import (
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SessionsIndex(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.SessionModel().With("user").OrderByDesc("expires_at").Paginate(page, limit))
}

func SessionsShow(c *fiber.Ctx) error {
	session := models.SessionModel().With("user").Find(c.Params("sessionID"))
	if session == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(session)
}

func SessionsRevoke(c *fiber.Ctx) error {
	// Get session
	session := models.SessionModel().Find(c.Params("sessionID"))
	if session == nil {
		return fiber.ErrNotFound
	}

	// Revoke session
	models.SessionModel().Where("id", session.ID).Update(database.Map{
		"expires_at": time.Now(),
	})
	return c.JSON(fiber.Map{"success": true})
}
