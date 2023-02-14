package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Album struct {
	ID          string    `column:"id,uuid" json:"id"`
	TypeInt     AlbumType `column:"type,int" json:"-"`
	Type        string    `json:"type"`
	Title       string    `column:"title,string" json:"title"`
	ReleasedAt  time.Time `column:"released_at,date" json:"released_at"`
	Explicit    bool      `column:"explicit,bool" json:"explicit"`
	Verified    bool      `column:"verified,bool" json:"verified"`
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

func AlbumModel(c *fiber.Ctx) *database.Model[Album] {
	return (&database.Model[Album]{
		TableName: "albums",
		Process: func(album *Album) {
			if album.TypeInt == AlbumTypeAlbum {
				album.Type = "album"
			}
			if album.TypeInt == AlbumTypeEP {
				album.Type = "ep"
			}
			if album.TypeInt == AlbumTypeSingle {
				album.Type = "single"
			}

			album.SmallCover = fmt.Sprintf("%s/albums/small/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)
			album.MediumCover = fmt.Sprintf("%s/albums/medium/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)
			album.LargeCover = fmt.Sprintf("%s/albums/large/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)

			if c != nil {
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
	}).Init()
}

// Album artist
type AlbumArtist struct {
	ID       string `column:"id,uuid"`
	AlbumID  string `column:"album_id,uuid"`
	ArtistID string `column:"artist_id,uuid"`
}

func AlbumArtistModel() *database.Model[AlbumArtist] {
	return (&database.Model[AlbumArtist]{
		TableName: "album_artist",
	}).Init()
}

// Album genre
type AlbumGenre struct {
	ID      string `column:"id,uuid"`
	AlbumID string `column:"album_id,uuid"`
	GenreID string `column:"genre_id,uuid"`
}

func AlbumGenreModel() *database.Model[AlbumGenre] {
	return (&database.Model[AlbumGenre]{
		TableName: "album_genre",
	}).Init()
}

// Album Like
type AlbumLike struct {
	ID        string    `column:"id,uuid"`
	AlbumID   string    `column:"album_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func AlbumLikeModel() *database.Model[AlbumLike] {
	return (&database.Model[AlbumLike]{
		TableName: "album_likes",
	}).Init()
}
