package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"strings"
	"time"

	"github.com/bplaat/bassiemusic/consts"
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

func getIP(c *fiber.Ctx) string {
	if ip, ok := c.GetReqHeaders()["X-Forwarded-For"]; ok {
		return ip
	}
	return c.IP()
}

type AuthLoginBody struct {
	Logon    string `form:"logon"`
	Password string `form:"password"`
}

func AuthLogin(c *fiber.Ctx) error {
	// Parse body
	var body AuthLoginBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Get user by username or email
	user := models.UserModel.Where("username", body.Logon).WhereOr("email", body.Logon).First()
	if user == nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong username, email or password",
		})
	}

	// Verify user password
	if !utils.VerifyPassword(body.Password, user.Password) {
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
	if err := utils.FetchJson("https://ipinfo.io/"+getIP(c)+"/json", &ipInfo); err != nil {
		log.Fatalln(err)
	}

	// Create new session
	agent := utils.ParseUserAgent(c)
	session := database.Map{
		"user_id":        user.ID,
		"token":          token,
		"ip":             ipInfo.IP,
		"client_os":      &agent.OS,
		"client_name":    &agent.Name,
		"client_version": &agent.Version,
		"expires_at":     time.Now().Add(consts.AUTH_TOKEN_EXPIRES_TIMEOUT),
	}
	if !ipInfo.Bogon {
		session["ip_latitude"] = strings.Split(ipInfo.Loc, ",")[0]
		session["ip_longitude"] = strings.Split(ipInfo.Loc, ",")[1]
		session["ip_country"] = strings.ToLower(ipInfo.Country)
		session["ip_city"] = ipInfo.City
	}
	models.SessionModel.Create(session)

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

func AuthValidate(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Get session agent
	session := c.Locals("session").(*models.Session)
	agent := utils.Agent{OS: *session.ClientOS, Name: *session.ClientName, Version: *session.ClientVersion}
	response := fiber.Map{
		"success":     true,
		"user":        authUser,
		"session_id	": session.ID,
		"agent":       agent,
	}

	// Get last played track and return it
	lastTackPlay := models.TrackPlayModel.Where("user_id", authUser.ID).OrderByDesc("created_at").First()
	if lastTackPlay != nil {
		response["last_track"] = models.TrackModel.WithArgs("liked", c.Locals("authUser")).
			With("artists", "album").Find(lastTackPlay.TrackID)
		response["last_track_position"] = lastTackPlay.Position
	}

	// Get last playlists
	response["last_playlists"] = models.PlaylistModel.Where("user_id", authUser.ID).OrderByDesc("updated_at").Limit(10).Get()

	// Return response
	return c.JSON(response)
}

func AuthLogout(c *fiber.Ctx) error {
	token := utils.ParseTokenVar(c)

	// Get session
	session := models.SessionModel.Where("token", token).First()
	if session == nil {
		return c.JSON(fiber.Map{"success": false})
	}

	// Revoke session
	models.SessionModel.Where("id", session.ID).Update(database.Map{
		"expires_at": time.Now(),
	})
	return c.JSON(fiber.Map{"success": true})
}
