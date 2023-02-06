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
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, BIN_TO_UUID(`avatar`), `role`, `theme`, `created_at` FROM `users` WHERE `username` = ? OR `email` = ?", params.Logon, params.Logon)
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
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		log.Fatalln(err)
	}
	token := base64.StdEncoding.EncodeToString(randomBytes)

	// Create new session
	ua := useragent.Parse(c.Get("User-Agent"))
	database.Exec("INSERT INTO `sessions` (`id`, `user_id`, `token`, `ip`, `client_os`, `client_name`, `client_version`, `expires_at`) VALUES "+
		"(UUID_TO_BIN(UUID()), UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?)",
		user.ID, token, c.IP(), ua.OS, ua.Name, ua.Version, time.Now().Add(365*24*60*60*time.Second).Format(time.DateTime))

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

func AuthValidate(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Get last track plays binding
	trackPlayQuery := database.Query("SELECT BIN_TO_UUID(`track_id`), `position` FROM `track_plays` WHERE `user_id` = UUID_TO_BIN(?) ORDER BY `created_at` DESC LIMIT 1", authUser.ID)
	defer trackPlayQuery.Close()

	// When we have a last played track get it
	if trackPlayQuery.Next() {
		var lastTrackID string
		var lastTrackPosition float32
		trackPlayQuery.Scan(&lastTrackID, &lastTrackPosition)
		trackQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", lastTrackID)
		defer trackQuery.Close()
		trackQuery.Next()

		// Return response
		return c.JSON(fiber.Map{
			"success":             true,
			"user":                authUser,
			"last_track":          models.TrackScan(c, trackQuery, true, true),
			"last_track_position": lastTrackPosition,
		})
	}

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"user":    authUser,
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
	return c.JSON(fiber.Map{"success": true})
}
