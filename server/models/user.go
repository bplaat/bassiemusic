package models

import (
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `column:"id,uuid" json:"id"`
	Username  string    `column:"username,string" json:"username"`
	Email     string    `column:"email,string" json:"email"`
	Password  string    `column:"password,string" json:"-"`
	AvatarID  string    `column:"avatar_id,uuid" json:"-"`
	Avatar    string    `json:"avatar,omitempty"`
	Role      string    `column:"role,enum:normal|admin" json:"role"`
	Theme     string    `column:"theme,enum:system|light|dark" json:"theme"`
	CreatedAt time.Time `column:"created_at,timestamp" json:"created_at"`
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
			if c != nil && user.AvatarID != "" {
				user.Avatar = fmt.Sprintf("%s/storage/avatars/%s.jpg", c.BaseURL(), user.AvatarID)
			}
		},
	}.Init()
}
