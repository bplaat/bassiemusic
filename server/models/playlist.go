package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
)

// Playlist
type Playlist struct {
	ID          string              `column:"id,uuid" json:"id"`
	UserID      string              `column:"user_id,uuid" json:"-"`
	Name        string              `column:"name,string" json:"name"`
	ImageID     database.NullString `column:"image,uuid" json:"-"`
	Public      bool                `column:"public,bool" json:"public"`
	CreatedAt   time.Time           `column:"created_at,timestamp" json:"created_at"`
	UpdatedAt   time.Time           `column:"updated_at,timestamp" json:"-"`
	SmallImage  *string             `json:"small_image"`
	MediumImage *string             `json:"medium_image"`
	Liked       *bool               `json:"liked,omitempty"`
	User        *User               `json:"user,omitempty"`
	Tracks      *[]Track            `json:"tracks,omitempty"`
}

var PlaylistModel *database.Model[Playlist]

func init() {
	PlaylistModel = (&database.Model[Playlist]{
		TableName: "playlists",
		Process: func(playlist *Playlist) {
			if playlist.ImageID.Valid {
				if _, err := os.Stat(fmt.Sprintf("storage/playlists/original/%s", playlist.ImageID.String)); err == nil {
					smallImage := fmt.Sprintf("%s/playlists/small/%s.jpg", os.Getenv("STORAGE_URL"), playlist.ImageID.String)
					playlist.SmallImage = &smallImage
					mediumImage := fmt.Sprintf("%s/playlists/medium/%s.jpg", os.Getenv("STORAGE_URL"), playlist.ImageID.String)
					playlist.MediumImage = &mediumImage
				}
			}
		},
		Relationships: map[string]database.ModelRelationshipFunc[Playlist]{
			"liked": func(playlist *Playlist, args []any) {
				if len(args) > 0 {
					authUser := args[0].(*User)
					liked := PlaylistLikeModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).Count() != 0
					playlist.Liked = &liked
				}
			},
			"user": func(playlist *Playlist, args []any) {
				playlist.User = UserModel.Find(playlist.UserID)
			},
			"tracks": func(playlist *Playlist, args []any) {
				playlistTracks := PlaylistTrackModel.Select("track_id").Where("playlist_id", playlist.ID).OrderBy("position").Get()
				if len(playlistTracks) == 0 {
					emptyTracks := []Track{}
					playlist.Tracks = &emptyTracks
					return
				}
				var trackIds []any
				for _, playlistTrack := range playlistTracks {
					trackIds = append(trackIds, playlistTrack.TrackID)
				}
				tracksQuery := TrackModel.With("artists", "album").WhereIn("id", trackIds)
				if len(args) > 0 {
					authUser := args[0].(*User)
					tracksQuery = tracksQuery.WithArgs("liked", authUser)
				}
				tracks := tracksQuery.Get()

				var orderedTracks []Track
				for _, playlistTrack := range playlistTracks {
					for _, track := range tracks {
						if track.ID == playlistTrack.TrackID {
							orderedTracks = append(orderedTracks, track)
							break
						}
					}
				}
				playlist.Tracks = &orderedTracks
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

var PlaylistTrackModel *database.Model[PlaylistTrack] = (&database.Model[PlaylistTrack]{
	TableName: "playlist_track",
}).Init()

// Playlist Like
type PlaylistLike struct {
	ID         string    `column:"id,uuid"`
	PlaylistID string    `column:"playlist_id,uuid"`
	UserID     string    `column:"user_id,uuid"`
	CreatedAt  time.Time `column:"created_at,timestamp"`
}

var PlaylistLikeModel *database.Model[PlaylistLike] = (&database.Model[PlaylistLike]{
	TableName: "playlist_likes",
}).Init()
