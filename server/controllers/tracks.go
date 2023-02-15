package controllers

import (
	"strconv"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func TracksIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.TrackModel(c).With("like", "artists", "album").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("`plays` DESC, `updated_at` DESC").Paginate(page, limit))
}

func TracksShow(c *fiber.Ctx) error {
	track := models.TrackModel(c).With("like", "artists", "album").Find(c.Params("trackID"))
	if track == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(track)
}

func TracksLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

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
	models.TrackLikeModel().Create(database.Map{
		"track_id": c.Params("trackID"),
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func TracksLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

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
	authUser := c.Locals("authUser").(*models.User)

	// Get position query variable
	var position float32
	if positionFloat, err := strconv.ParseFloat(c.Query("position", "0"), 32); err == nil {
		position = float32(positionFloat)
	}

	// Handle track play
	models.HandleTrackPlay(authUser, c.Params("trackID"), position)
	return c.JSON(fiber.Map{"success": true})
}
