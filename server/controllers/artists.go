package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func ArtistsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total artists count
	artistsCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `artists` WHERE `name` LIKE ?", "%"+query+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsCountQuery.Close()

	artistsCountQuery.Next()
	var total int64
	artistsCountQuery.Scan(&total)

	// Get artists
	var artistsQuery *sql.Rows
	artistsQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.ArtistsScan(c, artistsQuery, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func ArtistsShow(c *fiber.Ctx) error {
	artistQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("artistID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer artistQuery.Close()

	if !artistQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.ArtistScan(c, artistQuery, true))
}
