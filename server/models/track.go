package models

import (
	"fmt"
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

func TrackModel(c *fiber.Ctx) database.Model[Track] {
	return database.Model[Track]{
		TableName: "tracks",
		Process: func(track *Track) {
			if c != nil {
				track.Music = fmt.Sprintf("%s/storage/tracks/%s.jpg", c.BaseURL(), track.ID)

				track.Liked = TrackLikeModel().Where("track_id", track.ID).Where("user_id", AuthUser(c).ID).First() != nil
			}
		},
		Relationships: map[string]database.QueryBuilderProcess[Track]{
			"album": func(track *Track) {
				track.Album = AlbumModel(c).Find(track.AlbumID)
			},
			"artists": func(track *Track) {
				track.Artists = ArtistModel(c).WhereIn("track_artist", "artist_id", "track_id", track.ID).OrderByRaw("LOWER(`name`)").Get()
			},
		},
	}.Init()
}

// Track Like
type TrackLike struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func TrackLikeModel() database.Model[TrackLike] {
	return database.Model[TrackLike]{
		TableName: "track_likes",
	}.Init()
}

// Track Play
type TrackPlay struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	Position  float32   `column:"position,float"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

func TrackPlayModel() database.Model[TrackPlay] {
	return database.Model[TrackPlay]{
		TableName: "track_plays",
	}.Init()
}
