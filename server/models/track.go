package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/gofiber/fiber/v2"
)

// Track
type Track struct {
	ID        string    `column:"id,uuid" json:"id"`
	AlbumID   string    `column:"album_id,uuid" json:"-"`
	Title     string    `column:"title,string" json:"title"`
	Disk      int       `column:"disk,int" json:"disk"`
	Position  int       `column:"position,int" json:"position"`
	Duration  float32   `column:"duration,float" json:"duration"`
	Explicit  bool      `column:"explicit,bool" json:"explicit"`
	DeezerID  int64     `column:"deezer_id,bigint" json:"-"`
	YoutubeID string    `column:"youtube_id,string" json:"-"`
	Plays     int64     `column:"plays,bigint" json:"plays"`
	Music     string    `json:"music"`
	Liked     bool      `json:"liked"`
	CreatedAt time.Time `column:"created_at,timestamp" json:"created_at"`
	Album     *Album    `json:"album,omitempty"`
	Artists   []Artist  `json:"artists,omitempty"`
}

func TrackModel(c *fiber.Ctx) *database.Model[Track] {
	return (&database.Model[Track]{
		TableName: "tracks",
		Process: func(track *Track) {
			track.Music = fmt.Sprintf("%s/tracks/%s.m4a", os.Getenv("STORAGE_URL"), track.ID)

			if c != nil {
				authUser := c.Locals("authUser").(*User)
				track.Liked = TrackLikeModel().Where("track_id", track.ID).Where("user_id", authUser.ID).First() != nil
			}
		},
		Relationships: map[string]database.QueryBuilderProcess[Track]{
			"album": func(track *Track) {
				track.Album = AlbumModel(c).With("genres", "artists").Find(track.AlbumID)
			},
			"artists": func(track *Track) {
				track.Artists = ArtistModel(c).WhereIn("track_artist", "artist_id", "track_id", track.ID).OrderByRaw("LOWER(`name`)").Get()
			},
		},
	}).Init()
}

// Track artist
type TrackArtist struct {
	ID       string `column:"id,uuid"`
	TrackID  string `column:"track_id,uuid"`
	ArtistID string `column:"artist_id,uuid"`
}

func TrackArtistModel() *database.Model[TrackArtist] {
	return (&database.Model[TrackArtist]{
		TableName: "track_artist",
	}).Init()
}

// Track Like
type TrackLike struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func TrackLikeModel() *database.Model[TrackLike] {
	return (&database.Model[TrackLike]{
		TableName: "track_likes",
	}).Init()
}

// Track Play
type TrackPlay struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	Position  float32   `column:"position,float"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func TrackPlayModel() *database.Model[TrackPlay] {
	return (&database.Model[TrackPlay]{
		TableName: "track_plays",
	}).Init()
}

func HandleTrackPlay(authUser *User, trackID string, position float32) bool {
	// Check if track exists
	track := TrackModel(nil).Find(trackID)
	if track == nil {
		return false
	}

	// Get user last track play and update if latest
	trackPlay := TrackPlayModel().Where("user_id", authUser.ID).OrderByDesc("created_at").First()
	if trackPlay != nil {
		if track.ID == trackPlay.TrackID {
			TrackPlayModel().Where("id", trackPlay.ID).Update(database.Map{
				"position": position,
			})
			return true
		}
	}

	// Create new track play
	TrackPlayModel().Create(database.Map{
		"track_id": track.ID,
		"user_id":  authUser.ID,
		"position": position,
	})

	// Increment global track plays count
	TrackModel(nil).Where("id", track.ID).Update(database.Map{
		"plays": track.Plays + 1,
	})
	return true
}
