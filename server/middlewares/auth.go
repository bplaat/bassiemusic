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
	session := models.SessionModel().Where("token", token).WhereRaw("`expires_at` > ?", time.Now()).First()
	if session == nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	user := models.AuthUser(c)
	if user.Role != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
