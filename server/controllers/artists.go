package controllers

import (
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func ArtistsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.ArtistModel(c).WhereRaw("`name` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`name`)").Paginate(page, limit))
}

func ArtistsShow(c *fiber.Ctx) error {
	artist := models.ArtistModel(c).With("albums", "top_tracks").Find(c.Params("artistID"))
	if artist == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(artist)
}

func ArtistsLike(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

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
	newArtistLike := models.ArtistLike{
		ArtistID: c.Params("artistID"),
		UserID:   authUser.ID,
	}
	models.ArtistLikeModel().Create(&newArtistLike)

	return c.JSON(fiber.Map{"success": true})
}

func ArtistsLikeDelete(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

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
