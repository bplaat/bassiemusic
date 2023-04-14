package models

import (
	"time"

	"github.com/bplaat/bassiemusic/core/database"
)

type Session struct {
	ID            string    `column:"id,uuid" json:"id"`
	UserID        string    `column:"user_id,uuid" json:"-"`
	Token         string    `column:"token,string" json:"-"`
	IP            string    `column:"ip,string" json:"ip"`
	IPLatitude    *float64  `column:"ip_latitude,string" json:"ip_latitude"`
	IPLongitude   *float64  `column:"ip_longitude,string" json:"ip_longitude"`
	IPCountry     *string   `column:"ip_country,string" json:"ip_country"`
	IPCity        *string   `column:"ip_city,string" json:"ip_city"`
	ClientOS      *string   `column:"client_os,string" json:"client_os"`
	ClientName    *string   `column:"client_name,string" json:"client_name"`
	ClientVersion *string   `column:"client_version,string" json:"client_version"`
	ExpiresAt     time.Time `column:"expires_at,timestamp" json:"expires_at"`
	CreatedAt     time.Time `column:"created_at,timestamp" json:"created_at"`
	User          *User     `json:"user,omitempty"`
}

var SessionModel *database.Model[Session] = (&database.Model[Session]{
	TableName: "sessions",
	Relationships: map[string]database.ModelRelationshipFunc[Session]{
		"user": func(session *Session, args []any) {
			session.User = UserModel.Find(session.UserID)
		},
	},
}).Init()
