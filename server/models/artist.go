package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

type Artist struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DeezerID    int64     `json:"-"`
	SmallImage  string    `json:"small_image"`
	MediumImage string    `json:"medium_image"`
	LargeImage  string    `json:"large_image"`
	Liked       bool      `json:"liked"`
	CreatedAt   time.Time `json:"created_at"`
	Albums      []Album   `json:"albums,omitempty"`
	TopTracks   []Track   `json:"top_tracks,omitempty"`
}

func ArtistScan(c *fiber.Ctx, artistQuery *sql.Rows, withAlbums bool, withTopTracks bool) Artist {
	var artist Artist
	artistQuery.Scan(&artist.ID, &artist.Name, &artist.DeezerID, &artist.CreatedAt)
	if c != nil {
		artist.SmallImage = fmt.Sprintf("%s/storage/artists/small/%s.jpg", c.BaseURL(), artist.ID)
		artist.MediumImage = fmt.Sprintf("%s/storage/artists/medium/%s.jpg", c.BaseURL(), artist.ID)
		artist.LargeImage = fmt.Sprintf("%s/storage/artists/large/%s.jpg", c.BaseURL(), artist.ID)
		if withAlbums {
			artist.Albums = ArtistAlbums(c, &artist)
		}
		if withTopTracks {
			artist.TopTracks = ArtistTopTracks(c, &artist)
		}
		artist.Liked = ArtistLiked(c, &artist)
	}
	return artist
}

func ArtistsScan(c *fiber.Ctx, artistsQuery *sql.Rows, withAlbums bool, withTopTracks bool) []Artist {
	artists := []Artist{}
	for artistsQuery.Next() {
		artists = append(artists, ArtistScan(c, artistsQuery, withAlbums, withTopTracks))
	}
	return artists
}

func ArtistLiked(c *fiber.Ctx, artist *Artist) bool {
	authUser := AuthUser(c)
	artistLikeQuery := database.Query("SELECT `id` FROM `artist_likes` WHERE `artist_id` = UUID_TO_BIN(?) AND `user_id` = UUID_TO_BIN(?)", artist.ID, authUser.ID)
	defer artistLikeQuery.Close()
	return artistLikeQuery.Next()
}

func ArtistAlbums(c *fiber.Ctx, artist *Artist) []Album {
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums` "+
		"WHERE `id` IN (SELECT `album_id` FROM `album_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `released_at` DESC", artist.ID)
	defer albumsQuery.Close()
	return AlbumsScan(c, albumsQuery, true, true, false)
}

func ArtistTopTracks(c *fiber.Ctx, artist *Artist) []Track {
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks` "+
		"WHERE `id` IN (SELECT `track_id` FROM `track_artist` WHERE `artist_id` = UUID_TO_BIN(?)) ORDER BY `plays` DESC LIMIT 5", c.Params("artistID"))
	return TracksScan(c, tracksQuery, true, true)
}
