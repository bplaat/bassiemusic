package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func ArtistsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.ArtistModel.WithArgs("liked", c.Locals("authUser")).WhereRaw("`name` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "sync" {
		q = q.OrderByRaw("`sync` DESC, LOWER(`name`)")
	} else if c.Query("sort_by") == "sync_desc" {
		q = q.OrderByRaw("`sync`, LOWER(`name`)")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderBy("created_at")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByDesc("created_at")
	} else if c.Query("sort_by") == "name_desc" {
		q = q.OrderByRaw("LOWER(`name`) DESC")
	} else {
		q = q.OrderByRaw("LOWER(`name`)")
	}
	return c.JSON(q.Paginate(page, limit))
}

func ArtistsShow(c *fiber.Ctx) error {
	artist := models.ArtistModel.WithArgs("liked", c.Locals("authUser")).With("albums").
		WithArgs("top_tracks", c.Locals("authUser")).Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(artist)
}

type ArtistsUpdateBody struct {
	Name     *string `form:"name" validate:"min:2"`
	Sync     *string `form:"sync" validate:"boolean"`
	DeezerID *string `form:"deezer_id" validate:"integer"`
}

func ArtistsUpdate(c *fiber.Ctx) error {
	// Check if artist exists
	artist := models.ArtistModel.Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Parse body
	var body ArtistsUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, artist, &body); err != nil {
		return nil
	}

	// Run updates
	updates := database.Map{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Sync != nil {
		updates["sync"] = *body.Sync == "true"
	}
	if body.DeezerID != nil {
		deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
		updates["deezer_id"] = deezerID
	}
	models.ArtistModel.Where("id", artist.ID).Update(updates)

	// Get updated artist
	return c.JSON(models.ArtistModel.Find(artist.ID))
}

func ArtistsDelete(c *fiber.Ctx) error {
	// Check if artist exists
	artist := models.ArtistModel.Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Delete artist image if exists
	if _, err := os.Stat(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID)); err == nil {
		_ = os.Remove(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID))
		_ = os.Remove(fmt.Sprintf("storage/artists/medium/%s.jpg", artist.ID))
		_ = os.Remove(fmt.Sprintf("storage/artists/large/%s.jpg", artist.ID))
	}

	// Delete artist
	models.ArtistModel.Where("id", artist.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if artist exists
	artist := models.ArtistModel.Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Check if artist already liked
	artistLike := models.ArtistLikeModel.Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).First()
	if artistLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like artist
	models.ArtistLikeModel.Create(database.Map{
		"artist_id": c.Params("artistID"),
		"user_id":   authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if artist exists
	artist := models.ArtistModel.Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Check if artist not liked
	artistLike := models.ArtistLikeModel.Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).First()
	if artistLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.ArtistLikeModel.Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
