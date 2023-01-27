package controllers

import (
	"database/sql"
	"log"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/gofiber/fiber/v2"
)

func AlbumsIndex(c *fiber.Ctx) error {
	query, page, limit := utils.ParseIndexVars(c)

	// Get total albums count
	albumsCountQuery, err := database.Query("SELECT COUNT(`id`) FROM `albums` WHERE `title` LIKE ?", "%"+query+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsCountQuery.Close()

	albumsCountQuery.Next()
	var total int64
	albumsCountQuery.Scan(&total)

	// Get albums
	var albumsQuery *sql.Rows
	albumsQuery, err = database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `title` LIKE ? ORDER BY LOWER(`title`) LIMIT ?, ?", "%"+query+"%", (page-1)*limit, limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()

	// Return response
	return c.JSON(&fiber.Map{
		"data": models.AlbumsScan(c, albumsQuery, true, true, false),
		"pagination": &fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func AlbumsShow(c *fiber.Ctx) error {
	albumQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `id` = UUID_TO_BIN(?) LIMIT 1", c.Params("albumID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer albumQuery.Close()

	if !albumQuery.Next() {
		return fiber.ErrNotFound
	}
	return c.JSON(models.AlbumScan(c, albumQuery, true, true, true))
}
