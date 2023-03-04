package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/bplaat/bassiemusic/utils/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.UserModel().WhereRaw("`username` LIKE ?", "%"+query+"%").WhereOrRaw("`email` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`username`)").Paginate(page, limit))
}

type UsersCreateParams struct {
	Username      string `form:"username" validate:"required,min=2"`
	Email         string `form:"email" validate:"required,email"`
	Password      string `form:"password" validate:"required,min=6"`
	AllowExplicit string `form:"allow_explicit" validate:"required"`
	Role          string `form:"role" validate:"required"`
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

	// Validate allow_explicit is correct
	if params.AllowExplicit != "true" && params.AllowExplicit != "false" {
		log.Println("allow_explicit not valid")
		return fiber.ErrBadRequest
	}

	// Create user
	return c.JSON(models.UserModel().Create(database.Map{
		"username":       params.Username,
		"email":          params.Email,
		"password":       utils.HashPassword(params.Password),
		"allow_explicit": params.AllowExplicit == "true",
		"role":           userRole,
		"language":       "en",
		"theme":          models.UserThemeSystem,
	}))
}

func UsersShow(c *fiber.Ctx) error {
	user := models.UserModel().Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

type UsersEditParams struct {
	Username      string `form:"username" validate:"required,min=2"`
	Email         string `form:"email" validate:"required,email"`
	Password      string `form:"password" validate:"omitempty,min=6"`
	AllowExplicit string `form:"allow_explicit" validate:"omitempty,required"`
	Role          string `form:"role" validate:"omitempty,required"`
	Language      string `form:"language" validate:"required"`
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

	// Validate lang is correct
	if params.Language != "en" && params.Language != "nl" {
		log.Println("lang not valid")
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
		"language":       params.Language,
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
		_ = os.Remove(fmt.Sprintf("storage/avatars/original/%s", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/small/%s.jpg", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/medium/%s.jpg", *user.AvatarID))
	}

	// Save uploaded avatar file
	avatarID := uuid.New()
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = c.SaveFile(avatar, fmt.Sprintf("storage/avatars/original/%s", avatarID.String())); err != nil {
		log.Fatalln(err)
	}

	// Open uploaded image
	originalFile, err := os.Open(fmt.Sprintf("storage/avatars/original/%s", avatarID.String()))
	if err != nil {
		log.Fatalln(err)
	}
	defer originalFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var originalImage image.Image
	originalImage, _, err = image.Decode(originalFile)
	if err != nil {
		if err := os.Remove(fmt.Sprintf("storage/avatars/original/%s", avatarID.String())); err != nil {
			log.Fatalln(err)
		}
		return c.JSON(fiber.Map{"success": false})
	}

	// Save small resize
	smallFile, err := os.Create(fmt.Sprintf("storage/avatars/small/%s.jpg", avatarID.String()))
	if err != nil {
		log.Fatalln(err)
	}
	defer smallFile.Close()
	smallImage := resize.Resize(250, 250, originalImage, resize.Lanczos3)
	err = jpeg.Encode(smallFile, smallImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save small resize
	mediumFile, err := os.Create(fmt.Sprintf("storage/avatars/medium/%s.jpg", avatarID.String()))
	if err != nil {
		log.Fatalln(err)
	}
	defer mediumFile.Close()
	mediumImage := resize.Resize(500, 500, originalImage, resize.Lanczos3)
	err = jpeg.Encode(mediumFile, mediumImage, nil)
	if err != nil {
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
		_ = os.Remove(fmt.Sprintf("storage/avatars/original/%s", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/small/%s.jpg", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/medium/%s.jpg", *user.AvatarID))

		// Clear avatar id for user
		models.UserModel().Where("id", user.ID).Update(database.Map{
			"avatar": nil,
		})
	}
	return c.JSON(fiber.Map{"success": true})
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
	q := models.ArtistModel(c).Join("INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id`").
		WhereRaw("`artist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).WhereRaw("`artists`.`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "name" {
		q = q.OrderByRaw("LOWER(`artists`.`name`)")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`artists`.`name`) DESC")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderByRaw("`artists`.`created_at`")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByRaw("`artists`.`created_at` DESC")
	} else if c.Query("sort_by") == "liked_at" {
		q = q.OrderByRaw("`artist_likes`.`created_at`")
	} else {
		q = q.OrderByRaw("`artist_likes`.`created_at` DESC")
	}
	return c.JSON(q.Paginate(page, limit))
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
	q := models.AlbumModel(c).Join("INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id`").
		With("artists", "genres").WhereRaw("`album_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`albums`.`title` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "title" {
		q = q.OrderByRaw("LOWER(`albums`.`title`)")
	} else if c.Query("sort_by") == "title_desc" {
		q = q.OrderByRaw("LOWER(`albums`.`title`) DESC")
	} else if c.Query("sort_by") == "released_at" {
		q = q.OrderByRaw("`albums`.`released_at`")
	} else if c.Query("sort_by") == "released_at_desc" {
		q = q.OrderByRaw("`albums`.`released_at` DESC")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderByRaw("`albums`.`created_at`")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByRaw("`albums`.`created_at` DESC")
	} else if c.Query("sort_by") == "liked_at" {
		q = q.OrderByRaw("`album_likes`.`created_at`")
	} else {
		q = q.OrderByRaw("`album_likes`.`created_at` DESC")
	}
	return c.JSON(q.Paginate(page, limit))
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
	q := models.TrackModel(c).Join("INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id`").
		With("like", "artists", "album").WhereRaw("`track_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "title" {
		q = q.OrderByRaw("LOWER(`tracks`.`title`)")
	} else if c.Query("sort_by") == "title_desc" {
		q = q.OrderByRaw("LOWER(`tracks`.`title`) DESC")
	} else if c.Query("sort_by") == "plays" {
		q = q.OrderByRaw("`tracks`.`plays`, LOWER(`tracks`.`title`)")
	} else if c.Query("sort_by") == "plays_desc" {
		q = q.OrderByRaw("`tracks`.`plays` DESC, LOWER(`tracks`.`title`)")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderByRaw("`tracks`.`created_at`")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByRaw("`tracks`.`created_at` DESC")
	} else if c.Query("sort_by") == "liked_at" {
		q = q.OrderByRaw("`track_likes`.`created_at`")
	} else {
		q = q.OrderByRaw("`track_likes`.`created_at` DESC")
	}
	return c.JSON(q.Paginate(page, limit))
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
	q := models.PlaylistModel(c).Join("INNER JOIN `playlist_likes` ON `playlists`.`id` = `playlist_likes`.`playlist_id`").
		With("like").WhereRaw("`playlist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`playlists`.`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "name" {
		q = q.OrderByRaw("LOWER(`playlists`.`name`)")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`playlists`.`name`) DESC")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderByRaw("`playlists`.`created_at`")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByRaw("`playlists`.`created_at` DESC")
	} else if c.Query("sort_by") == "updated_at" {
		q = q.OrderByRaw("`playlists`.`updated_at`")
	} else if c.Query("sort_by") == "updated_at_desc" {
		q = q.OrderByRaw("`playlists`.`updated_at` DESC")
	} else if c.Query("sort_by") == "liked_at" {
		q = q.OrderByRaw("`playlist_likes`.`created_at`")
	} else {
		q = q.OrderByRaw("`playlist_likes`.`created_at` DESC")
	}
	return c.JSON(q.Paginate(page, limit))
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

func UsersActiveSessions(c *fiber.Ctx) error {
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
	userSessions := models.SessionModel().Where("user_id", user.ID).WhereRaw("`expires_at` > ?", time.Now()).
		OrderByDesc("created_at").Paginate(page, limit)
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
