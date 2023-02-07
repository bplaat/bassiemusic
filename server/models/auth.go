package models

// This needs to be in the utils package but then you get an import cycle error
// So to put it here is a quick and dirty fix

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseTokenVar(c *fiber.Ctx) string {
	if strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
		return c.Get("Authorization")[7:]
	}
	return ""
}

func AuthUser(c *fiber.Ctx) User {
	token := ParseTokenVar(c)
	return SessionModel(c).With("user").Where("token", token).First().User
}
