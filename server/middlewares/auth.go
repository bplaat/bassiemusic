package middlewares

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func IsAuthed(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// Get all active session by token
	sessionQuery := database.Query("SELECT `id` FROM `sessions` WHERE `token` = ? AND `expires_at` > NOW()", token)
	defer sessionQuery.Close()

	// When a session doesn't exist return unauthorized error
	if !sessionQuery.Next() {
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
