package models

import (
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	ID            string    `column:"id,uuid" json:"id"`
	UserID        string    `column:"user_id,uuid" json:"-"`
	Token         string    `column:"token_string" json:"-"`
	IP            string    `column:"ip,string" json:"ip"`
	IPLatitude    string    `column:"ip_latitude,string" json:"ip_latitude"`
	IPLongitude   string    `column:"ip_longitude,string" json:"ip_longitude"`
	IPCountry     string    `column:"ip_country,string" json:"ip_country"`
	IPCity        string    `column:"ip_city,string" json:"ip_city"`
	ClientOS      string    `column:"client_os,string" json:"client_os"`
	ClientName    string    `column:"client_name,string" json:"client_name"`
	ClientVersion string    `column:"client_version,string" json:"client_version"`
	ExpiresAt     time.Time `column:"expires_at,timestamp" json:"expires_at"`
	CreatedAt     time.Time `column:"created_at,timestamp" json:"created_at"`
	User          User      `json:"user,omitempty"`
}

func SessionModel(c *fiber.Ctx) database.Model[Session] {
	return database.Model[Session]{
		TableName: "sessions",
		Relationships: map[string]database.QueryBuilderProcess[Session]{
			"user": func(session *Session) {
				session.User = *UserModel(c).Find(session.UserID)
			},
		},
	}.Init()
}
