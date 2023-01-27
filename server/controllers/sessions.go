package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func SessionsIndex(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Get total sessions count
	sessionsCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `sessions`")
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionsCountQuery.Close()

	sessionsCountQuery.Next()
	var total int64
	sessionsCountQuery.Scan(&total)

	// Get sessions
	var sessionsQuery *sql.Rows
	sessionsQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`user_id`), `token`, `ip`, `ip_latitude`, `ip_longitude`, `ip_country`, `ip_city`, `client_os`, `client_name`, `client_version`, `expires_at`, `created_at` FROM `sessions` ORDER BY `created_at` LIMIT ?, ?", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.SessionsScan(c, sessionsQuery, true),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func SessionsShow(c *fiber.Ctx) error {
	sessionQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`user_id`), `token`, `ip`, `ip_latitude`, `ip_longitude`, `ip_country`, `ip_city`, `client_os`, `client_name`, `client_version`, `expires_at`, `created_at` FROM `sessions` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("sessionID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionQuery.Close()

	if !sessionQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.SessionScan(c, sessionQuery, true))
}

func SessionsRevoke(c *fiber.Ctx) error {
	// Get session
	sessionQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`user_id`), `token`, `ip`, `ip_latitude`, `ip_longitude`, `ip_country`, `ip_city`, `client_os`, `client_name`, `client_version`, `expires_at`, `created_at` FROM `sessions` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("sessionID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionQuery.Close()

	if !sessionQuery.Next() {
		return fiber.ErrNotFound
	}
	session := models.SessionScan(c, sessionQuery, true)

	// Revoke session
	database.Exec("UPDATE `sessions` SET `expires_at` = NOW() WHERE `id` = UUID_TO_BIN(?)", session.ID)

	// Return response
	return c.JSON(fiber.Map{
		"success": true,
	})
}
