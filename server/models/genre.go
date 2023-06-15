package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

// Genre
type Genre struct {
	ID          uuid.Uuid `column:"id" json:"id"`
	Name        string    `column:"name" json:"name"`
	DeezerID    int64     `column:"deezer_id" json:"deezer_id"`
	CreatedAt   time.Time `column:"created_at" json:"created_at"`
	SmallImage  *string   `json:"small_image,omitempty"`
	MediumImage *string   `json:"medium_image,omitempty"`
	LargeImage  *string   `json:"large_image,omitempty"`
	Liked       *bool     `json:"liked,omitempty"`
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
				liked := GenreLikeModel.Where("genre_id", genre.ID).Where("user_id", authUser.ID).Count() != 0
				genre.Liked = &liked
			}
		},
	},
}).Init()

// Genre Like
type GenreLike struct {
	ID        uuid.Uuid `column:"id"`
	GenreID   uuid.Uuid `column:"genre_id"`
	UserID    uuid.Uuid `column:"user_id"`
	CreatedAt time.Time `column:"created_at"`
}

var GenreLikeModel *database.Model[GenreLike] = (&database.Model[GenreLike]{
	TableName: "genre_likes",
}).Init()
