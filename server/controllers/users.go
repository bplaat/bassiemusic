package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/bplaat/bassiemusic/utils/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.UserModel().WhereRaw("`username` LIKE ?", "%"+query+"%").WhereOrRaw("`email` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`username`)").Paginate(page, limit))
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
	if models.UserModel().Where("username", params.Username).First() != nil {
		log.Println("username not unique")
		return fiber.ErrBadRequest
	}

	// Validate email is unique
	if models.UserModel().Where("email", params.Email).First() != nil {
		log.Println("email not unique")
		return fiber.ErrBadRequest
	}

	// Validate role is correct
	if params.Role != "normal" && params.Role != "admin" {
		log.Println("role not valid")
		return fiber.ErrBadRequest
	}
	var userRole models.UserRole
	if params.Role == "normal" {
		userRole = models.UserRoleNormal
	}
	if params.Role == "admin" {
		userRole = models.UserRoleAdmin
	}

	// Validate theme is correct
	if params.Theme != "system" && params.Theme != "light" && params.Theme != "dark" {
		log.Println("theme not valid")
		return fiber.ErrBadRequest
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

	// Create user
	return c.JSON(models.UserModel().Create(database.Map{
		"username": params.Username,
		"email":    params.Email,
		"password": utils.HashPassword(params.Password),
		"role":     userRole,
		"theme":    userTheme,
	}))
}

func UsersShow(c *fiber.Ctx) error {
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

func UsersLikedArtists(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked artists
	likedArtists := models.ArtistModel(c).Join("INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id`").
		WhereRaw("`artist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).WhereRaw("`artists`.`name` LIKE ?", "%"+query+"%").
		OrderByRaw("`artist_likes`.`created_at` DESC").Paginate(page, limit)
	return c.JSON(likedArtists)
}

func UsersLikedAlbums(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked albums
	likedAlbums := models.AlbumModel(c).Join("INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id`").
		With("artists", "genres").WhereRaw("`album_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`albums`.`title` LIKE ?", "%"+query+"%").OrderByRaw("`album_likes`.`created_at` DESC").Paginate(page, limit)
	return c.JSON(likedAlbums)
}

func UsersLikedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked tracks
	likedTracks := models.TrackModel(c).Join("INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id`").
		With("like", "artists", "album").WhereRaw("`track_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%").OrderByRaw("`track_likes`.`created_at` DESC").Paginate(page, limit)
	return c.JSON(likedTracks)
}

func UsersLikedPlaylists(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked playlists
	likedPlaylists := models.PlaylistModel(c).Join("INNER JOIN `playlist_likes` ON `playlists`.`id` = `playlist_likes`.`playlist_id`").
		With("like").WhereRaw("`playlist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`playlists`.`name` LIKE ?", "%"+query+"%").OrderByRaw("`playlist_likes`.`created_at` DESC").Paginate(page, limit)
	return c.JSON(likedPlaylists)
}

func UsersPlayedTracks(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get played tracks
	playedTracks := models.TrackModel(c).Join("INNER JOIN `track_plays` ON `tracks`.`id` = `track_plays`.`track_id`").
		With("like", "artists", "album").WhereRaw("`track_plays`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%").OrderByRaw("`track_plays`.`updated_at` DESC").Paginate(page, limit)
	return c.JSON(playedTracks)
}

func UsersSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user sessions
	userSessions := models.SessionModel().Where("user_id", user.ID).OrderByDesc("created_at").Paginate(page, limit)
	return c.JSON(userSessions)
}

func UsersPlaylists(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user playlists
	userPlaylists := models.PlaylistModel(c).With("like").Where("user_id", user.ID).OrderByRaw("LOWER(`name`)").Paginate(page, limit)
	return c.JSON(userPlaylists)
}

type UsersEditParams struct {
	Username      string `form:"username" validate:"required,min=2"`
	Email         string `form:"email" validate:"required,email"`
	Password      string `form:"password" validate:"omitempty,min=6"`
	AllowExplicit string `form:"allow_explicit" validate:"omitempty,required"`
	Role          string `form:"role" validate:"omitempty,required"`
	Theme         string `form:"theme" validate:"required"`
}

func UsersEdit(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
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
	if user.Username != params.Username && models.UserModel().Where("username", params.Username).First() != nil {
		log.Println("username not unique")
		return fiber.ErrBadRequest
	}

	// Validate email is unique
	if user.Email != params.Email && models.UserModel().Where("email", params.Email).First() != nil {
		log.Println("email not unique")
		return fiber.ErrBadRequest
	}

	// Validate allow_explicit is correct
	if params.AllowExplicit != "" && params.AllowExplicit != "true" && params.AllowExplicit != "false" {
		log.Println("allow_explicit not valid")
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

	updates := database.Map{
		"username":       params.Username,
		"email":          params.Email,
		"allow_explicit": params.AllowExplicit == "true",
		"theme":          userTheme,
	}
	if params.Password != "" {
		updates["password"] = utils.HashPassword(params.Password)
	}
	if authUser.Role == "admin" && params.Role != "" {
		updates["role"] = userRole
	}
	models.UserModel().Where("id", user.ID).Update(updates)

	// Get updated user
	return c.JSON(models.UserModel().Find(user.ID))
}

func UsersAvatar(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Remove old avatar file
	if user.AvatarID != nil && *user.AvatarID != "" {
		if err := os.Remove(fmt.Sprintf("storage/avatars/%s.jpg", *user.AvatarID)); err != nil {
			log.Fatalln(err)
		}
	}

	// Save uploaded avatar file
	avatarID := uuid.New()
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = c.SaveFile(avatar, fmt.Sprintf("storage/avatars/%s.jpg", avatarID.String())); err != nil {
		log.Fatalln(err)
	}

	// Save avatar id for user
	models.UserModel().Where("id", user.ID).Update(database.Map{
		"avatar": avatarID.String(),
	})
	return c.JSON(fiber.Map{"success": true})
}

func UsersAvatarDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Check if user has avatar
	if user.AvatarID != nil && *user.AvatarID != "" {
		// Remove old avatar file
		if err := os.Remove(fmt.Sprintf("storage/avatars/%s.jpg", *user.AvatarID)); err != nil {
			log.Fatalln(err)
		}

		// Clear avatar id for user
		models.UserModel().Where("id", user.ID).Update(database.Map{
			"avatar": nil,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

func UsersDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Delete user
	models.UserModel().Where("id", user.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
