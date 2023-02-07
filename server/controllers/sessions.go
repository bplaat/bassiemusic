package controllers

import (
	"time"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SessionsIndex(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.SessionModel(c).Paginate(page, limit))
}

func SessionsShow(c *fiber.Ctx) error {
	session := models.SessionModel(c).With("albums").Find(c.Params("sessionID"))
	if session == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(session)
}

func SessionsRevoke(c *fiber.Ctx) error {
	// Get session
	session := models.SessionModel(c).With("albums").Find(c.Params("sessionID"))
	if session == nil {
		return fiber.ErrNotFound
	}

	// Revoke session
	session.ExpiresAt = time.Now()
	models.SessionModel(c).Update(session)
	return c.JSON(fiber.Map{"success": true})
}
