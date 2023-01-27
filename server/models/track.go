package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Track struct {
	ID        string    `json:"id"`
	AlbumID   string    `json:"-"`
	Title     string    `json:"title"`
	Disk      int       `json:"disk"`
	Position  int       `json:"position"`
	Duration  int       `json:"duration"`
	Explicit  bool      `json:"explicit"`
	Plays     int64     `json:"plays"`
	Music     string    `json:"music"`
	CreatedAt time.Time `json:"created_at"`
	Album     *Album    `json:"album,omitempty"`
	Artists   []Artist  `json:"artists,omitempty"`
}

func TrackScan(c *fiber.Ctx, trackQuery *sql.Rows, withAlbum bool, withArtists bool) Track {
	var track Track
	trackQuery.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Disk, &track.Position, &track.Duration, &track.Explicit, &track.Plays, &track.CreatedAt)
	track.Music = fmt.Sprintf("%s/storage/tracks/%s.m4a", c.BaseURL(), track.ID)
	if withAlbum {
		album := TrackAlbum(c, &track)
		track.Album = &album
	}
	if withArtists {
		track.Artists = TrackArtists(c, &track)
	}
	return track
}

func TracksScan(c *fiber.Ctx, tracksQuery *sql.Rows, withAlbum bool, withArtists bool) []Track {
	tracks := []Track{}
	for tracksQuery.Next() {
		tracks = append(tracks, TrackScan(c, tracksQuery, withAlbum, withArtists))
	}
	return tracks
}

func TrackAlbum(c *fiber.Ctx, track *Track) Album {
	albumQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `created_at` FROM `albums` WHERE `id` = UUID_TO_BIN(?)", track.AlbumID)
	if err != nil {
		log.Fatalln(err)
	}
	defer albumQuery.Close()

	albumQuery.Next()
	return AlbumScan(c, albumQuery, true, true, false)
}

func TrackArtists(c *fiber.Ctx, track *Track) []Artist {
	artistsQuery, err := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `created_at` FROM `artists` WHERE `id` IN (SELECT `artist_id` FROM `track_artist` WHERE `track_id` = UUID_TO_BIN(?)) ORDER BY LOWER(`name`)", track.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer artistsQuery.Close()

	return ArtistsScan(c, artistsQuery, false)
}
