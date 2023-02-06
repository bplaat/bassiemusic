package models

// This needs to be in the utils package but then you get an import cycle error
// So to put it here is a quick and dirty fix

import (
	"strings"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

func ParseTokenVar(c *fiber.Ctx) string {
	if strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
		return c.Get("Authorization")[7:]
	}
	return ""
}

func AuthUser(c *fiber.Ctx) User {
	token := ParseTokenVar(c)

	// Get user from session
	sessionQuery := database.Query("SELECT BIN_TO_UUID(`user_id`) FROM `sessions` WHERE `token` = ?", token)
	defer sessionQuery.Close()
	sessionQuery.Next()
	var userID string
	sessionQuery.Scan(&userID)

	// Get user
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, BIN_TO_UUID(`avatar`), `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", userID)
	defer userQuery.Close()
	userQuery.Next()
	return UserScan(c, userQuery)
}
