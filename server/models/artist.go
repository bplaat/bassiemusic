package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
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
					liked := ArtistLikeModel.Where("artist_id", artist.ID).Where("user_id", authUser.ID).First() != nil
					artist.Liked = &liked
				}
			},
			"albums": func(artist *Artist, args []any) {
				albums := AlbumModel.With("artists", "genres").WhereIn("album_artist", "album_id", "artist_id", artist.ID).OrderByDesc("released_at").Get()
				artist.Albums = &albums
			},
			"top_tracks": func(artist *Artist, args []any) {
				topTracksQuery := TrackModel.With("artists", "album")
				if len(args) > 0 {
					authUser := args[0].(*User)
					topTracksQuery = topTracksQuery.WithArgs("liked", authUser)
				}
				topTracks := topTracksQuery.WhereIn("track_artist", "track_id", "artist_id", artist.ID).OrderByDesc("plays").Limit(5).Get()
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

var ArtistLikeModel *database.Model[ArtistLike] = (&database.Model[ArtistLike]{
	TableName: "artist_likes",
}).Init()
