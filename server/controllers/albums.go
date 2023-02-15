package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func AlbumsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	return c.JSON(models.AlbumModel(c).With("like", "artists", "genres").WhereRaw("`title` LIKE ?", "%"+query+"%").OrderByRaw("LOWER(`title`)").Paginate(page, limit))
}

func AlbumsShow(c *fiber.Ctx) error {
	album := models.AlbumModel(c).With("like", "artists", "genres", "tracks").Find(c.Params("albumID"))
	if album == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(album)
}

func AlbumsLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if album exists
	album := models.AlbumModel(c).Find(c.Params("albumID"))
	if album == nil {
		return fiber.ErrNotFound
	}

	// Check if album already liked
	albumLike := models.AlbumLikeModel().Where("album_id", c.Params("albumID")).Where("user_id", authUser.ID).First()
	if albumLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like album
	models.AlbumLikeModel().Create(database.Map{
		"album_id": c.Params("albumID"),
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func AlbumsLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Check if album exists
	album := models.AlbumModel(c).Find(c.Params("albumID"))
	if album == nil {
		return fiber.ErrNotFound
	}

	// Check if album not liked
	albumLike := models.AlbumLikeModel().Where("album_id", c.Params("albumID")).Where("user_id", authUser.ID).First()
	if albumLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.AlbumLikeModel().Where("album_id", c.Params("albumID")).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
