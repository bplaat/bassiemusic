package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func ArtistsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.ArtistModel(c).With("liked").WhereRaw("`name` LIKE ?", "%"+query+"%")
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
	artist := models.ArtistModel(c).With("liked", "albums", "top_tracks").Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(artist)
}

type ArtistsUpdateParams struct {
	Name     string `form:"name"`
	Synced   string `form:"synced"`
	DeezerID string `form:"deezer_id"`
}

func ArtistsUpdate(c *fiber.Ctx) error {
	// Check if artist exists
	artist := models.ArtistModel(c).Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Parse body
	var params ArtistsUpdateParams
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// TODO

	// Get updated artist
	return c.JSON(models.ArtistModel(c).Find(artist.ID))
}

func ArtistsDelete(c *fiber.Ctx) error {
	// Check if artist exists
	artist := models.ArtistModel(c).Find(c.Params("artistID"))
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
	models.ArtistModel(c).Where("id", artist.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if artist exists
	artist := models.ArtistModel(c).Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Check if artist already liked
	artistLike := models.ArtistLikeModel().Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).First()
	if artistLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like artist
	models.ArtistLikeModel().Create(database.Map{
		"artist_id": c.Params("artistID"),
		"user_id":   authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if artist exists
	artist := models.ArtistModel(c).Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}

	// Check if artist not liked
	artistLike := models.ArtistLikeModel().Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).First()
	if artistLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.ArtistLikeModel().Where("artist_id", c.Params("artistID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
