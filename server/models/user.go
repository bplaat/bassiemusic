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
	Avatar        string    `json:"avatar,omitempty"`
	AllowExplicit bool      `column:"allow_explicit,bool" json:"allow_explicit"`
	RoleInt       UserRole  `column:"role,int" json:"-"`
	Role          string    `json:"role"`
	ThemeInt      UserTheme `column:"theme,int" json:"-"`
	Language      string    `column:"language,string" json:"language"`
	Theme         string    `json:"theme"`
	CreatedAt     string    `column:"created_at,timestamp" json:"created_at"`
}

type UserRole int

const UserRoleNormal UserRole = 0
const UserRoleAdmin UserRole = 1

type UserTheme int

const UserThemeSystem UserTheme = 0
const UserThemeLight UserTheme = 1
const UserThemeDark UserTheme = 2

func UserModel() *database.Model[User] {
	return (&database.Model[User]{
		TableName: "users",
		Process: func(user *User) {
			if user.RoleInt == UserRoleNormal {
				user.Role = "normal"
			}
			if user.RoleInt == UserRoleAdmin {
				user.Role = "admin"
			}

			if user.ThemeInt == UserThemeSystem {
				user.Theme = "system"
			}
			if user.ThemeInt == UserThemeLight {
				user.Theme = "light"
			}
			if user.ThemeInt == UserThemeDark {
				user.Theme = "dark"
			}

			if user.AvatarID != nil && *user.AvatarID != "" {
				user.Avatar = fmt.Sprintf("%s/avatars/%s.jpg", os.Getenv("STORAGE_URL"), *user.AvatarID)
			}
		},
	}).Init()
}
