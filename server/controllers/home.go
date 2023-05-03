package controllers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":    os.Getenv("APP_NAME"),
		"version": os.Getenv("APP_VERSION"),
	})
}
