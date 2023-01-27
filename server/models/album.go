package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Album struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	ReleasedAt time.Time `json:"released_at"`
	Explicit   bool      `json:"explicit"`
	Cover      string    `json:"cover"`
	CreatedAt  time.Time `json:"created_at"`
	Artists    []Artist  `json:"artists,omitempty"`
	Genres     []Genre   `json:"genres,omitempty"`
	Tracks     []Track   `json:"tracks,omitempty"`
}

type AlbumType int

const AlbumTypeAlbum AlbumType = 0
const AlbumTypeEP AlbumType = 1
const AlbumTypeSingle AlbumType = 2

func AlbumScan(c *fiber.Ctx, albumsQuery *sql.Rows, withArtists bool, withGenres bool, withTracks bool) Album {
	var album Album
	var albumType AlbumType
	albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt)
	if albumType == AlbumTypeAlbum {
		album.Type = "album"
	}
	if albumType == AlbumTypeEP {
		album.Type = "ep"
	}
	if albumType == AlbumTypeSingle {
		album.Type = "single"
	}
	album.Cover = fmt.Sprintf("%s/storage/albums/%s.jpg", c.BaseURL(), album.ID)
	if withArtists {
		album.Artists = AlbumArtists(c, &album)
	}
	if withGenres {
		album.Genres = AlbumGenres(c, &album)
	}
	if withTracks {
		album.Tracks = AlbumTracks(c, &album)
	}
	return album
}

func AlbumsScan(c *fiber.Ctx, albumsQuery *sql.Rows, withArtists bool, withGenres bool, withTracks bool) []Album {
	albums := []Album{}
	for albumsQuery.Next() {
		albums = append(albums, AlbumScan(c, albumsQuery, withArtists, withGenres, withTracks))
	}
	return albums
}

func AlbumArtists(c *fiber.Ctx, album *Album) []Artist {
	artistsQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `id` IN (SELECT `artist_id` FROM `album_artist` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	return ArtistsScan(c, artistsQuery, false)
}

func AlbumGenres(c *fiber.Ctx, album *Album) []Genre {
	genresQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `genres` WHERE `id` IN (SELECT `genre_id` FROM `album_genre` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()

	return GenresScan(c, genresQuery, false)
}

func AlbumTracks(c *fiber.Ctx, album *Album) []Track {
	tracksQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at` FROM `tracks` WHERE `album_id` = UUID_TO_BIN(?) ORDER BY `disk`, `position`", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracksQuery.Close()

	return TracksScan(c, tracksQuery, false, true)
}
