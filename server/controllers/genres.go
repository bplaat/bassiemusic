package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func GenresIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total genres count
	genresCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `genres` WHERE `name` LIKE ?", "%"+query+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer genresCountQuery.Close()

	genresCountQuery.Next()
	var total int64
	genresCountQuery.Scan(&total)

	// Get genres
	var genresQuery *sql.Rows
	genresQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `genres` WHERE `name` LIKE ? ORDER BY LOWER(`name`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.GenresScan(c, genresQuery, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GenresShow(c *fiber.Ctx) error {
	genreQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `genres` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("genreID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer genreQuery.Close()

	if !genreQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.GenreScan(c, genreQuery, true))
}
