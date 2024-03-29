package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

// Artist
type Artist struct {
	ID          uuid.Uuid `column:"id" json:"id"`
	Name        string    `column:"name" json:"name"`
	Sync        bool      `column:"sync" json:"sync"`
	DeezerID    int64     `column:"deezer_id" json:"deezer_id"`
	CreatedAt   time.Time `column:"created_at" json:"created_at"`
	SmallImage  *string   `json:"small_image"`
	MediumImage *string   `json:"medium_image"`
	LargeImage  *string   `json:"large_image"`
	Liked       *bool     `json:"liked,omitempty"`
	Albums      *[]Album  `json:"albums,omitempty"`
	TopTracks   *[]Track  `json:"top_tracks,omitempty"`
}

var ArtistModel *database.Model[Artist]

func init() {
	ArtistModel = (&database.Model[Artist]{
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
		Relationships: map[string]database.ModelRelationshipFunc[Artist]{
			"liked": func(artist *Artist, args []any) {
				if len(args) > 0 {
					authUser := args[0].(*User)
					liked := ArtistLikeModel.Where("artist_id", artist.ID).Where("user_id", authUser.ID).Count() != 0
					artist.Liked = &liked
				}
			},
			"albums": func(artist *Artist, args []any) {
				albums := AlbumModel.With("artists", "genres").WhereInQuery("id", AlbumArtistModel.Select("album_id").Where("artist_id", artist.ID)).OrderByDesc("released_at").Get()
				artist.Albums = &albums
			},
			"top_tracks": func(artist *Artist, args []any) {
				topTracksQuery := TrackModel.With("artists", "album")
				if len(args) > 0 {
					authUser := args[0].(*User)
					topTracksQuery = topTracksQuery.WithArgs("liked", authUser)
				}
				topTracks := topTracksQuery.WhereInQuery("id", TrackArtistModel.Select("track_id").Where("artist_id", artist.ID)).OrderByDesc("plays").Limit(25).Get()
				artist.TopTracks = &topTracks
			},
		},
	}).Init()
}

// Artist Like
type ArtistLike struct {
	ID        uuid.Uuid `column:"id"`
	ArtistID  uuid.Uuid `column:"artist_id"`
	UserID    uuid.Uuid `column:"user_id"`
	CreatedAt time.Time `column:"created_at"`
}

var ArtistLikeModel *database.Model[ArtistLike] = (&database.Model[ArtistLike]{
	TableName: "artist_likes",
}).Init()
