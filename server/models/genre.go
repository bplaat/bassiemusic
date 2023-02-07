package models

import (
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Genre struct {
	ID          string    `column:"id,uuid" json:"id"`
	Name        string    `column:"name,string" json:"name"`
	DeezerID    int64     `column:"deezer_id,bigint" json:"-"`
	SmallImage  string    `json:"small_image"`
	MediumImage string    `json:"medium_image"`
	LargeImage  string    `json:"large_image"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
	Albums      []Album   `json:"albums,omitempty"`
}

func GenreModel(c *fiber.Ctx) database.Model[Genre] {
	return database.Model[Genre]{
		TableName: "genres",
		Process: func(genre *Genre) {
			if c != nil {
				genre.SmallImage = fmt.Sprintf("%s/storage/artists/small/%s.jpg", c.BaseURL(), genre.ID)
				genre.MediumImage = fmt.Sprintf("%s/storage/artists/medium/%s.jpg", c.BaseURL(), genre.ID)
				genre.LargeImage = fmt.Sprintf("%s/storage/artists/large/%s.jpg", c.BaseURL(), genre.ID)
			}
		},
		Relationships: map[string]database.QueryBuilderProcess[Genre]{
			"albums": func(genre *Genre) {
				genre.Albums = AlbumModel(c).With("artists", "genres").WhereIn("album_genre", "album_id", "genre_id", genre.ID).OrderByDesc("released_at").Get()
			},
		},
	}.Init()
}
