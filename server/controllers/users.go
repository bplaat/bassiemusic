package controllers

import (
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total users
	total := database.Count("SELECT COUNT(`id`) FROM `users` WHERE `username` LIKE ? OR `email` LIKE ?", "%"+query+"%", "%"+query+"%")

	// Get users
	usersQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `username` LIKE ? OR `email` LIKE ? ORDER BY LOWER(`username`) LIMIT ?, ?", "%"+query+"%", "%"+query+"%", (page-1)*limit, limit)
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

type UsersCreateParams struct {
	Username string `form:"username" validate:"required,min=2"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
	Role     string `form:"role" validate:"required"`
	Theme    string `form:"theme" validate:"required"`
}

func UsersCreate(c *fiber.Ctx) error {
	// Parse body
	var params UsersCreateParams
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate values
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate username is unique
	usernameQuery := database.Query("SELECT `id` FROM `users` WHERE `username` = ?", params.Username)
	defer usernameQuery.Close()
	if usernameQuery.Next() {
		log.Println("username not unique")
		return fiber.ErrBadRequest
	}

	// Validate email is unique
	emailQuery := database.Query("SELECT `id` FROM `users` WHERE `email` = ?", params.Email)
	defer emailQuery.Close()
	if emailQuery.Next() {
		log.Println("email not unique")
		return fiber.ErrBadRequest
	}

	// Validate role is correct
	if params.Role != "normal" && params.Role != "admin" {
		log.Println("role not valid")
		return fiber.ErrBadRequest
	}

	// Validate theme is correct
	if params.Theme != "system" && params.Theme != "light" && params.Theme != "dark" {
		log.Println("theme not valid")
		return fiber.ErrBadRequest
	}

	// Create user
	userID := uuid.NewV4()

	var userRole models.UserRole
	if params.Role == "normal" {
		userRole = models.UserRoleNormal
	}
	if params.Role == "admin" {
		userRole = models.UserRoleAdmin
	}

	var userTheme models.UserTheme
	if params.Theme == "system" {
		userTheme = models.UserThemeSystem
	}
	if params.Theme == "light" {
		userTheme = models.UserThemeLight
	}
	if params.Theme == "dark" {
		userTheme = models.UserThemeDark
	}

	database.Exec("INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`, `theme`) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		userID.String(), params.Username, params.Email, utils.HashPassword(params.Password), userRole, userTheme)

	// Get new created user and send response
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", userID.String())
	defer userQuery.Close()
	userQuery.Next()
	return c.JSON(models.UserScan(c, userQuery))
}

func UsersShow(c *fiber.Ctx) error {
	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.UserScan(c, userQuery))
}

type UsersEditParams struct {
	Username string `form:"username" validate:"required,min=2"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"omitempty,min=6"`
	Role     string `form:"role" validate:"omitempty,required"`
	Theme    string `form:"theme" validate:"required"`
}

func UsersEdit(c *fiber.Ctx) error {
	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get user
	user := models.UserScan(c, userQuery)

	// Parse body
	var params UsersEditParams
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate values
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// Validate username is unique when diffrent
	if user.Username != params.Username {
		usernameQuery := database.Query("SELECT `id` FROM `users` WHERE `username` = ?", params.Username)
		defer usernameQuery.Close()
		if usernameQuery.Next() {
			log.Println("username not unique")
			return fiber.ErrBadRequest
		}
	}

	// Validate email is unique
	if user.Email != params.Email {
		emailQuery := database.Query("SELECT `id` FROM `users` WHERE `email` = ?", params.Email)
		defer emailQuery.Close()
		if emailQuery.Next() {
			log.Println("email not unique")
			return fiber.ErrBadRequest
		}
	}

	// Validate role is correct
	if params.Role != "" && params.Role != "normal" && params.Role != "admin" {
		log.Println("role not valid")
		return fiber.ErrBadRequest
	}

	// Validate theme is correct
	if params.Theme != "system" && params.Theme != "light" && params.Theme != "dark" {
		log.Println("theme not valid")
		return fiber.ErrBadRequest
	}

	// Update user
	var userRole models.UserRole
	if params.Role == "normal" {
		userRole = models.UserRoleNormal
	}
	if params.Role == "admin" {
		userRole = models.UserRoleAdmin
	}

	var userTheme models.UserTheme
	if params.Theme == "system" {
		userTheme = models.UserThemeSystem
	}
	if params.Theme == "light" {
		userTheme = models.UserThemeLight
	}
	if params.Theme == "dark" {
		userTheme = models.UserThemeDark
	}

	if params.Role != "" {
		if params.Password != "" {
			database.Exec("UPDATE `users` SET `username` = ?, `email` = ?, `password` = ?, `role` = ?, theme = ? WHERE `id` = UUID_TO_BIN(?)", params.Username, params.Email, utils.HashPassword(params.Password), userRole, userTheme, user.ID)
		} else {
			database.Exec("UPDATE `users` SET `username` = ?, `email` = ?, `role` = ?, `theme` = ? WHERE `id` = UUID_TO_BIN(?)", params.Username, params.Email, userRole, userTheme, user.ID)
		}
	} else {
		database.Exec("UPDATE `users` SET `username` = ?, `email` = ?, `theme` = ? WHERE `id` = UUID_TO_BIN(?)", params.Username, params.Email, userTheme, user.ID)
	}

	// Get edited user and send response
	updatedUserQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", user.ID)
	defer updatedUserQuery.Close()
	updatedUserQuery.Next()
	return c.JSON(models.UserScan(c, updatedUserQuery))
}

