package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

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
	user := models.UserModel(c).Where("username", params.Logon).WhereOr("email", params.Logon).First()
	if user == nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong username, email or password",
		})
	}

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
	agent := useragent.Parse(c.Get("User-Agent"))
	session := models.Session{
		UserID:        user.ID,
		Token:         token,
		IP:            c.IP(),
		ClientOS:      agent.OS,
		ClientName:    agent.Name,
		ClientVersion: agent.Version,
		ExpiresAt:     time.Now().Add(365 * 24 * 60 * 60 * time.Second),
	}
	models.SessionModel(c).Create(&session)

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

func AuthValidate(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Get last track play
	lastTackplay := models.TrackPlayModel().Where("user_id", authUser.ID).OrderByDesc("created_at").First()

	// When we have a last played track get it
	if lastTackplay != nil {
		return c.JSON(fiber.Map{
			"success":             true,
			"user":                authUser,
			"last_track":          models.TrackModel(c).With("artists", "album").Find(lastTackplay.TrackID),
			"last_track_position": lastTackplay.Position,
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

	// Get session
	session := models.SessionModel(c).Where("token", token).First()
	if session == nil {
		return c.JSON(fiber.Map{"success": false})
	}

	// Revoke session
	session.ExpiresAt = time.Now()
	models.SessionModel(c).Update(session)
	return c.JSON(fiber.Map{"success": true})
}
