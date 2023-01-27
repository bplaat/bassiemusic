package utils

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseTokenVar(c *fiber.Ctx) string {
	if strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
		return c.Get("Authorization")[7:]
	}
	return ""
}

func ParseIndexVars(c *fiber.Ctx) (string, int, int) {
	query := c.Query("q")

	page := 1
	if pageInt, err := strconv.Atoi(c.Query("page")); err == nil {
		page = pageInt
		if page < 1 {
			page = 1
		}
	}

	limit := 20
	if limitInt, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = limitInt
		if limit < 1 {
			limit = 1
		}
		if limit > 50 {
			limit = 50
		}
	}

	return query, page, limit
}
