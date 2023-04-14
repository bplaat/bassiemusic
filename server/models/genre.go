package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
)

// Genre
type Genre struct {
	ID          string    `column:"id,uuid" json:"id"`
	Name        string    `column:"name,string" json:"name"`
	DeezerID    int64     `column:"deezer_id,bigint" json:"-"`
	SmallImage  *string   `json:"small_image,omitempty"`
	MediumImage *string   `json:"medium_image,omitempty"`
	LargeImage  *string   `json:"large_image,omitempty"`
	Liked       *bool     `json:"liked,omitempty"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
}

var GenreModel *database.Model[Genre] = (&database.Model[Genre]{
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
	Relationships: map[string]database.ModelRelationshipFunc[Genre]{
		"liked": func(genre *Genre, args []any) {
			if len(args) > 0 {
				authUser := args[0].(*User)
				liked := GenreLikeModel.Where("genre_id", genre.ID).Where("user_id", authUser.ID).First() != nil
				genre.Liked = &liked
			}
		},
	},
}).Init()

// Genre Like
type GenreLike struct {
	ID        string    `column:"id,uuid"`
	GenreID   string    `column:"genre_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

var GenreLikeModel *database.Model[GenreLike] = (&database.Model[GenreLike]{
	TableName: "genre_likes",
}).Init()
