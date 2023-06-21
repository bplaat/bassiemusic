package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/controllers/websocket"
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
		updates["album_id"] = *body.AlbumID
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

	// Checks if youtube id has been updated
	if body.YoutubeID != nil && track.YoutubeID.String != *body.YoutubeID {
		data := Data{
			YoutubeID: *body.YoutubeID,
			TrackID:   track.ID,
		}
		jsonData, _ := json.Marshal(data)
		downloadTask := models.DownloadTaskModel.Create(database.Map{
			"type":         models.DownloadTaskTypeYoutubeTrack,
			"data":         jsonData,
			"display_name": *body.Title,
			"status":       models.DownloadTaskStatusPending,
			"progress":     0,
		})

		// Broadcast new task message to all listening admins
		if err := websocket.BroadcastAdmin("download_tasks.create", downloadTask); err != nil {
			log.Println(err)
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
