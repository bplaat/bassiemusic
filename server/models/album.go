package models

import (
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Album struct {
	ID          string    `column:"id,uuid" json:"id"`
	Type        string    `column:"type,enum:album|ep|single" json:"type"`
	Title       string    `column:"title,string" json:"title"`
	ReleasedAt  time.Time `column:"released_at,date" json:"released_at"`
	Explicit    bool      `column:"explicit,bool" json:"explicit"`
	DeezerID    int64     `column:"deezer_id,bigint" json:"-"`
	SmallCover  string    `json:"small_cover"`
	MediumCover string    `json:"medium_cover"`
	LargeCover  string    `json:"large_cover"`
	Liked       bool      `json:"liked"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
	Artists     []Artist  `json:"artists,omitempty"`
	Genres      []Genre   `json:"genres,omitempty"`
	Tracks      []Track   `json:"tracks,omitempty"`
}

type AlbumType int

const AlbumTypeAlbum AlbumType = 0
const AlbumTypeEP AlbumType = 1
const AlbumTypeSingle AlbumType = 2

func AlbumModel(c *fiber.Ctx) database.Model[Album] {
	return database.Model[Album]{
		TableName: "albums",
		Process: func(album *Album) {
			if c != nil {
				album.SmallCover = fmt.Sprintf("%s/storage/albums/small/%s.jpg", c.BaseURL(), album.ID)
				album.MediumCover = fmt.Sprintf("%s/storage/albums/medium/%s.jpg", c.BaseURL(), album.ID)
				album.LargeCover = fmt.Sprintf("%s/storage/albums/large/%s.jpg", c.BaseURL(), album.ID)

				album.Liked = AlbumLikeModel().Where("album_id", album.ID).Where("user_id", AuthUser(c).ID).First() != nil
			}
		},
		Relationships: map[string]database.QueryBuilderProcess[Album]{
			"artists": func(album *Album) {
				album.Artists = ArtistModel(c).WhereIn("album_artist", "artist_id", "album_id", album.ID).OrderByRaw("LOWER(`name`)").Get()
			},
			"genres": func(album *Album) {
				album.Genres = GenreModel(c).WhereIn("album_genre", "genre_id", "album_id", album.ID).OrderByRaw("LOWER(`name`)").Get()
			},
			"tracks": func(album *Album) {
				album.Tracks = TrackModel(c).With("artists").Where("album_id", album.ID).OrderByRaw("`disk`, `position`").Get()
			},
		},
	}.Init()
}

// Album Like
type AlbumLike struct {
	ID        string    `column:"id,uuid"`
	AlbumID   string    `column:"album_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func AlbumLikeModel() database.Model[AlbumLike] {
	return database.Model[AlbumLike]{
		TableName: "album_likes",
	}.Init()
}
