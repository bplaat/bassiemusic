package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Artist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	Albums    []Album   `json:"albums,omitempty"`
}

func ArtistScan(c *fiber.Ctx, artistQuery *sql.Rows, withAlbums bool) Artist {
	var artist Artist
	artistQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt)
	artist.Image = fmt.Sprintf("%s/storage/artists/%s.jpg", c.BaseURL(), artist.ID)
	if withAlbums {
		artist.Albums = ArtistAlbums(c, &artist)
	}
	return artist
}

func ArtistsScan(c *fiber.Ctx, artistsQuery *sql.Rows, withAlbums bool) []Artist {
	artists := []Artist{}
	for artistsQuery.Next() {
		artists = append(artists, ArtistScan(c, artistsQuery, withAlbums))
	}
	return artists
}

func ArtistAlbums(c *fiber.Ctx, artist *Artist) []Album {
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", artist.ID)
	defer albumsQuery.Close()
	return AlbumsScan(c, albumsQuery, true, true, false)
}
