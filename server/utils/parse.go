package utils

import (
	"log"
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

func ParseIndexVars(c *fiber.Ctx) (string, int64, int64) {
	query := c.Query("q")

	page := int64(1)
	if pageInt, err := strconv.ParseInt(c.Query("page"), 10, 64); err == nil {
		page = pageInt
		if page < 1 {
			page = 1
		}
	}

	limit := int64(20)
	if limitInt, err := strconv.ParseInt(c.Query("limit"), 10, 64); err == nil {
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

func ParseBytes(bytes string) int64 {
	if strings.HasSuffix(bytes, "KB") {
		value, err := strconv.ParseInt(bytes[:len(bytes)-2], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		return value * 1000
	}
	if strings.HasSuffix(bytes, "MB") {
		value, err := strconv.ParseInt(bytes[:len(bytes)-2], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		return value * 1000000

	}
	if strings.HasSuffix(bytes, "GB") {
		value, err := strconv.ParseInt(bytes[:len(bytes)-2], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		return value * 1000000000

	}
	if strings.HasSuffix(bytes, "TB") {
		value, err := strconv.ParseInt(bytes[:len(bytes)-2], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		return value * 1000000000000
	}
	// Bytes
	value, err := strconv.ParseInt(bytes, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return value
}
