package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type Album struct {
	ID          uuid.Uuid `column:"id" json:"id"`
	Type        AlbumType `column:"type" json:"-"`
	TypeString  string    `json:"type"`
	Title       string    `column:"title" json:"title"`
	ReleasedAt  time.Time `column:"released_at" json:"released_at"`
	Explicit    bool      `column:"explicit" json:"explicit"`
	DeezerID    int64     `column:"deezer_id" json:"deezer_id"`
	CreatedAt   time.Time `column:"created_at" json:"created_at"`
	SmallCover  *string   `json:"small_cover"`
	MediumCover *string   `json:"medium_cover"`
	LargeCover  *string   `json:"large_cover"`
	Liked       *bool     `json:"liked,omitempty"`
	Artists     *[]Artist `json:"artists,omitempty"`
	Genres      *[]Genre  `json:"genres,omitempty"`
	Tracks      *[]Track  `json:"tracks,omitempty"`
}

type AlbumType int

const AlbumTypeAlbum AlbumType = 0
const AlbumTypeEP AlbumType = 1
const AlbumTypeSingle AlbumType = 2

var AlbumModel *database.Model[Album]

func init() {
	AlbumModel = (&database.Model[Album]{
		TableName: "albums",
		Process: func(album *Album) {
			if album.Type == AlbumTypeAlbum {
				album.TypeString = "album"
			}
			if album.Type == AlbumTypeEP {
				album.TypeString = "ep"
			}
			if album.Type == AlbumTypeSingle {
				album.TypeString = "single"
			}

			if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); err == nil {
				smallCover := fmt.Sprintf("%s/albums/small/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)
				album.SmallCover = &smallCover
				mediumCover := fmt.Sprintf("%s/albums/medium/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)
				album.MediumCover = &mediumCover
				largeCover := fmt.Sprintf("%s/albums/large/%s.jpg", os.Getenv("STORAGE_URL"), album.ID)
				album.LargeCover = &largeCover
			}
		},
		Relationships: map[string]database.ModelRelationshipFunc[Album]{
			"liked": func(album *Album, args []any) {
				if len(args) > 0 {
					authUser := args[0].(*User)
					liked := AlbumLikeModel.Where("album_id", album.ID).Where("user_id", authUser.ID).Count() != 0
					album.Liked = &liked
				}
			},
			"artists": func(album *Album, args []any) {
				albumArtists := AlbumArtistModel.Select("artist_id").Where("album_id", album.ID).OrderBy("position").Get()
				if len(albumArtists) == 0 {
					emptyArtists := []Artist{}
					album.Artists = &emptyArtists
					return
				}
				var artistIds []any
				for _, albumArtist := range albumArtists {
					artistIds = append(artistIds, albumArtist.ArtistID)
				}
				artists := ArtistModel.WhereIn("id", artistIds).Get()

				var orderedArtists []Artist
				for _, albumArtist := range albumArtists {
					for _, artist := range artists {
						if artist.ID.Equals(albumArtist.ArtistID) {
							orderedArtists = append(orderedArtists, artist)
							break
						}
					}
				}
				album.Artists = &orderedArtists
			},
			"genres": func(album *Album, args []any) {
				genres := GenreModel.WhereInQuery("id", AlbumGenreModel.Select("genre_id").Where("album_id", album.ID)).OrderByRaw("LOWER(`name`)").Get()
				album.Genres = &genres
			},
			"tracks": func(album *Album, args []any) {
				tracksQuery := TrackModel.With("artists")
				if len(args) > 0 {
					authUser := args[0].(*User)
					tracksQuery = tracksQuery.WithArgs("liked", authUser)
				}
				tracks := tracksQuery.Where("album_id", album.ID).OrderByRaw("`disk`, `position`").Get()
				album.Tracks = &tracks
			},
		},
	}).Init()
}

// Album artist
type AlbumArtist struct {
	ID       uuid.Uuid `column:"id"`
	AlbumID  uuid.Uuid `column:"album_id"`
	ArtistID uuid.Uuid `column:"artist_id"`
	Position int       `column:"position"`
}

var AlbumArtistModel *database.Model[AlbumArtist] = (&database.Model[AlbumArtist]{
	TableName: "album_artist",
}).Init()

// Album genre
type AlbumGenre struct {
	ID      uuid.Uuid `column:"id"`
	AlbumID uuid.Uuid `column:"album_id"`
	GenreID uuid.Uuid `column:"genre_id"`
}

var AlbumGenreModel *database.Model[AlbumGenre] = (&database.Model[AlbumGenre]{
	TableName: "album_genre",
}).Init()

// Album Like
type AlbumLike struct {
	ID        uuid.Uuid `column:"id"`
	AlbumID   uuid.Uuid `column:"album_id"`
	UserID    uuid.Uuid `column:"user_id"`
	CreatedAt time.Time `column:"created_at"`
}

var AlbumLikeModel *database.Model[AlbumLike] = (&database.Model[AlbumLike]{
	TableName: "album_likes",
}).Init()
