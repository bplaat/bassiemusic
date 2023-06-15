package models

import (
	"fmt"
	"os"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type User struct {
	ID            uuid.Uuid     `column:"id" json:"id"`
	Username      string        `column:"username" json:"username"`
	Email         string        `column:"email" json:"email"`
	Password      string        `column:"password" json:"-"`
	AvatarID      uuid.NullUuid `column:"avatar" json:"-"`
	AllowExplicit bool          `column:"allow_explicit" json:"allow_explicit"`
	Role          UserRole      `column:"role" json:"-"`
	RoleString    string        `json:"role"`
	Language      string        `column:"language" json:"language"`
	Theme         UserTheme     `column:"theme" json:"-"`
	ThemeString   string        `json:"theme"`
	CreatedAt     string        `column:"created_at" json:"created_at"`
	SmallAvatar   *string       `json:"small_avatar"`
	MediumAvatar  *string       `json:"medium_avatar"`
}

type UserRole int

const UserRoleNormal UserRole = 0
const UserRoleAdmin UserRole = 1

type UserTheme int

const UserThemeSystem UserTheme = 0
const UserThemeLight UserTheme = 1
const UserThemeDark UserTheme = 2

var UserModel *database.Model[User] = (&database.Model[User]{
	TableName: "users",
	Process: func(user *User) {
		if user.Role == UserRoleNormal {
			user.RoleString = "normal"
		}
		if user.Role == UserRoleAdmin {
			user.RoleString = "admin"
		}

		if user.Theme == UserThemeSystem {
			user.ThemeString = "system"
		}
		if user.Theme == UserThemeLight {
			user.ThemeString = "light"
		}
		if user.Theme == UserThemeDark {
			user.ThemeString = "dark"
		}

		if user.AvatarID.Valid {
			avatarIDString := user.AvatarID.Uuid.String()
			if _, err := os.Stat(fmt.Sprintf("storage/avatars/original/%s", avatarIDString)); err == nil {
				smallAvatar := fmt.Sprintf("%s/avatars/small/%s.jpg", os.Getenv("STORAGE_URL"), avatarIDString)
				user.SmallAvatar = &smallAvatar
				mediumAvatar := fmt.Sprintf("%s/avatars/medium/%s.jpg", os.Getenv("STORAGE_URL"), avatarIDString)
				user.MediumAvatar = &mediumAvatar
			}
		}
	},
}).Init()
