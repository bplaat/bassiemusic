package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Genre struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DeezerID    int64     `json:"-"`
	SmallImage  string    `json:"small_image"`
	MediumImage string    `json:"medium_image"`
	LargeImage  string    `json:"large_image"`
	CreatedAt   time.Time `json:"created_at"`
	Albums      []Album   `json:"albums,omitempty"`
}

func GenreScan(c *fiber.Ctx, genreQuery *sql.Rows, withAlbums bool) Genre {
	var genre Genre
	genreQuery.Scan(&genre.ID, &genre.Name, &genre.DeezerID, &genre.CreatedAt)
	if c != nil {
		genre.SmallImage = fmt.Sprintf("%s/storage/genres/small/%s.jpg", c.BaseURL(), genre.ID)
		genre.MediumImage = fmt.Sprintf("%s/storage/genres/medium/%s.jpg", c.BaseURL(), genre.ID)
		genre.LargeImage = fmt.Sprintf("%s/storage/genres/large/%s.jpg", c.BaseURL(), genre.ID)
	}
	if withAlbums {
		genre.Albums = GenreAlbums(c, &genre)
	}
	return genre
}

func GenresScan(c *fiber.Ctx, genresQuery *sql.Rows, withAlbums bool) []Genre {
	genres := []Genre{}
	for genresQuery.Next() {
		genres = append(genres, GenreScan(c, genresQuery, withAlbums))
	}
	return genres
}

func GenreAlbums(c *fiber.Ctx, genre *Genre) []Album {
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums` WHERE `id` IN (SELECT `album_id` FROM `album_genre` WHERE `genre_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", genre.ID)
	defer albumsQuery.Close()
	return AlbumsScan(c, albumsQuery, true, true, false)
}
