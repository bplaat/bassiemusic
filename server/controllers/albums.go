package controllers

import (
	"fmt"
	"os"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func AlbumsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.AlbumModel(c).With("like", "artists", "genres").WhereRaw("`title` LIKE ?", "%"+query+"%")
	if c.Query("sort_by") == "released_at" {
		q = q.OrderBy("released_at")
	} else if c.Query("sort_by") == "released_at_desc" {
		q = q.OrderByDesc("released_at")
	} else if c.Query("sort_by") == "created_at" {
		q = q.OrderBy("created_at")
	} else if c.Query("sort_by") == "created_at_desc" {
		q = q.OrderByDesc("created_at")
	} else if c.Query("sort_by") == "title_desc" {
		q = q.OrderByRaw("LOWER(`title`) DESC")
	} else {
		q = q.OrderByRaw("LOWER(`title`)")
	}
	return c.JSON(q.Paginate(page, limit))
}

func AlbumsShow(c *fiber.Ctx) error {
	album := models.AlbumModel(c).With("like", "artists", "genres", "tracks").Find(c.Params("albumID"))
	if album == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(album)
}

func AlbumsDelete(c *fiber.Ctx) error {
	// Check if album exists
	album := models.AlbumModel(c).With("tracks").Find(c.Params("albumID"))
	if album == nil {
		return fiber.ErrNotFound
	}

	// Delete album cover if exists
	if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); err == nil {
		_ = os.Remove(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID))
		_ = os.Remove(fmt.Sprintf("storage/albums/medium/%s.jpg", album.ID))
		_ = os.Remove(fmt.Sprintf("storage/albums/large/%s.jpg", album.ID))
	}

	// Delete album tracks music if exists (the models will be delete with album delete)
	for _, track := range *album.Tracks {
		if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); err == nil {
			_ = os.Remove(fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
		}
	}

	// Delete album
	models.AlbumModel(c).Where("id", album.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
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
