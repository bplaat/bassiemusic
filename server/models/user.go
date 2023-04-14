package models

import (
	"fmt"
	"os"

	"github.com/bplaat/bassiemusic/database"
)

type User struct {
	ID            string    `column:"id,uuid" json:"id"`
	Username      string    `column:"username,string" json:"username"`
	Email         string    `column:"email,string" json:"email"`
	Password      string    `column:"password,string" json:"-"`
	AvatarID      *string   `column:"avatar,uuid" json:"-"`
	SmallAvatar   *string   `json:"small_avatar"`
	MediumAvatar  *string   `json:"medium_avatar"`
	AllowExplicit bool      `column:"allow_explicit,bool" json:"allow_explicit"`
	Role          UserRole  `column:"role,int" json:"-"`
	RoleString    string    `json:"role"`
	Language      string    `column:"language,string" json:"language"`
	Theme         UserTheme `column:"theme,int" json:"-"`
	ThemeString   string    `json:"theme"`
	CreatedAt     string    `column:"created_at,timestamp" json:"created_at"`
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

		if user.AvatarID != nil && *user.AvatarID != "" {
			if _, err := os.Stat(fmt.Sprintf("storage/avatars/original/%s", *user.AvatarID)); err == nil {
				smallAvatar := fmt.Sprintf("%s/avatars/small/%s.jpg", os.Getenv("STORAGE_URL"), *user.AvatarID)
				user.SmallAvatar = &smallAvatar
				mediumAvatar := fmt.Sprintf("%s/avatars/medium/%s.jpg", os.Getenv("STORAGE_URL"), *user.AvatarID)
				user.MediumAvatar = &mediumAvatar
			}
		}
	},
}).Init()
