package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total users count
	usersCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `users` WHERE `username` LIKE ? OR `email` LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer usersCountQuery.Close()

	usersCountQuery.Next()
	var total int64
	usersCountQuery.Scan(&total)

	// Get users
	var usersQuery *sql.Rows
	usersQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `username` LIKE ? OR `email` LIKE ? ORDER BY LOWER(`username`) LIMIT ?, ?", "%"+query+"%", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer usersQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.UsersScan(c, usersQuery),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersShow(c *fiber.Ctx) error {
	userQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("userID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	if !userQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.UserScan(c, userQuery))
}

func UsersSessions(c *fiber.Ctx) error {
	userQuery, err := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("userID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer userQuery.Close()

	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	var sessionsQuery *sql.Rows
	sessionsQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`user_id`), `token`, `ip`, `ip_latitude`, `ip_longitude`, `ip_country`, `ip_city`, `client_os`, `client_name`, `client_version`, `expires_at`, `created_at` FROM `sessions` WHERE `id` = UUID_TO_BIN(?) ORDER BY `created_at`", c.Params("userID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer sessionsQuery.Close()
	return c.JSON(models.SessionsScan(c, sessionsQuery, false))
}
