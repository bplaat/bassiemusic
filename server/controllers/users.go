package controllers

import (
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total users
	total := database.Count("SELECT COUNT(`id`) FROM `users` WHERE `username` LIKE ? OR `email` LIKE ?", "%"+query+"%", "%"+query+"%")

	// Get users
	usersQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `username` LIKE ? OR `email` LIKE ? ORDER BY LOWER(`username`) LIMIT ?, ?", "%"+query+"%", "%"+query+"%", (page-1)*limit, limit)
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

func UsersCreate(c *fiber.Ctx) error {
	// Parse and validate body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Fatalln(err)
	}

	// Create user
	userId := uuid.NewV4()
	var userRole models.UserRole
	if user.Role == "normal" {
		userRole = models.UserRoleNormal
	}
	if user.Role == "admin" {
		userRole = models.UserRoleAdmin
	}
	database.Exec("INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)",
		userId.String(), user.Username, user.Email, utils.HashPassword(user.Password), userRole)

	// Get new created user and send response
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", userId.String())
	defer userQuery.Close()
	return c.JSON(models.UserScan(c, userQuery))
}

func UsersShow(c *fiber.Ctx) error {
	// Check if user exists
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.UserScan(c, userQuery))
}

func UsersEdit(c *fiber.Ctx) error {
	// Check if user exists
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get user
	user := models.UserScan(c, userQuery)

	// Parse and validate body
	var updates models.User
	if err := c.BodyParser(&updates); err != nil {
		log.Fatalln(err)
	}

	// TODO

	// Get edited user and send response
	updatedUserQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", user.ID)
	defer updatedUserQuery.Close()
	return c.JSON(models.UserScan(c, updatedUserQuery))
}

func UsersSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get total user sessions
	total := database.Count("SELECT COUNT(`id`) FROM `sessions` WHERE `user_id` = UUID_TO_BIN(?)", c.Params("userID"))

	// Get user sessions
	sessionsQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`user_id`), `token`, `ip`, `ip_latitude`, `ip_longitude`, `ip_country`, `ip_city`, `client_os`, `client_name`, `client_version`, `expires_at`, `created_at` FROM `sessions` WHERE `user_id` = UUID_TO_BIN(?) ORDER BY `created_at` DESC LIMIT ?, ?", c.Params("userID"), (page-1)*limit, limit)
	defer sessionsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.UsersScan(c, sessionsQuery),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersDelete(c *fiber.Ctx) error {
	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Delete user
	database.Exec("DELETE FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))

	// Return response
	return c.JSON(&fiber.Map{
		"success": true,
	})
}