func UsersLikedArtists(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get total liked artists
	total := database.Count("SELECT COUNT(`artists`.`id`) FROM `artists` INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id` "+
		"WHERE `artists`.`name` LIKE ?", "%"+query+"%")

	// Get liked artists
	artistsQuery := database.Query("SELECT BIN_TO_UUID(`artists`.`id`), `artists`.`name`, `artists`.`deezer_id`, `artists`.`created_at` FROM `artists` INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id` "+
		"WHERE `artists`.`name` LIKE ? ORDER BY LOWER(`artists`.`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer artistsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.ArtistsScan(c, artistsQuery, false, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersLikedAlbums(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get total liked albums
	total := database.Count("SELECT COUNT(`albums`.`id`) FROM `albums` INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id` "+
		"WHERE `albums`.`title` LIKE ?", "%"+query+"%")

	// Get liked albums
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`albums`.`id`), `albums`.`type`, `albums`.`title`, `albums`.`released_at`, `albums`.`explicit`, `albums`.`deezer_id`, `albums`.`created_at` FROM `albums` "+
		"INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id` WHERE `albums`.`title` LIKE ? ORDER BY LOWER(`albums`.`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer albumsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.AlbumsScan(c, albumsQuery, true, true, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersLikedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get total liked tracks
	total := database.Count("SELECT COUNT(`tracks`.`id`) FROM `tracks` INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id` "+
		"WHERE `tracks`.`title` LIKE ?", "%"+query+"%")

	// Get liked tracks
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`tracks`.`id`), BIN_TO_UUID(`tracks`.`album_id`), `tracks`.`title`, `tracks`.`disk`, `tracks`.`position`, `tracks`.`duration`, `tracks`.`explicit`, `tracks`.`deezer_id`, `tracks`.`youtube_id`, `tracks`.`plays`, `tracks`.`created_at` FROM `tracks` "+
		"INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id` WHERE `tracks`.`title` LIKE ? ORDER BY `tracks`.`plays` DESC, LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer tracksQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.TracksScan(c, tracksQuery, true, true),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersPlayedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

	// Check if user exists
	userQuery := database.Query("SELECT `id` FROM `users` WHERE `id` = UUID_TO_BIN(?)", c.Params("userID"))
	defer userQuery.Close()
	if !userQuery.Next() {
		return fiber.ErrNotFound
	}

	// Get total played tracks
	total := database.Count("SELECT COUNT(`tracks`.`id`) FROM `tracks` INNER JOIN `track_plays` ON `tracks`.`id` = `track_plays`.`track_id` "+
		"WHERE `tracks`.`title` LIKE ?", "%"+query+"%")

	// Get played tracks
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`tracks`.`id`), BIN_TO_UUID(`tracks`.`album_id`), `tracks`.`title`, `tracks`.`disk`, `tracks`.`position`, `tracks`.`duration`, `tracks`.`explicit`, `tracks`.`deezer_id`, `tracks`.`youtube_id`, `tracks`.`plays`, `tracks`.`created_at` FROM `tracks` "+
		"INNER JOIN `track_plays` ON `tracks`.`id` = `track_plays`.`track_id` WHERE `tracks`.`title` LIKE ? ORDER BY `track_plays`.`created_at` DESC LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer tracksQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.TracksScan(c, tracksQuery, true, true),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func UsersSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != c.Params("userID") {
		return fiber.ErrUnauthorized
	}

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
		"data": models.SessionsScan(c, sessionsQuery, false),
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
	return c.JSON(fiber.Map{"success": true})
}
