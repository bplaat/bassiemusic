package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string    `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Artist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Albums    []Album   `json:"albums,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Album struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	ReleasedAt time.Time `json:"released_at"`
	Explicit   bool      `json:"explicit"`
	Cover      string    `json:"cover"`
	Genres     []Genre   `json:"genres,omitempty"`
	Artists    []Artist  `json:"artists,omitempty"`
	Tracks     []Track   `json:"tracks,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AlbumType int

const AlbumTypeAlbum AlbumType = 0
const AlbumTypeEP AlbumType = 1
const AlbumTypeSingle AlbumType = 2

type Genre struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Albums    []Album   `json:"albums,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Track struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Disk      int       `json:"disk"`
	Position  int       `json:"position"`
	Duration  int       `json:"duration"`
	Explicit  bool      `json:"explicit"`
	Plays     int64     `json:"plays"`
	Music     string    `json:"music"`
	Album     *Album    `json:"album,omitempty"`
	Artists   []Artist  `json:"artists,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func artistAlbums(artist *Artist, c *fiber.Ctx) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", artist.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
	for albumsQuery.Next() {
		var album Album
		var albumType AlbumType
		albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt)
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
		albumArtists(&album, c)
		artist.Albums = append(artist.Albums, album)
	}
}

func albumGenres(album *Album, c *fiber.Ctx) {
	genresQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `genres` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `genre_id` FROM `album_genre` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer genresQuery.Close()
	for genresQuery.Next() {
		var genre Genre
		genresQuery.Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		genre.Image = fmt.Sprintf("%s/storage/genres/%s.jpg", c.BaseURL(), genre.ID)
		album.Genres = append(album.Genres, genre)
	}
}

func albumArtists(album *Album, c *fiber.Ctx) {
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `album_artist` WHERE `album_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("%s/storage/artists/%s.jpg", c.BaseURL(), artist.ID)
		album.Artists = append(album.Artists, artist)
	}
}

func albumTracks(album *Album, c *fiber.Ctx) {
	tracksQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `title`, `disk`, `position`, `duration`, `explicit`, `plays`, `created_at`, `updated_at` FROM `tracks` WHERE `deleted_at` IS NULL AND `album_id` = UUID_TO_BIN(?) ORDER BY `disk`, `position`", album.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracksQuery.Close()
	for tracksQuery.Next() {
		var track Track
		tracksQuery.Scan(&track.ID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt, &track.UpdatedAt)
		track.Music = fmt.Sprintf("%s/storage/tracks/%s.m4a", c.BaseURL(), track.ID)
		trackArtists(&track, c)
		album.Tracks = append(album.Tracks, track)
	}
}

func genreAlbums(genre *Genre, c *fiber.Ctx) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `album_id` FROM `album_genre` WHERE `genre_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", genre.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
	for albumsQuery.Next() {
		var album Album
		var albumType AlbumType
		albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt)
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
		albumArtists(&album, c)
		genre.Albums = append(genre.Albums, album)
	}
}

func trackAlbum(track *Track, albumID string, c *fiber.Ctx) {
	albumsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at`, `updated_at` FROM `albums` WHERE `deleted_at` IS NULL AND `id` = UUID_TO_BIN(?)", albumID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumsQuery.Close()
	var album Album
	var albumType AlbumType
	albumsQuery.Next()
	albumsQuery.Scan(&album.ID, &albumType, &album.Title, &album.ReleasedAt, &album.Explicit, &album.CreatedAt, &album.UpdatedAt)
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
	track.Album = &album
}

func trackArtists(track *Track, c *fiber.Ctx) {
	artistsQuery, err := db.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at`, `updated_at` FROM `artists` WHERE `deleted_at` IS NULL AND `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()
	for artistsQuery.Next() {
		var artist Artist
		artistsQuery.Scan(&artist.ID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt)
		artist.Image = fmt.Sprintf("%s/storage/artists/%s.jpg", c.BaseURL(), artist.ID)
		track.Artists = append(track.Artists, artist)
	}
}
