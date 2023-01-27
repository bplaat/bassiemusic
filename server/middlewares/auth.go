package middlewares

import (
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func IsAuthed(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)
	if token == "" {
		return fiber.ErrUnauthorized
	}

	// Get all active session by token
	sessionQuery, err := database.Query("SELECT `id` FROM `sessions` WHERE `token` = ? AND `expires_at` > NOW() LIMIT 1", token)
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionQuery.Close()

	// When a session doesn't exist return unauthorized error
	if !sessionQuery.Next() {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	user := utils.AuthUser(c)
	if user.Role != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
