package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"strings"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Bogon    bool   `json:"bogon"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

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
	user := models.UserModel().Where("username", params.Logon).WhereOr("email", params.Logon).First()
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
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		log.Fatalln(err)
	}
	token := base64.StdEncoding.EncodeToString(randomBytes)

	// Fetch ip info
	var ipInfo IPInfo
	utils.FetchJson("https://ipinfo.io/"+c.IP()+"/json", &ipInfo)

	// Create new session
	agent := utils.ParseUserAgent(c)
	session := database.Map{
		"user_id":        user.ID,
		"token":          token,
		"ip":             ipInfo.IP,
		"client_os":      &agent.OS,
		"client_name":    &agent.Name,
		"client_version": &agent.Version,
		"expires_at":     time.Now().Add(365 * 24 * 60 * 60 * time.Second),
	}
	if !ipInfo.Bogon {
		session["ip_latitude"] = strings.Split(ipInfo.Loc, ",")[0]
		session["ip_longitude"] = strings.Split(ipInfo.Loc, ",")[1]
		session["ip_country"] = ipInfo.Country
		session["ip_city"] = ipInfo.City
	}
	models.SessionModel().Create(session)

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

func AuthValidate(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Get session agent
	session := models.AuthSession(c)
	agent := utils.Agent{OS: *session.ClientOS, Name: *session.ClientName, Version: *session.ClientVersion}

	// Get last track play
	lastTackplay := models.TrackPlayModel().Where("user_id", authUser.ID).OrderByDesc("created_at").First()

	// When we have a last played track get it
	if lastTackplay != nil {
		return c.JSON(fiber.Map{
			"success":             true,
			"user":                authUser,
			"session":             session,
			"agent":               agent,
			"last_track":          models.TrackModel(c).With("artists", "album").Find(lastTackplay.TrackID),
			"last_track_position": lastTackplay.Position,
		})
	}

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"user":    authUser,
		"session": session,
		"agent":   agent,
	})
}

func AuthLogout(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)

	// Get session
	session := models.SessionModel().Where("token", token).First()
	if session == nil {
		return c.JSON(fiber.Map{"success": false})
	}

	// Revoke session
	models.SessionModel().Where("id", session.ID).Update(database.Map{
		"expires_at": time.Now(),
	})
	return c.JSON(fiber.Map{"success": true})
}
