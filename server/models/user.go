package models

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole int

const UserRoleNormal UserRole = 0
const UserRoleAdmin UserRole = 1

type UserTheme int

const UserThemeSystem UserTheme = 0
const UserThemeLight UserTheme = 1
const UserThemeDark UserTheme = 2

func UserScan(c *fiber.Ctx, userQuery *sql.Rows) User {
	var user User
	var userRole UserRole
	var userTheme UserTheme
	userQuery.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &userRole, &userTheme, &user.CreatedAt)

	if userRole == UserRoleNormal {
		user.Role = "normal"
	}
	if userRole == UserRoleAdmin {
		user.Role = "admin"
	}

	if userTheme == UserThemeSystem {
		user.Theme = "system"
	}
	if userTheme == UserThemeLight {
		user.Theme = "light"
	}
	if userTheme == UserThemeDark {
		user.Theme = "dark"
	}
	return user
}

func UsersScan(c *fiber.Ctx, usersQuery *sql.Rows) []User {
	users := []User{}
	for usersQuery.Next() {
		users = append(users, UserScan(c, usersQuery))
	}
	return users
}
