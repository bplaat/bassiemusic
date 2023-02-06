package controllers

import (
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func TracksIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total tracks
	total := database.Count("SELECT COUNT(`id`) FROM `tracks` WHERE `title` LIKE ?", "%"+query+"%")

	// Get tracks
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks` WHERE `title` LIKE ? ORDER BY `plays` DESC, LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	defer tracksQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.TracksScan(c, tracksQuery, true, true),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func TracksShow(c *fiber.Ctx) error {
	// Check if track exists
	trackQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", c.Params("trackID"))
	defer trackQuery.Close()
	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}

	// Return response
	return c.JSON(models.TrackScan(c, trackQuery, true, true))
}

func TracksPlay(c *fiber.Ctx) error {
	trackQuery := database.Query("SELECT `plays` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", c.Params("trackID"))
	defer trackQuery.Close()

	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}
	var plays int64
	trackQuery.Scan(&plays)
	database.Exec("UPDATE `tracks` SET `plays` = ? WHERE `id` = UUID_TO_BIN(?)", plays+1, c.Params("trackID"))

	return c.JSON(fiber.Map{"success": true})
}

func TracksLike(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if track exists
	trackQuery := database.Query("SELECT `id` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", c.Params("trackID"))
	defer trackQuery.Close()
	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if track_likes binding exists
	trackLikeQuery := database.Query("SELECT `id` FROM `track_likes` WHERE `track_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("trackID"), authUser.ID)
	defer trackLikeQuery.Close()
	if trackLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Create track_likes binding
	database.Exec("INSERT INTO `track_likes` (`id`, `track_id`, `user_id`) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), UUID_TO_BIN(?))", c.Params("trackID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}

func TracksLikeDelete(c *fiber.Ctx) error {
	authUser := models.AuthUser(c)

	// Check if track exists
	trackQuery := database.Query("SELECT `id` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", c.Params("trackID"))
	defer trackQuery.Close()
	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}

	// Check if track_likes binding doesn't exists
	trackLikeQuery := database.Query("SELECT `id` FROM `track_likes` WHERE `track_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("trackID"), authUser.ID)
	defer trackLikeQuery.Close()
	if !trackLikeQuery.Next() {
		return c.JSON(fiber.Map{"success": true})
	}

	// Delete track_likes binding
	database.Exec("DELETE FROM `track_likes` WHERE `track_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", c.Params("trackID"), authUser.ID)

	// Send successfull response
	return c.JSON(fiber.Map{"success": true})
}
