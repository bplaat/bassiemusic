package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func TracksIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total tracks count
	tracksCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `tracks` WHERE `title` LIKE ?", "%"+query+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer tracksCountQuery.Close()

	tracksCountQuery.Next()
	var total int64
	tracksCountQuery.Scan(&total)

	// Get tracks
	var tracksQuery *sql.Rows
	tracksQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at` FROM `tracks` WHERE `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
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
	trackQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at` FROM `tracks` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("trackID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer trackQuery.Close()

	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.TrackScan(c, trackQuery, true, true))
}

func TracksPlay(c *fiber.Ctx) error {
	trackQuery, err := database.Query("SELECT `plays` FROM `tracks` WHERE `id` = UUID_TO_BIN(?)", c.Params("trackID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer trackQuery.Close()

	if !trackQuery.Next() {
		return fiber.ErrNotFound
	}
	var plays int64
	trackQuery.Scan(&plays)
	database.Exec("UPDATE `tracks` SET `plays` = ? WHERE `id` = UUID_TO_BIN(?)", plays+1, c.Params("trackID"))

	return c.JSON(fiber.Map{
		"success": true,
	})
}
