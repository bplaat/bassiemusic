package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func UsersIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.UserModel.WhereRaw("`username` LIKE ?", "%"+query+"%").WhereOrRaw("`email` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`username`)").Paginate(page, limit))
}

type UsersCreateBody struct {
	Username      string `form:"username" validate:"required|min:2|unique:users,username"`
	Email         string `form:"email" validate:"required|email|unique:users,email"`
	Password      string `form:"password" validate:"required|min:6"`
	AllowExplicit string `form:"allow_explicit" validate:"required|boolean"`
	Role          string `form:"role" validate:"required|enum:normal,admin"`
}

func UsersCreate(c *fiber.Ctx) error {
	// Parse body
	var body UsersCreateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create user
	var userRole models.UserRole
	if body.Role == "normal" {
		userRole = models.UserRoleNormal
	}
	if body.Role == "admin" {
		userRole = models.UserRoleAdmin
	}
	return c.JSON(models.UserModel.Create(database.Map{
		"username":       body.Username,
		"email":          body.Email,
		"password":       utils.HashPassword(body.Password),
		"allow_explicit": body.AllowExplicit == "true",
		"role":           userRole,
		"language":       "en",
		"theme":          models.UserThemeSystem,
	}))
}

func UsersShow(c *fiber.Ctx) error {
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

type UsersUpdateBody struct {
	Username      *string `form:"username" validate:"min:2|unique:users,username,Username"`
	Email         *string `form:"email" validate:"email|unique:users,email,Email"`
	Password      *string `form:"password" validate:"min:6"`
	AllowExplicit *string `form:"allow_explicit" validate:"boolean"`
	Role          *string `form:"role" validate:"enum:normal,admin"`
	Language      *string `form:"language" validate:"enum:en,nl"`
	Theme         *string `form:"theme" validate:"enum:system,light,dark"`
}

func UsersUpdate(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var body UsersUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, user, &body); err != nil {
		return nil
	}

	// Run updates
	updates := database.Map{}
	if body.Username != nil {
		updates["username"] = *body.Username
	}
	if body.Email != nil {
		updates["email"] = *body.Email
	}
	if body.Password != nil {
		updates["password"] = utils.HashPassword(*body.Password)
	}
	if body.AllowExplicit != nil {
		updates["allow_explicit"] = *body.AllowExplicit == "true"
	}
	if body.Role != nil {
		if *body.Role == "normal" {
			updates["role"] = models.UserRoleNormal
		}
		if *body.Role == "admin" {
			updates["role"] = models.UserRoleAdmin
		}
	}
	if body.Language != nil {
		updates["language"] = *body.Language
	}
	if body.Theme != nil {
		if *body.Theme == "system" {
			updates["theme"] = models.UserThemeSystem
		}
		if *body.Theme == "light" {
			updates["theme"] = models.UserThemeLight
		}
		if *body.Theme == "dark" {
			updates["theme"] = models.UserThemeDark
		}
	}
	models.UserModel.Where("id", user.ID).Update(updates)

	// Get updated user
	return c.JSON(models.UserModel.Find(user.ID))
}

func UsersDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Delete user
	models.UserModel.Where("id", user.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func UsersAvatar(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
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
	models.UserModel.Where("id", user.ID).Update(database.Map{
		"avatar": avatarID.String(),
	})
	return c.JSON(fiber.Map{"success": true})
}

func UsersAvatarDelete(c *fiber.Ctx) error {
	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Check if user has avatar
	if user.AvatarID != nil && *user.AvatarID != "" {
		// Remove old avatar file
		_ = os.Remove(fmt.Sprintf("storage/avatars/original/%s", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/small/%s.jpg", *user.AvatarID))
		_ = os.Remove(fmt.Sprintf("storage/avatars/medium/%s.jpg", *user.AvatarID))

		// Clear avatar id for user
		models.UserModel.Where("id", user.ID).Update(database.Map{
			"avatar": nil,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

func UsersLikedArtists(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked artists
	q := models.ArtistModel.Join("INNER JOIN `artist_likes` ON `artists`.`id` = `artist_likes`.`artist_id`").
		WhereRaw("`artist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).WhereRaw("`artists`.`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "name" {
		q = q.OrderByRaw("LOWER(`artists`.`name`)")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`artists`.`name`) DESC")
	} else if c.Query("sort_by") == "sync" {
		q = q.OrderByRaw("`artists`.`sync` DESC, LOWER(`artists`.`name`)")
	} else if c.Query("sort_by") == "sync_desc" {
		q = q.OrderByRaw("`artists`.`sync`, LOWER(`artists`.`name`)")
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

func UsersLikedGenres(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked genres
	q := models.GenreModel.Join("INNER JOIN `genre_likes` ON `genres`.`id` = `genre_likes`.`genre_id`").
		WhereRaw("`genre_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).WhereRaw("`genres`.`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "name" {
		q = q.OrderByRaw("LOWER(`genres`.`name`)")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`genres`.`name`) DESC")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderByRaw("`genres`.`created_at`")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByRaw("`genres`.`created_at` DESC")
	} else if c.Query("sort_by") == "liked_at" {
		q = q.OrderByRaw("`genre_likes`.`created_at`")
	} else {
		q = q.OrderByRaw("`genre_likes`.`created_at` DESC")
	}
	return c.JSON(q.Paginate(page, limit))
}

func UsersLikedAlbums(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked albums
	q := models.AlbumModel.Join("INNER JOIN `album_likes` ON `albums`.`id` = `album_likes`.`album_id`").
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
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked tracks
	q := models.TrackModel.Join("INNER JOIN `track_likes` ON `tracks`.`id` = `track_likes`.`track_id`").
		With("liked_true", "artists", "album").WhereRaw("`track_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
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
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get liked playlists
	q := models.PlaylistModel.Join("INNER JOIN `playlist_likes` ON `playlists`.`id` = `playlist_likes`.`playlist_id`").
		WhereRaw("`playlist_likes`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`playlists`.`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "name" {
		q = q.OrderByRaw("LOWER(`playlists`.`name`)")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`playlists`.`name`) DESC")
	} else if c.Query("sort_by") == "public" {
		q = q.OrderByRaw("`playlists`.`public` DESC, LOWER(`playlists`.`name`)")
	} else if c.Query("sort_by") == "public_desc" {
		q = q.OrderByRaw("`playlists`.`public`, LOWER(`playlists`.`name`)")
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
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get played tracks
	playedTracks := models.TrackModel.Join("INNER JOIN `track_plays` ON `tracks`.`id` = `track_plays`.`track_id`").
		WithArgs("liked", c.Locals("authUser")).With("artists", "album").WhereRaw("`track_plays`.`user_id` = UUID_TO_BIN(?)", authUser.ID).
		WhereRaw("`tracks`.`title` LIKE ?", "%"+query+"%").OrderByRaw("`track_plays`.`updated_at` DESC").Paginate(page, limit)
	return c.JSON(playedTracks)
}

func UsersSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user sessions
	userSessions := models.SessionModel.Where("user_id", user.ID).OrderByDesc("created_at").Paginate(page, limit)
	return c.JSON(userSessions)
}

func UsersActiveSessions(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user sessions
	userSessions := models.SessionModel.Where("user_id", user.ID).WhereRaw("`expires_at` > ?", time.Now()).
		OrderByDesc("created_at").Paginate(page, limit)
	return c.JSON(userSessions)
}

func UsersPlaylists(c *fiber.Ctx) error {
	_, page, limit := utils.ParseIndexVars(c)

	// Check if user exists
	user := models.UserModel.Find(c.Params("userID"))
	if user == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != models.UserRoleAdmin && authUser.ID != user.ID {
		return fiber.ErrUnauthorized
	}

	// Get user playlists
	q := models.PlaylistModel.WithArgs("liked", c.Locals("authUser")).Where("user_id", user.ID)
	if c.Query("sort_by") == "public" {
		q = q.OrderByRaw("`public` DESC, LOWER(`name`)")
	} else if c.Query("sort_by") == "public_desc" {
		q = q.OrderByRaw("`public`, LOWER(`name`)")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderBy("created_at")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByDesc("created_at")
	} else if c.Query("sort_by") == "updated_at" {
		q = q.OrderBy("updated_at")
	} else if c.Query("sort_by") == "updated_at_desc" {
		q = q.OrderByDesc("updated_at")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`name`) DESC")
	} else {
		q = q.OrderByRaw("LOWER(`name`)")
	}
	return c.JSON(q.Paginate(page, limit))
}
