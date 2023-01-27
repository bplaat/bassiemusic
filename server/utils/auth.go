package utils

import (
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/gofiber/fiber/v2"
)

func AuthUser(c *fiber.Ctx) models.User {
	token := ParseTokenVar(c)

	// Get user from session
	sessionQuery, err := database.Query("SELECT BIN_TO_UUID(`user_id`) FROM `sessions` WHERE `token` = ? LIMIT 1", token)
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionQuery.Close()

	sessionQuery.Next()
	var userID string
	sessionQuery.Scan(&userID)

	// Get user
	userQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", userID)
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	userQuery.Next()
	return models.UserScan(c, userQuery)
}
