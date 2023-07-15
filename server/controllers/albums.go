package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bplaat/bassiemusic/controllers/websocket"
	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/core/validation"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func AlbumsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)
	q := models.AlbumModel.WithArgs("liked", c.Locals("authUser")).
		With("artists", "genres").WhereRaw("`title` LIKE ?", "%"+query+"%")
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

type AlbumsCreateBody struct {
	Title      *string `form:"title" validate:"min:2"`
	Type       *string `form:"type" validate:"enum:album,ep,single"`
	ReleasedAt *string `form:"released_at" validate:"date"`
	Explicit   *string `form:"explicit" validate:"boolean"`
	DeezerID   *string `form:"deezer_id" validate:"integer"`
	Cover      *string `form:"cover" validate:"nullable"`
}

func AlbumsCreate(c *fiber.Ctx) error {
	// Parse body
	var body AlbumsCreateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStruct(c, &body); err != nil {
		return nil
	}

	// Create album
	albumID := uuid.New()
	var albumType models.AlbumType
	if body.Type != nil {
		if *body.Type == "album" {
			albumType = models.AlbumTypeAlbum
		}
		if *body.Type == "ep" {
			albumType = models.AlbumTypeEP
		}
		if *body.Type == "single" {
			albumType = models.AlbumTypeSingle
		}
	}
	releasedAt, _ := time.Parse("2006-01-02T15:04:05Z", *body.ReleasedAt)
	deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
	album := models.AlbumModel.Create(database.Map{
		"id":          albumID,
		"title":       *body.Title,
		"type":        albumType,
		"released_at": releasedAt,
		"explicit":    *body.Explicit == "true",
		"deezer_id":   deezerID,
	})

	// Store new album image
	if imageFile, err := c.FormFile("image"); err == nil {
		if err := utils.StoreUploadedImage(c, "albums", albumID, imageFile, true); err != nil {
			return err
		}
	}

	return c.JSON(album)
}

func AlbumsShow(c *fiber.Ctx) error {
	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.WithArgs("liked", c.Locals("authUser")).With("artists", "genres").
		WithArgs("tracks", c.Locals("authUser")).Find(albumID)
	if album == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(album)
}

func AlbumsTracksUpdate(c *fiber.Ctx) error {
	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.Find(albumID)
	if album == nil {
		return fiber.ErrNotFound
	}

	// Generate update task
	jsonData, err := json.Marshal(models.DownloadTaskData{
		AlbumID:  albumID.String(),
		DeezerID: album.DeezerID,
	})
	if err != nil {
		log.Fatalln(err)
	}

	downloadTask := models.DownloadTaskModel.Create(database.Map{
		"type":         models.DownloadTaskTypeUpdateDeezerAlbum,
		"data":         jsonData,
		"display_name": album.Title,
		"status":       models.DownloadTaskStatusPending,
		"progress":     0,
	})

	// Broadcast new task message to all listening admins
	if err := websocket.BroadcastAdmin("download_tasks.create", downloadTask); err != nil {
		log.Println(err)
	}

	// Return success response
	return c.JSON(fiber.Map{"success": true})
}

type AlbumsUpdateBody struct {
	Title      *string `form:"title" validate:"min:2"`
	Type       *string `form:"type" validate:"enum:album,ep,single"`
	ReleasedAt *string `form:"released_at" validate:"date"`
	Explicit   *string `form:"explicit" validate:"boolean"`
	DeezerID   *string `form:"deezer_id" validate:"integer"`
	Cover      *string `form:"cover" validate:"nullable"`
}

func AlbumsUpdate(c *fiber.Ctx) error {
	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.Find(albumID)
	if album == nil {
		return fiber.ErrNotFound
	}

	// Parse body
	var body AlbumsUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate body
	if err := validation.ValidateStructUpdates(c, album, &body); err != nil {
		return nil
	}

	// Run updates
	updates := database.Map{}
	if body.Title != nil {
		updates["title"] = *body.Title
	}
	if body.Type != nil {
		if *body.Type == "album" {
			updates["type"] = models.AlbumTypeAlbum
		}
		if *body.Type == "ep" {
			updates["type"] = models.AlbumTypeEP
		}
		if *body.Type == "single" {
			updates["type"] = models.AlbumTypeSingle
		}
	}
	if body.ReleasedAt != nil {
		releasedAt, _ := time.Parse("2006-01-02T15:04:05Z", *body.ReleasedAt)
		updates["released_at"] = releasedAt.Format("2006-01-02 15:04:05")
	}
	if body.Explicit != nil {
		updates["explicit"] = *body.Explicit == "true"
	}
	if body.DeezerID != nil {
		deezerID, _ := strconv.ParseInt(*body.DeezerID, 10, 64)
		updates["deezer_id"] = deezerID
	}
	if coverFile, err := c.FormFile("cover"); err == nil {
		// Remove old cover file
		if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); err == nil {
			_ = os.Remove(fmt.Sprintf("storage/albums/original/%s", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/medium/%s.jpg", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/large/%s.jpg", album.ID))
		}

		// Store new cover
		if err := utils.StoreUploadedImage(c, "albums", album.ID, coverFile, true); err != nil {
			return err
		}
	}
	if body.Cover != nil && *body.Cover == "" {
		// Remove old cover file
		if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); err == nil {
			_ = os.Remove(fmt.Sprintf("storage/albums/original/%s", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/medium/%s.jpg", album.ID))
			_ = os.Remove(fmt.Sprintf("storage/albums/large/%s.jpg", album.ID))
		}
	}
	models.AlbumModel.Where("id", album.ID).Update(updates)

	// Get updated album
	return c.JSON(models.AlbumModel.Find(album.ID))
}

func AlbumsDelete(c *fiber.Ctx) error {
	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.With("tracks").Find(albumID)
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
	models.AlbumModel.Where("id", album.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}

func AlbumsLike(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.Find(albumID)
	if album == nil {
		return fiber.ErrNotFound
	}

	// Check if album already liked
	albumLike := models.AlbumLikeModel.Where("album_id", album.ID).Where("user_id", authUser.ID).First()
	if albumLike != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Like album
	models.AlbumLikeModel.Create(database.Map{
		"album_id": album.ID,
		"user_id":  authUser.ID,
	})

	return c.JSON(fiber.Map{"success": true})
}

func AlbumsLikeDelete(c *fiber.Ctx) error {
	authUser := c.Locals("authUser").(*models.User)

	// Parse album id uuid
	albumID, err := uuid.Parse(c.Params("albumID"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Check if album exists
	album := models.AlbumModel.Find(albumID)
	if album == nil {
		return fiber.ErrNotFound
	}

	// Check if album not liked
	albumLike := models.AlbumLikeModel.Where("album_id", album.ID).Where("user_id", authUser.ID).First()
	if albumLike == nil {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete like
	models.AlbumLikeModel.Where("album_id", album.ID).Where("user_id", authUser.ID).Delete()
	return c.JSON(fiber.Map{"success": true})
}
