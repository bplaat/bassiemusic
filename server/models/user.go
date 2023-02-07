package models

import (
	"fmt"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `column:"id,uuid" json:"id"`
	Username  string    `column:"username,string" json:"username"`
	Email     string    `column:"email,string" json:"email"`
	Password  string    `column:"password,string" json:"-"`
	AvatarID  *string   `column:"avatar,uuid" json:"-"`
	Avatar    string    `json:"avatar,omitempty"`
	RoleInt   UserRole  `column:"role,int" json:"-"`
	Role      string    `json:"role"`
	ThemeInt  UserTheme `column:"theme,int" json:"-"`
	Theme     string    `json:"theme"`
	CreatedAt string    `column:"created_at,timestamp" json:"created_at"`
}

type UserRole int

const UserRoleNormal UserRole = 0
const UserRoleAdmin UserRole = 1

type UserTheme int

const UserThemeSystem UserTheme = 0
const UserThemeLight UserTheme = 1
const UserThemeDark UserTheme = 2

func UserModel(c *fiber.Ctx) database.Model[User] {
	return database.Model[User]{
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

			if c != nil && user.AvatarID != nil {
				user.Avatar = fmt.Sprintf("%s/storage/avatars/%s.jpg", c.BaseURL(), *user.AvatarID)
			}
		},
	}.Init()
}
