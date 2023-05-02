package controllers

import (
	"fmt"
	_ "image/png"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func PlaylistsIndex(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)
	query, page, limit := utils.ParseIndexVars(c)
	q := models.PlaylistModel.WithArgs("liked", c.Locals("authUser")).With("user").WhereRaw("`name` LIKE ?", "%"+query+"%")
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

type PlaylistsCreateBody struct {
	Name   string `form:"name" validate:"required|min:2"`
	Public string `form:"public" validate:"required|boolean"`
}

func PlaylistsCreate(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse body
	var body PlaylistsCreateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create playlist
	fields := database.Map{
		"user_id": authUser.ID,
		"name":    body.Name,
		"public":  body.Public == "true",
	}
	if imageFile, err := c.FormFile("image"); err == nil {
		// Store image when given
		imageID := uuid.New()
		if err := utils.StoreUploadedImage(c, "playlists", imageID.String(), imageFile, false); err != nil {
			return err
		}
		fields["image"] = imageID.String()
	}
	return c.JSON(models.PlaylistModel.Create(fields))
}

func PlaylistsShow(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.WithArgs("liked", c.Locals("authUser")).With("user").
		WithArgs("tracks", c.Locals("authUser")).Find(c.Params("playlistID"))
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

type PlaylistsUpdateBody struct {
	Name   *string `form:"name" validate:"min:2"`
	Public *string `form:"public" validate:"boolean"`
	Image  *string `form:"image" validate:"nullable"`
}

func PlaylistsUpdate(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var body PlaylistsUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, playlist, &body); err != nil {
		return nil
	}

	// Update playlist
	updates := database.Map{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Public != nil {
		updates["public"] = *body.Public == "true"
	}
	if imageFile, err := c.FormFile("image"); err == nil {
		// Store new image
		imageID := uuid.New()
		if err := utils.StoreUploadedImage(c, "playlists", imageID.String(), imageFile, false); err != nil {
			return err
		}
		updates["image"] = imageID.String()

		// Remove old image file
		if playlist.ImageID.Valid {
			_ = os.Remove(fmt.Sprintf("storage/playlists/original/%s", playlist.ImageID.String))
			_ = os.Remove(fmt.Sprintf("storage/playlists/small/%s.jpg", playlist.ImageID.String))
			_ = os.Remove(fmt.Sprintf("storage/playlists/medium/%s.jpg", playlist.ImageID.String))
		}
	}
	if body.Image != nil && *body.Image == "" {
		if playlist.ImageID.Valid {
			// Remove old image file
			_ = os.Remove(fmt.Sprintf("storage/playlists/original/%s", playlist.ImageID.String))
			_ = os.Remove(fmt.Sprintf("storage/playlists/small/%s.jpg", playlist.ImageID.String))
			_ = os.Remove(fmt.Sprintf("storage/playlists/medium/%s.jpg", playlist.ImageID.String))

			updates["image"] = nil
		}
	}
	models.PlaylistModel.Where("id", playlist.ID).Update(updates)

	// Get updated playlist
	return c.JSON(models.PlaylistModel.Find(playlist.ID))
}

func PlaylistsDelete(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Delete playlist
	models.PlaylistModel.Where("id", playlist.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsAppendTrackBody struct {
	TrackID string `form:"track_id" validate:"required|uuid|exists:tracks,id"`
}

func PlaylistsAppendTrack(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var body PlaylistsAppendTrackBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Trigger playlist update
	oldName := playlist.Name
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Create playlist track
	models.PlaylistTrackModel.Create(database.Map{
		"playlist_id": playlist.ID,
		"track_id":    body.TrackID,
		"position":    models.PlaylistTrackModel.Where("playlist_id", playlist.ID).Count() + 1,
	})
	return c.JSON(fiber.Map{"success": true})
}

type PlaylistsInsertTrackBody struct {
	TrackID string `form:"track_id" validate:"required|uuid|exists:tracks,id"`
}

func PlaylistsInsertTrack(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Parse body
	var body PlaylistsInsertTrackBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
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
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Move all existing track ids to the left
	playlistTracks := models.PlaylistTrackModel.Where("playlist_id", playlist.ID).WhereRaw("`position` >= ?", position).Get()
	for _, playlistTrack := range playlistTracks {
		models.PlaylistTrackModel.Where("id", playlistTrack.ID).Update(database.Map{
			"position": playlistTrack.Position + 1,
		})
	}

	// Create playlist track
	models.PlaylistTrackModel.Create(database.Map{
		"playlist_id": playlist.ID,
		"track_id":    body.TrackID,
		"position":    position,
	})
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsRemoveTrack(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
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
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": "~" + oldName,
	})
	models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
		"name": oldName,
	})

	// Remove playlist track
	models.PlaylistTrackModel.Where("playlist_id", playlist.ID).Where("position", position).Delete()

	// Move all existing track ids to the right
	playlistTracks := models.PlaylistTrackModel.Where("playlist_id", playlist.ID).WhereRaw("`position` > ?", position).Get()
	for _, playlistTrack := range playlistTracks {
		models.PlaylistTrackModel.Where("id", playlistTrack.ID).Update(database.Map{
			"position": playlistTrack.Position - 1,
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsLike(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(playlist.Public || authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Check if playlist already liked
	playlistLike := models.PlaylistLikeModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).First()
	if playlistLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like playlist
	models.PlaylistLikeModel.Create(database.Map{
		"playlist_id": playlist.ID,
		"user_id":     authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func PlaylistsLikeDelete(c *fiber.Ctx) error {
	// Check if playlist id is valid uuid
	if !uuid.IsValid(c.Params("playlistID")) {
		return fiber.ErrBadRequest
	}

	// Check if playlist exists
	playlist := models.PlaylistModel.Find(c.Params("playlistID"))
	if playlist == nil {
		return fiber.ErrNotFound
	}

	// Check auth
	authUser := c.Locals("authUser").(*models.User)
	if !(playlist.Public || authUser.Role == models.UserRoleAdmin || playlist.UserID == authUser.ID) {
		return fiber.ErrUnauthorized
	}

	// Check if playlist not liked
	playlistLike := models.PlaylistLikeModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).First()
	if playlistLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.PlaylistLikeModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
