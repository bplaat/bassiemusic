package controllers

import (
	"strconv"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func TracksIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.TrackModel(c).With("artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`title`)").Paginate(page, limit))
}

func TracksShow(c *fiber.Ctx) error {
	track := models.TrackModel(c).With("artists", "album").Find(c.Params("trackID"))
	if track == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(track)
}

func TracksLike(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if track exists
	track := models.TrackModel(c).Find(c.Params("trackID"))
	if track == nil {
		return fiber.ErrNotFound
	}

	// Check if track already liked
	trackLike := models.TrackLikeModel().Where("track_id", c.Params("trackID")).Where("user_id", authUser.ID).First()
	if trackLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like track
	newTrackLike := models.TrackLike{
		TrackID: c.Params("trackID"),
		UserID:  authUser.ID,
	}
	models.TrackLikeModel().Create(&newTrackLike)

	return c.JSON(fiber.Map{"success": true})
}

func TracksLikeDelete(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if track exists
	track := models.TrackModel(c).Find(c.Params("trackID"))
	if track == nil {
		return fiber.ErrNotFound
	}

	// Check if track not liked
	trackLike := models.TrackLikeModel().Where("track_id", c.Params("trackID")).Where("user_id", authUser.ID).First()
	if trackLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.TrackLikeModel().Where("track_id", c.Params("trackID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func TracksPlay(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if track exists
	track := models.TrackModel(c).Find(c.Params("trackID"))
	if track == nil {
		return fiber.ErrNotFound
	}

	// Parse position get variable
	var position float32
	if positionFloat, err := strconv.ParseFloat(c.Query("position", "0"), 32); err == nil {
		position = float32(positionFloat)
	}

	// Get user last track play and update if latest
	trackPlay := models.TrackPlayModel().Where("user_id", authUser.ID).OrderByDesc("created_at").First()
	if trackPlay != nil {
		if track.ID == trackPlay.TrackID {
			trackPlay.Position = position
			models.TrackPlayModel().Update(trackPlay)
			return c.JSON(fiber.Map{"success": true})
		}
	}

	// Create new track play
	newTrackPlay := models.TrackPlay{
		TrackID:  track.ID,
		UserID:   authUser.ID,
		Position: position,
	}
	models.TrackPlayModel().Create(&newTrackPlay)

	// Increment global track plays count
	track.Plays += 1
	models.TrackModel(c).Update(track)
	return c.JSON(fiber.Map{"success": true})
}
