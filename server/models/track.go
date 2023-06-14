package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
)

// Track
type Track struct {
	ID        string              `column:"id,uuid" json:"id"`
	AlbumID   string              `column:"album_id,uuid" json:"-"`
	Title     string              `column:"title,string" json:"title"`
	Disk      int                 `column:"disk,int" json:"disk"`
	Position  int                 `column:"position,int" json:"position"`
	Duration  float32             `column:"duration,float" json:"duration"`
	Explicit  bool                `column:"explicit,bool" json:"explicit"`
	DeezerID  int64               `column:"deezer_id,bigint" json:"deezer_id"`
	YoutubeID database.NullString `column:"youtube_id,string" json:"youtube_id"`
	Plays     int64               `column:"plays,bigint" json:"plays"`
	CreatedAt time.Time           `column:"created_at,timestamp" json:"created_at"`
	Music     *string             `json:"music"`
	Liked     *bool               `json:"liked,omitempty"`
	Album     *Album              `json:"album,omitempty"`
	Artists   *[]Artist           `json:"artists,omitempty"`
}

var TrackModel *database.Model[Track]

func init() {
	TrackModel = (&database.Model[Track]{
		TableName: "tracks",
		Process: func(track *Track) {
			if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); err == nil {
				music := fmt.Sprintf("%s/tracks/%s.m4a", os.Getenv("STORAGE_URL"), track.ID)
				track.Music = &music
			}
		},
		Relationships: map[string]database.ModelRelationshipFunc[Track]{
			"liked": func(track *Track, args []any) {
				if len(args) > 0 {
					authUser := args[0].(*User)
					liked := TrackLikeModel.Where("track_id", track.ID).Where("user_id", authUser.ID).Count() != 0
					track.Liked = &liked
				}
			},
			"liked_true": func(track *Track, args []any) {
				liked := true
				track.Liked = &liked
			},
			"album": func(track *Track, args []any) {
				track.Album = AlbumModel.With("genres", "artists").Find(track.AlbumID)
			},
			"artists": func(track *Track, args []any) {
				trackArtists := TrackArtistModel.Select("artist_id").Where("track_id", track.ID).OrderBy("position").Get()
				if len(trackArtists) == 0 {
					emptyArtists := []Artist{}
					track.Artists = &emptyArtists
					return
				}
				var artistIds []any
				for _, trackArtist := range trackArtists {
					artistIds = append(artistIds, trackArtist.ArtistID)
				}
				artists := ArtistModel.WhereIn("id", artistIds).Get()

				var orderedArtists []Artist
				for _, trackArtist := range trackArtists {
					for _, artist := range artists {
						if artist.ID == trackArtist.ArtistID {
							orderedArtists = append(orderedArtists, artist)
							break
						}
					}
				}
				track.Artists = &orderedArtists
			},
		},
	}).Init()
}

// Track artist
type TrackArtist struct {
	ID       string `column:"id,uuid"`
	TrackID  string `column:"track_id,uuid"`
	ArtistID string `column:"artist_id,uuid"`
	Position int    `column:"position,int"`
}

var TrackArtistModel *database.Model[TrackArtist] = (&database.Model[TrackArtist]{
	TableName: "track_artist",
}).Init()

// Track Like
type TrackLike struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

var TrackLikeModel *database.Model[TrackLike] = (&database.Model[TrackLike]{
	TableName: "track_likes",
}).Init()

// Track Play
type TrackPlay struct {
	ID        string    `column:"id,uuid"`
	TrackID   string    `column:"track_id,uuid"`
	UserID    string    `column:"user_id,uuid"`
	Position  float32   `column:"position,float"`
	CreatedAt time.Time `column:"created_at,timestamp"`
}

var TrackPlayModel *database.Model[TrackPlay] = (&database.Model[TrackPlay]{
	TableName: "track_plays",
}).Init()

func HandleTrackPlay(authUser *User, trackID string, position float32) bool {
	// Check if track exists
	track := TrackModel.Find(trackID)
	if track == nil {
		return false
	}

	// Get user last track play and update if latest
	trackPlay := TrackPlayModel.Where("user_id", authUser.ID).OrderByDesc("created_at").First()
	if trackPlay != nil {
		if track.ID == trackPlay.TrackID {
			TrackPlayModel.Where("id", trackPlay.ID).Update(database.Map{
				"position": position,
			})
			return true
		}
	}

	// Create new track play
	TrackPlayModel.Create(database.Map{
		"track_id": track.ID,
		"user_id":  authUser.ID,
		"position": position,
	})

	// Increment global track plays count
	TrackModel.Where("id", track.ID).Update(database.Map{
		"plays": track.Plays + 1,
	})
	return true
}
