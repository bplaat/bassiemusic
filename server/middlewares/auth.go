package middlewares

import (
	"time"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func IsAuthed(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// Get active session by token
	session := models.SessionModel.With("user").Where("token", token).WhereRaw("`expires_at` > ?", time.Now()).First()
	if session == nil {
		return fiber.ErrUnauthorized
	}
	c.Locals("session", session)
	c.Locals("authUser", session.User)
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
