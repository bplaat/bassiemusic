package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func TracksIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.TrackModel.WithArgs("liked", c.Locals("authUser")).With("artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "title" {
		q = q.OrderByRaw("LOWER(`title`)")
	} else if c.Query("sort_by") == "title_desc" {
		q = q.OrderByRaw("LOWER(`title`) DESC")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderBy("created_at")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByDesc("created_at")
	} else if c.Query("sort_by") == "plays" {
		q = q.OrderByRaw("`plays`, LOWER(`title`)")
	} else {
		q = q.OrderByRaw("`plays` DESC, LOWER(`title`)")
	}
	return c.JSON(q.Paginate(page, limit))
}

type TracksCreateBody struct {
	Title     *string `form:"title" validate:"min:2"`
	AlbumID   *string `form:"album_id" validate:"uuid|exists:albums,id"`
	Disk      *string `form:"disk" validate:"integer"`
	Position  *string `form:"position" validate:"integer"`
	Explicit  *string `form:"explicit" validate:"boolean"`
	DeezerID  *string `form:"deezer_id" validate:"integer"`
	YoutubeID *string `form:"youtube_id" validate:"nullable|min:11"`
}

func TracksCreate(c *fiber.Ctx) error {
	// Parse body
	var body TracksCreateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create track
	trackID := uuid.New()
	albumID, _ := uuid.Parse(*body.AlbumID)
	disk, _ := strconv.ParseInt(*body.Disk, 10, 32)
	position, _ := strconv.ParseInt(*body.Position, 10, 32)
	deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
	var youtubeID *string
	if *body.YoutubeID != "" {
		youtubeID = body.YoutubeID
	} else {
		youtubeID = nil
	}
	models.TrackModel.Create(database.Map{
		"id":         trackID,
		"title":      *body.Title,
		"album_id":   albumID,
		"disk":       int32(disk),
		"position":   int32(position),
		"explicit":   *body.Explicit == "true",
		"youtube_id": youtubeID,
		"deezer_id":  deezerID,
	})

	// Store new track image
	if imageFile, err := c.FormFile("image"); err == nil {
		if err := utils.StoreUploadedImage(c, "tracks", trackID, imageFile, true); err != nil {
			return err
		}
	}

	// Get new track
	return c.JSON(models.TrackModel.Find(trackID))
}

func TracksShow(c *fiber.Ctx) error {
	// Parse track id uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if track exists
	track := models.TrackModel.WithArgs("liked", c.Locals("authUser")).With("artists", "album").Find(trackID)
	if track == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(track)
}

type TracksUpdateBody struct {
	Title     *string `form:"title" validate:"min:2"`
	AlbumID   *string `form:"album_id" validate:"uuid|exists:albums,id"`
	Disk      *string `form:"disk" validate:"integer"`
	Position  *string `form:"position" validate:"integer"`
	Explicit  *string `form:"explicit" validate:"boolean"`
	DeezerID  *string `form:"deezer_id" validate:"integer"`
	YoutubeID *string `form:"youtube_id" validate:"nullable|min:11"`
}

func TracksUpdate(c *fiber.Ctx) error {
	// Parse track id uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if track exists
	track := models.TrackModel.Find(trackID)
	if track == nil {
		return fiber.ErrNotFound
	}

	// Parse body
	var body TracksUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, track, &body); err != nil {
		return nil
	}

	// Run updates
	updates := database.Map{}
	if body.Title != nil {
		updates["title"] = *body.Title
	}
	if body.AlbumID != nil {
		albumID, _ := uuid.Parse(*body.AlbumID)
		updates["album_id"] = albumID
	}
	if body.Disk != nil {
		disk, _ := strconv.ParseInt(*body.Disk, 10, 32)
		updates["disk"] = int32(disk)
	}
	if body.Position != nil {
		position, _ := strconv.ParseInt(*body.Position, 10, 32)
		updates["position"] = int32(position)
	}
	if body.Explicit != nil {
		updates["explicit"] = *body.Explicit == "true"
	}
	if body.DeezerID != nil {
		deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
		updates["deezer_id"] = deezerID
	}
	if body.YoutubeID != nil {
		if *body.YoutubeID != "" {
			updates["youtube_id"] = *body.YoutubeID
		} else {
			updates["youtube_id"] = nil
		}
	}
	models.TrackModel.Where("id", track.ID).Update(updates)

	// Get updated track
	return c.JSON(models.TrackModel.Find(track.ID))
}

func TracksDelete(c *fiber.Ctx) error {
	// Parse track id uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if track exists
	track := models.TrackModel.Find(trackID)
	if track == nil {
		return fiber.ErrNotFound
	}

	// Delete track music if exists
	if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); err == nil {
		_ = os.Remove(fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
	}

	// Delete track
	models.TrackModel.Where("id", track.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func TracksPlay(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if track id is valid uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Get position query variable
	position, err := strconv.ParseFloat(c.Query("position", "0"), 32)
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Handle track play
	models.HandleTrackPlay(authUser, trackID, float32(position))
	return c.JSON(fiber.Map{"success": true})
}

func TracksLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse track id uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if track exists
	track := models.TrackModel.Find(trackID)
	if track == nil {
		return fiber.ErrNotFound
	}

	// Check if track already liked
	trackLike := models.TrackLikeModel.Where("track_id", track.ID).Where("user_id", authUser.ID).First()
	if trackLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like track
	models.TrackLikeModel.Create(database.Map{
		"track_id": track.ID,
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func TracksLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse track id uuid
	trackID, err := uuid.Parse(c.Params("trackID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if track exists
	track := models.TrackModel.Find(trackID)
	if track == nil {
		return fiber.ErrNotFound
	}

	// Check if track not liked
	trackLike := models.TrackLikeModel.Where("track_id", track.ID).Where("user_id", authUser.ID).First()
	if trackLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.TrackLikeModel.Where("track_id", track.ID).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
