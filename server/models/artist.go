package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

// Artist
type Artist struct {
	ID          string    `column:"id,uuid" json:"id"`
	Name        string    `column:"name,string" json:"name"`
	Sync        bool      `column:"sync,bool" json:"sync"`
	DeezerID    int64     `column:"deezer_id,bigint" json:"-"`
	SmallImage  *string   `json:"small_image"`
	MediumImage *string   `json:"medium_image"`
	LargeImage  *string   `json:"large_image"`
	Liked       *bool     `json:"liked,omitempty"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
	Albums      *[]Album  `json:"albums,omitempty"`
	TopTracks   *[]Track  `json:"top_tracks,omitempty"`
}

func ArtistModel(c *fiber.Ctx) *database.Model[Artist] {
	return (&database.Model[Artist]{
		TableName: "artists",
		Process: func(artist *Artist) {
			if _, err := os.Stat(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID)); err == nil {
				smallImage := fmt.Sprintf("%s/artists/small/%s.jpg", os.Getenv("STORAGE_URL"), artist.ID)
				artist.SmallImage = &smallImage
				mediumImage := fmt.Sprintf("%s/artists/medium/%s.jpg", os.Getenv("STORAGE_URL"), artist.ID)
				artist.MediumImage = &mediumImage
				largeImage := fmt.Sprintf("%s/artists/large/%s.jpg", os.Getenv("STORAGE_URL"), artist.ID)
				artist.LargeImage = &largeImage
			}
		},
		Relationships: map[string]database.QueryBuilderProcess[Artist]{
			"like": func(artist *Artist) {
				if c != nil {
					authUser := c.Locals("authUser").(*User)
					liked := ArtistLikeModel().Where("artist_id", artist.ID).Where("user_id", authUser.ID).First() != nil
					artist.Liked = &liked
				}
			},
			"albums": func(artist *Artist) {
				albums := AlbumModel(c).With("artists", "genres").WhereIn("album_artist", "album_id", "artist_id", artist.ID).OrderByDesc("released_at").Get()
				artist.Albums = &albums
			},
			"top_tracks": func(artist *Artist) {
				topTracks := TrackModel(c).With("like", "artists", "album").WhereIn("track_artist", "track_id", "artist_id", artist.ID).OrderByDesc("plays").Limit("5").Get()
				artist.TopTracks = &topTracks
			},
		},
	}).Init()
}

// Artist Like
type ArtistLike struct {
	ID        string    `column:"id,uuid"`
	ArtistID  string    `column:"artist_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func ArtistLikeModel() *database.Model[ArtistLike] {
	return (&database.Model[ArtistLike]{
		TableName: "artist_likes",
	}).Init()
}
