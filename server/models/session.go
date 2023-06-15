package models

import (
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type Session struct {
	ID            uuid.Uuid            `column:"id" json:"id"`
	UserID        uuid.Uuid            `column:"user_id" json:"-"`
	Token         string               `column:"token" json:"-"`
	IP            string               `column:"ip" json:"ip"`
	IPLatitude    database.NullFloat64 `column:"ip_latitude" json:"ip_latitude"`
	IPLongitude   database.NullFloat64 `column:"ip_longitude" json:"ip_longitude"`
	IPCountry     database.NullString  `column:"ip_country" json:"ip_country"`
	IPCity        database.NullString  `column:"ip_city" json:"ip_city"`
	ClientOS      database.NullString  `column:"client_os" json:"client_os"`
	ClientName    database.NullString  `column:"client_name" json:"client_name"`
	ClientVersion database.NullString  `column:"client_version" json:"client_version"`
	ExpiresAt     time.Time            `column:"expires_at" json:"expires_at"`
	CreatedAt     time.Time            `column:"created_at" json:"created_at"`
	User          *User                `json:"user,omitempty"`
}

var SessionModel *database.Model[Session] = (&database.Model[Session]{
	TableName: "sessions",
	Relationships: map[string]database.ModelRelationshipFunc[Session]{
		"user": func(session *Session, args []any) {
			session.User = UserModel.Find(session.UserID)
		},
	},
}).Init()
