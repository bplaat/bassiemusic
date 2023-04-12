package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

// Playlist
type Playlist struct {
	ID          string    `column:"id,uuid" json:"id"`
	UserID      string    `column:"user_id,uuid" json:"-"`
	Name        string    `column:"name,string" json:"name"`
	ImageID     *string   `column:"image,uuid" json:"-"`
	SmallImage  *string   `json:"small_image"`
	MediumImage *string   `json:"medium_image"`
	Public      bool      `column:"public,bool" json:"public"`
	Liked       *bool     `json:"liked,omitempty"`
	CreatedAt   time.Time `column:"created_at,timestamp" json:"created_at"`
	UpdatedAt   time.Time `column:"updated_at,timestamp" json:"-"`
	User        *User     `json:"user,omitempty"`
	Tracks      *[]Track  `json:"tracks,omitempty"`
}

func PlaylistModel(c *fiber.Ctx) *database.Model[Playlist] {
	return (&database.Model[Playlist]{
		TableName: "playlists",
		Process: func(playlist *Playlist) {
			if playlist.ImageID != nil && *playlist.ImageID != "" {
				if _, err := os.Stat(fmt.Sprintf("storage/playlists/original/%s", *playlist.ImageID)); err == nil {
					smallImage := fmt.Sprintf("%s/playlists/small/%s.jpg", os.Getenv("STORAGE_URL"), *playlist.ImageID)
					playlist.SmallImage = &smallImage
					mediumImage := fmt.Sprintf("%s/playlists/medium/%s.jpg", os.Getenv("STORAGE_URL"), *playlist.ImageID)
					playlist.MediumImage = &mediumImage
				}
			}
		},
		Relationships: map[string]database.ModelProcessFunc[Playlist]{
			"liked": func(playlist *Playlist) {
				if c != nil {
					authUser := c.Locals("authUser").(*User)
					liked := PlaylistLikeModel().Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).First() != nil
					playlist.Liked = &liked
				}
			},
			"user": func(playlist *Playlist) {
				playlist.User = UserModel().Find(playlist.UserID)
			},
			"tracks": func(playlist *Playlist) {
				tracks := TrackModel(c).Join("INNER JOIN `playlist_track` ON `tracks`.`id` = `playlist_track`.`track_id`").
					With("liked", "artists", "album").WhereRaw("`playlist_track`.`playlist_id` = UUID_TO_BIN(?)", playlist.ID).
					OrderByRaw("`playlist_track`.`position`").Get()
				playlist.Tracks = &tracks
			},
		},
	}).Init()
}

// Playlist Track
type PlaylistTrack struct {
	ID         string    `column:"id,uuid"`
	PlaylistID string    `column:"playlist_id,uuid"`
	Position   int       `column:"position,int"`
	TrackID    string    `column:"track_id,uuid"`
	CreatedAt  time.Time `column:"created_at,timestamp"`
}

func PlaylistTrackModel() *database.Model[PlaylistTrack] {
	return (&database.Model[PlaylistTrack]{
		TableName: "playlist_track",
	}).Init()
}

// Playlist Like
type PlaylistLike struct {
	ID         string    `column:"id,uuid"`
	PlaylistID string    `column:"playlist_id,uuid"`
	UserID     string    `column:"user_id,uuid"`
	CreatedAt  time.Time `column:"created_at,timestamp"`
}

func PlaylistLikeModel() *database.Model[PlaylistLike] {
	return (&database.Model[PlaylistLike]{
		TableName: "playlist_likes",
	}).Init()
}
