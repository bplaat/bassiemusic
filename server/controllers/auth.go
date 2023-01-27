package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
)

type AuthLoginParams struct {
	Logon    string `form:"logon"`
	Password string `form:"password"`
}

func AuthLogin(c *fiber.Ctx) error {
	var params AuthLoginParams
	if err := c.BodyParser(&params); err != nil {
		log.Fatalln(err)
	}

	// Get user by username or email
	userQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `username` = ? OR `email` = ?", params.Logon, params.Logon)
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	if !userQuery.Next() {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong username, email or password",
		})
	}
	user := models.UserScan(c, userQuery)

	// Verify user password
	if !utils.VerifyPassword(params.Password, user.Password) {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong email or password",
		})
	}

	// Generate new token
	randomBytes := make([]byte, 128)
	_, err = io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		log.Fatalln(err)
	}
	token := base64.StdEncoding.EncodeToString(randomBytes)

	// Create new session
	ua := useragent.Parse(c.Get("User-Agent"))
	_, err = database.Exec("INSERT INTO `sessions` (`id`, `user_id`, `token`, `ip`, `client_os`, `client_name`, `client_version`, `expires_at`) VALUES "+
		"(UUID_TO_BIN(UUID()), UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?)",
		user.ID, token, c.IP(), ua.OS, ua.Name, ua.Version, time.Now().Add(365*24*60*60*time.Second).Format(time.RFC3339))
	if err != nil {
		log.Fatalln(err)
	}

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

func AuthValidate(c *fiber.Ctx) error {
	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"user":    utils.AuthUser(c),
	})
}

func AuthLogout(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)
	if token == "" {
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	// Revoke session
	database.Exec("UPDATE `sessions` SET `expires_at` = NOW() WHERE `token` = ?", token)

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
	})
}
