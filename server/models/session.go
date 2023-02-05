package models

import (
	"database/sql"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	ID            string    `json:"id"`
	UserID        string    `json:"-"`
	Token         string    `json:"-"`
	IP            string    `json:"ip"`
	IPLatitude    string    `json:"ip_latitude"`
	IPLongitude   string    `json:"ip_longitude"`
	IPCountry     string    `json:"ip_country"`
	IPCity        string    `json:"ip_city"`
	ClientOS      string    `json:"client_os"`
	ClientName    string    `json:"client_name"`
	ClientVersion string    `json:"client_version"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	User          User      `json:"user,omitempty"`
}

func SessionScan(c *fiber.Ctx, sessionQuery *sql.Rows, withUser bool) Session {
	var session Session
	sessionQuery.Scan(&session.ID, &session.UserID, &session.Token, &session.IP, &session.IPLatitude,
		&session.IPLongitude, &session.IPCountry, &session.IPCity, &session.ClientOS, &session.ClientName,
		&session.ClientVersion, &session.ExpiresAt, &session.CreatedAt)
	if withUser {
		session.User = SessionUser(c, &session)
	}
	return session
}

func SessionsScan(c *fiber.Ctx, sessionsQuery *sql.Rows, withUser bool) []Session {
	sessions := []Session{}
	for sessionsQuery.Next() {
		sessions = append(sessions, SessionScan(c, sessionsQuery, withUser))
	}
	return sessions
}

func SessionUser(c *fiber.Ctx, session *Session) User {
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", session.UserID)
	defer userQuery.Close()
	userQuery.Next()
	return UserScan(c, userQuery)
}
