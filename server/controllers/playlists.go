package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/bplaat/bassiemusic/utils/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func PlaylistsIndex(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)
	query, page, limit := utils.ParseIndexVars(c)
	q := models.PlaylistModel(c).With("liked", "user").WhereRaw("`name` LIKE ?", "%"+query+"%")
	if authUser.Role != models.UserRoleAdmin {
		q = q.Where("public", true)
	}
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

type PlaylistsCreateParams struct {
	Name   string `form:"name" validate:"required,min=2"`
	Public string `form:"public" validate:"required"`
}

func PlaylistsCreate(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse body
	var params PlaylistsCreateParams
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

	// Validate public is correct
	if params.Public != "true" && params.Public != "false" {
		log.Println("public not valid")
		return fiber.ErrBadRequest
	}

	// Create playlist
	return c.JSON(models.PlaylistModel(c).Create(database.Map{
		"user_id": authUser.ID,
		"name":    params.Name,
		"public":  params.Public == "true",
	}))
}

func PlaylistsShow(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).With("liked", "user", "tracks").Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(playlist.Public || authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	return c.JSON(playlist)
}

type PlaylistsUpdateParams struct {
	Name   string `form:"name" validate:"required,min=2"`
	Public string `form:"public" validate:"required"`
}

func PlaylistsUpdate(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var params PlaylistsUpdateParams
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

	// Validate public is correct
	if params.Public != "true" && params.Public != "false" {
		log.Println("public not valid")
		return fiber.ErrBadRequest
	}

	// Update playlist
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name":   params.Name,
		"public": params.Public == "true",
	})

	// Get updated playlist
	return c.JSON(models.PlaylistModel(c).Find(playlist.ID))
}

func PlaylistsDelete(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Delete playlist
	models.PlaylistModel(c).Where("id", playlist.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsImage(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Remove old image file
	if playlist.ImageID != nil && *playlist.ImageID != "" {
		_ = os.Remove(fmt.Sprintf("storage/playlists/original/%s", *playlist.ImageID))
		_ = os.Remove(fmt.Sprintf("storage/playlists/small/%s.jpg", *playlist.ImageID))
		_ = os.Remove(fmt.Sprintf("storage/playlists/medium/%s.jpg", *playlist.ImageID))
	}

	// Save uploaded image file
	imageID := uuid.New()
	uploadedImage, err := c.FormFile("image")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = c.SaveFile(uploadedImage, fmt.Sprintf("storage/playlists/original/%s", imageID.String())); err != nil {
		log.Fatalln(err)
	}

	// Open uploaded image
	originalFile, err := os.Open(fmt.Sprintf("storage/playlists/original/%s", imageID.String()))
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
		if err := os.Remove(fmt.Sprintf("storage/playlists/original/%s", imageID.String())); err != nil {
			log.Fatalln(err)
		}
		return c.JSON(fiber.Map{"success": false})
	}

	// Save small resize
	smallFile, err := os.Create(fmt.Sprintf("storage/playlists/small/%s.jpg", imageID.String()))
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
	mediumFile, err := os.Create(fmt.Sprintf("storage/playlists/medium/%s.jpg", imageID.String()))
	if err != nil {
		log.Fatalln(err)
	}
	defer mediumFile.Close()
	mediumImage := resize.Resize(500, 500, originalImage, resize.Lanczos3)
	err = jpeg.Encode(mediumFile, mediumImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save image id for playlist
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"image": imageID.String(),
	})
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsImageDelete(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Check if playlist has image
	if playlist.ImageID != nil && *playlist.ImageID != "" {
		// Remove old image file
		_ = os.Remove(fmt.Sprintf("storage/playlists/original/%s", *playlist.ImageID))
		_ = os.Remove(fmt.Sprintf("storage/playlists/small/%s.jpg", *playlist.ImageID))
		_ = os.Remove(fmt.Sprintf("storage/playlists/medium/%s.jpg", *playlist.ImageID))

		// Clear image id for playlist
		models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
			"image": nil,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsAppendTrackParams struct {
	TrackID string `form:"track_id" validate:"required"`
}

func PlaylistsAppendTrack(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var params PlaylistsInsertTrackParams
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

	// Validate track id exists
	if models.TrackModel(c).Find(params.TrackID) == nil {
		log.Println("track_id not valid")
		return fiber.ErrBadRequest
	}

	// Trigger playlist update
	oldName := playlist.Name
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Create playlist track
	models.PlaylistTrackModel().Create(database.Map{
		"playlist_id": playlist.ID,
		"track_id":    params.TrackID,
		"position":    models.PlaylistTrackModel().Where("playlist_id", playlist.ID).Count() + 1,
	})
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsInsertTrackParams struct {
	TrackID string `form:"track_id" validate:"required"`
}

func PlaylistsInsertTrack(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var params PlaylistsInsertTrackParams
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

	// Validate track id exists
	if models.TrackModel(c).Find(params.TrackID) == nil {
		log.Println("track_id not valid")
		return fiber.ErrBadRequest
	}

	// Parse position
	var position int64
	if positionInt, err := strconv.ParseInt(c.Params("position"), 10, 64); err == nil {
		position = positionInt
		if position < 0 {
			position = 0
		}
	}

	// Trigger playlist update
	oldName := playlist.Name
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Move all existing track ids to the left
	playlistTracks := models.PlaylistTrackModel().Where("playlist_id", playlist.ID).WhereRaw("`position` >= ?", position).Get()
	for _, playlistTrack := range playlistTracks {
		models.PlaylistTrackModel().Where("id", playlistTrack.ID).Update(database.Map{
			"position": playlistTrack.Position + 1,
		})
	}

	// Create playlist track
	models.PlaylistTrackModel().Create(database.Map{
		"playlist_id": playlist.ID,
		"track_id":    params.TrackID,
		"position":    position,
	})
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsRemoveTrack(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse position
	var position int64
	if positionInt, err := strconv.ParseInt(c.Params("position"), 10, 64); err == nil {
		position = positionInt
		if position < 0 {
			position = 0
		}
	}

	// Trigger playlist update
	oldName := playlist.Name
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Remove playlist track
	models.PlaylistTrackModel().Where("playlist_id", playlist.ID).Where("position", position).Delete()

	// Move all existing track ids to the right
	playlistTracks := models.PlaylistTrackModel().Where("playlist_id", playlist.ID).WhereRaw("`position` > ?", position).Get()
	for _, playlistTrack := range playlistTracks {
		models.PlaylistTrackModel().Where("id", playlistTrack.ID).Update(database.Map{
			"position": playlistTrack.Position - 1,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsLike(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(playlist.Public || authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Check if playlist already liked
	playlistLike := models.PlaylistLikeModel().Where("playlist_id", c.Params("playlistID")).Where("user_id", authUser.ID).First()
	if playlistLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like playlist
	models.PlaylistLikeModel().Create(database.Map{
		"playlist_id": c.Params("playlistID"),
		"user_id":     authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsLikeDelete(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(playlist.Public || authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Check if playlist not liked
	playlistLike := models.PlaylistLikeModel().Where("playlist_id", c.Params("playlistID")).Where("user_id", authUser.ID).First()
	if playlistLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.PlaylistLikeModel().Where("playlist_id", c.Params("playlistID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
