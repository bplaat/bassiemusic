package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Genre struct {
	ID          string    `column:"id,uuid" json:"id"`
	Name        string    `column:"name,string" json:"name"`
	DeezerID    int64     `column:"deezer_id,bigint" json:"-"`
	SmallImage  *string   `json:"small_image,omitempty"`
	MediumImage *string   `json:"medium_image,omitempty"`
	LargeImage  *string   `json:"large_image,omitempty"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
}

func GenreModel(*fiber.Ctx) *database.Model[Genre] {
	return (&database.Model[Genre]{
		TableName: "genres",
		Process: func(genre *Genre) {
			if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); err == nil {
				smallImage := fmt.Sprintf("%s/genres/small/%s.jpg", os.Getenv("STORAGE_URL"), genre.ID)
				genre.SmallImage = &smallImage
				mediumImage := fmt.Sprintf("%s/genres/medium/%s.jpg", os.Getenv("STORAGE_URL"), genre.ID)
				genre.MediumImage = &mediumImage
				largeImage := fmt.Sprintf("%s/genres/large/%s.jpg", os.Getenv("STORAGE_URL"), genre.ID)
				genre.LargeImage = &largeImage
			}
		},
	}).Init()
}
