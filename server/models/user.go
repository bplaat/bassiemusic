package models

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" form:"username"`
	Email     string    `json:"email" form:"email"`
	Password  string    `json:"-" form:"password"`
	Role      string    `json:"role" form:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole int

const UserRoleNormal UserRole = 0
const UserRoleAdmin UserRole = 1

func UserScan(c *fiber.Ctx, userQuery *sql.Rows) User {
	var user User
	var userRole UserRole
	userQuery.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &userRole, &user.CreatedAt)
	if userRole == UserRoleNormal {
		user.Role = "normal"
	}
	if userRole == UserRoleAdmin {
		user.Role = "admin"
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
