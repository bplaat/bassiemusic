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

func PlaylistsIndex(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)
	query, page, limit := utils.ParseIndexVars(c)
	q := models.PlaylistModel(c).With("like", "user").WhereRaw("`name` LIKE ?", "%"+query+"%")
	if authUser.Role != "admin" {
		q = q.Where("public", true)
	}
	if c.Query("sort_by") == "created_at" {
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
	playlist := models.PlaylistModel(c).With("like", "user", "tracks").Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !playlist.Public && (authUser.Role != "admin" || playlist.UserID != authUser.ID) {
		return fiber.ErrUnauthorized
	}

	return c.JSON(playlist)
}

type PlaylistsEditParams struct {
	Name   string `form:"name" validate:"required,min=2"`
	Public string `form:"public" validate:"required"`
}

func PlaylistsEdit(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" || playlist.UserID != authUser.ID {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var params PlaylistsEditParams
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
	updates := database.Map{
		"name":   params.Name,
		"public": params.Public == "true",
	}
	models.PlaylistModel(c).Where("id", playlist.ID).Update(updates)

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
	if authUser.Role != "admin" || playlist.UserID != authUser.ID {
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
	if authUser.Role != "admin" && authUser.ID != playlist.ID {
		return fiber.ErrUnauthorized
	}

	// Remove old image file
	if playlist.ImageID != nil && *playlist.ImageID != "" {
		if err := os.Remove(fmt.Sprintf("storage/playlists/%s.jpg", *playlist.ImageID)); err != nil {
			log.Fatalln(err)
		}
	}

	// Save uploaded image file
	imageID := uuid.New()
	image, err := c.FormFile("image")
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = c.SaveFile(image, fmt.Sprintf("storage/playlists/%s.jpg", imageID.String())); err != nil {
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
	if authUser.Role != "admin" && authUser.ID != playlist.ID {
		return fiber.ErrUnauthorized
	}

	// Check if playlist has image
	if playlist.ImageID != nil && *playlist.ImageID != "" {
		// Remove old image file
		if err := os.Remove(fmt.Sprintf("storage/playlists/%s.jpg", *playlist.ImageID)); err != nil {
			log.Fatalln(err)
		}

		// Clear image id for playlist
		models.PlaylistModel(c).Where("id", playlist.ID).Update(database.Map{
			"image": nil,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsInsertTrackParams struct {
	TrackID  string `form:"track_id" validate:"required"`
	Position string `form:"position" validate:"required,integer"`
}

func PlaylistsInsertTrack(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" || playlist.UserID != authUser.ID {
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
	if models.TrackModel(c).Find(params.TrackID) != nil {
		log.Println("track_id not valid")
		return fiber.ErrBadRequest
	}

	// Move all exisiting track ids to the left
	playlistTracks := models.PlaylistTrackModel().Where("playlist_id", playlist.ID).WhereRaw("`position` >= ?", params.Position).Get()
	for _, playlistTrack := range playlistTracks {
		models.PlaylistTrackModel().Where("id", playlistTrack.ID).Update(database.Map{
			"position": playlistTrack.Position + 1,
		})
	}

	// Create playlist track
	models.PlaylistTrackModel().Create(database.Map{
		"playlist_id": playlist.ID,
		"track_id":    params.TrackID,
		"position":    params.TrackID,
	})
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsRemoveTrackParams struct {
	Position string `form:"position" validate:"required,integer"`
}

func PlaylistsRemoveTrack(c *fiber.Ctx) error {
	// Check if playlist exists
	playlist := models.PlaylistModel(c).Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if authUser.Role != "admin" || playlist.UserID != authUser.ID {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var params PlaylistsRemoveTrackParams
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

	// Remove playlist track
	models.PlaylistTrackModel().Where("playlist_id", playlist.ID).Where("position", params.Position).Delete()

	// Move all exisiting track ids to the right
	playlistTracks := models.PlaylistTrackModel().Where("playlist_id", playlist.ID).WhereRaw("`position` > ?", params.Position).Get()
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
	if !playlist.Public && (authUser.Role != "admin" || playlist.UserID != authUser.ID) {
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
	if !playlist.Public && (authUser.Role != "admin" || playlist.UserID != authUser.ID) {
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
