package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.UserModel(c).WhereRaw("`username` LIKE ?", "%"+query+"%").WhereOrRaw("`email` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`username`)").Paginate(page, limit))
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
	if models.UserModel(c).Where("username", params.Username).First() != nil {
		log.Println("username not unique")
		return fiber.ErrBadRequest
	}

	// Validate email is unique
	if models.UserModel(c).Where("email", params.Email).First() != nil {
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
	user := models.User{
		Username: params.Username,
		Email:    params.Email,
		Password: utils.HashPassword(params.Password),
		Role:     params.Role,
	}
	return c.JSON(models.UserModel(c).Create(&user))
}

func UsersShow(c *fiber.Ctx) error {
	user := models.UserModel(c).With("albums", "top_tracks").Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

type UsersEditParams struct {
	Username string `form:"username" validate:"required,min=2"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"omitempty,min=6"`
	Role     string `form:"role" validate:"omitempty,required"`
	Theme    string `form:"theme" validate:"required"`
}

func UsersEdit(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

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
	if user.Username != params.Username && models.UserModel(c).Where("username", params.Username).First() != nil {
		log.Println("username not unique")
		return fiber.ErrBadRequest
	}

	// Validate email is unique
	if user.Email != params.Email && models.UserModel(c).Where("email", params.Email).First() != nil {
		log.Println("email not unique")
		return fiber.ErrBadRequest
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
	updatedUserQuery := database.Query("SELECT BIN_TO_UUID(`id`), `username`, `email`, `password`, BIN_TO_UUID(`avatar`), `role`, `theme`, `created_at` FROM `users` WHERE `id` = UUID_TO_BIN(?)", user.ID)
	defer updatedUserQuery.Close()
	updatedUserQuery.Next()
	return c.JSON(models.UserScan(c, updatedUserQuery))
}

func UsersAvatar(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Remove old avatar file
	if user.AvatarID != "" {
		if err := os.Remove(fmt.Sprintf("storage/avatars/%s.jpg", user.AvatarID)); err != nil {
			log.Fatalln(err)
		}
	}

	// Save uploaded avatar file
	avatarID := uuid.NewV4()
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = c.SaveFile(avatar, fmt.Sprintf("storage/avatars/%s.jpg", avatarID.String())); err != nil {
		log.Fatalln(err)
	}

	// Save avatar id for user
	user.AvatarID = avatarID.String()
	models.UserModel(c).Update(user)
	return c.JSON(fiber.Map{"success": true})
}

func UsersAvatarDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Check if user has avatar
	if user.AvatarID != "" {
		// Remove old avatar file
		if err := os.Remove(fmt.Sprintf("storage/avatars/%s.jpg", user.AvatarID)); err != nil {
			log.Fatalln(err)
		}

		// Clear avatar id for user
		database.Exec("UPDATE `users` SET `avatar` = NULL WHERE `id` = UUID_TO_BIN(?)", user.ID)
	}
	return c.JSON(fiber.Map{"success": true})
}

func UsersLikedArtists(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked artists
	likedArtists := models.ArtistModel(c).Join("INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id`").
		WhereRaw("`artists`.`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`artists`.`name`)").Paginate(page, limit)
	return c.JSON(likedArtists)
}

func UsersLikedAlbums(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked albums
	likedAlbums := models.AlbumModel(c).Join("INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id`").
		With("artists", "genres").WhereRaw("`albums`.`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`albums`.`title`)").Paginate(page, limit)
	return c.JSON(likedAlbums)
}

func UsersLikedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked tracks
	likedTracks := models.TrackModel(c).Join("INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id`").
		With("artists", "album").WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%").OrderByDesc("plays").Paginate(page, limit)
	return c.JSON(likedTracks)
}

func UsersPlayedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get played tracks
	playedTracks := models.TrackModel(c).Join("INNER JOIN `track_plays` ON `tracks`.`id` = `track_plays`.`track_id`").
		With("artists", "album").WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%").OrderByDesc("plays").Paginate(page, limit)
	return c.JSON(playedTracks)
}

func UsersSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := models.AuthUser(c)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user sessions
	userSessions := models.SessionModel(c).Where("user_id", user.ID).OrderByDesc("created_at").Paginate(page, limit)
	return c.JSON(userSessions)
}

func UsersDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel(c).Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Delete user
	models.UserModel(c).Where("id", user.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
