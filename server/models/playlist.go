package models

import (
	"fmt"
	"os"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

// Playlist
type Playlist struct {
	ID          uuid.Uuid           `column:"id" json:"id"`
	Name        string              `column:"name" json:"name"`
	ImageID     uuid.NullUuid       `column:"image" json:"-"`
	Public      bool                `column:"public" json:"public"`
	CreatedAt   time.Time           `column:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `column:"updated_at" json:"-"`
	SmallImage  *string             `json:"small_image"`
	MediumImage *string             `json:"medium_image"`
	Liked       *bool               `json:"liked,omitempty"`
	Owners      *[]User             `json:"owners,omitempty"`
	Users       *[]PlaylistUserItem `json:"users,omitempty"`
	Tracks      *[]Track            `json:"tracks,omitempty"`
}

type PlaylistUserItem struct {
	User User   `json:"user"`
	Role string `json:"role"`
}

var PlaylistModel *database.Model[Playlist]

func init() {
	PlaylistModel = (&database.Model[Playlist]{
		TableName: "playlists",
		Process: func(playlist *Playlist) {
			if playlist.ImageID.Valid {
				if _, err := os.Stat(fmt.Sprintf("storage/playlists/original/%s", playlist.ImageID.Uuid)); err == nil {
					smallImage := fmt.Sprintf("%s/playlists/small/%s.jpg", os.Getenv("STORAGE_URL"), playlist.ImageID.Uuid)
					playlist.SmallImage = &smallImage
					mediumImage := fmt.Sprintf("%s/playlists/medium/%s.jpg", os.Getenv("STORAGE_URL"), playlist.ImageID.Uuid)
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
			"owners": func(playlist *Playlist, args []any) {
				playlistOwners := PlaylistUserModel.Select("user_id", "role").Where("playlist_id", playlist.ID).Where("role", PlaylistUserRoleOwner).OrderBy("created_at").Get()
				var ownerIds []any
				for _, playlistOwner := range playlistOwners {
					ownerIds = append(ownerIds, playlistOwner.UserID)
				}
				owners := UserModel.WhereIn("id", ownerIds).Get()
				playlist.Owners = &owners
			},
			"users": func(playlist *Playlist, args []any) {
				playlistUsers := PlaylistUserModel.Select("user_id", "role").Where("playlist_id", playlist.ID).OrderByRaw("`role`, `created_at`").Get()
				var userIds []any
				for _, playlistUser := range playlistUsers {
					userIds = append(userIds, playlistUser.UserID)
				}
				users := UserModel.WhereIn("id", userIds).Get()

				playlistUserItems := []PlaylistUserItem{}
				for index, user := range users {
					playlistUserItems = append(playlistUserItems, PlaylistUserItem{User: user, Role: playlistUsers[index].RoleString})
				}
				playlist.Users = &playlistUserItems
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

// Playlist User
type PlaylistUser struct {
	ID         uuid.Uuid        `column:"id"`
	PlaylistID uuid.Uuid        `column:"playlist_id"`
	UserID     uuid.Uuid        `column:"user_id"`
	Role       PlaylistUserRole `column:"role" json:"-"`
	RoleString string           `json:"role"`
	CreatedAt  time.Time        `column:"created_at"`
}

type PlaylistUserRole int

const PlaylistUserRoleViewer PlaylistUserRole = 0
const PlaylistUserRoleEditor PlaylistUserRole = 1
const PlaylistUserRoleOwner PlaylistUserRole = 2

var PlaylistUserModel *database.Model[PlaylistUser] = (&database.Model[PlaylistUser]{
	TableName: "playlist_user",
	Process: func(playlistUser *PlaylistUser) {
		if playlistUser.Role == PlaylistUserRoleViewer {
			playlistUser.RoleString = "viewer"
		}
		if playlistUser.Role == PlaylistUserRoleEditor {
			playlistUser.RoleString = "editor"
		}
		if playlistUser.Role == PlaylistUserRoleOwner {
			playlistUser.RoleString = "owner"
		}
	},
}).Init()

// Playlist Track
type PlaylistTrack struct {
	ID         uuid.Uuid `column:"id"`
	PlaylistID uuid.Uuid `column:"playlist_id"`
	TrackID    uuid.Uuid `column:"track_id"`
	Position   int       `column:"position"`
	CreatedAt  time.Time `column:"created_at"`
}

var PlaylistTrackModel *database.Model[PlaylistTrack] = (&database.Model[PlaylistTrack]{
	TableName: "playlist_track",
}).Init()

// Playlist Like
type PlaylistLike struct {
	ID         uuid.Uuid `column:"id"`
	PlaylistID uuid.Uuid `column:"playlist_id"`
	UserID     uuid.Uuid `column:"user_id"`
	CreatedAt  time.Time `column:"created_at"`
}

var PlaylistLikeModel *database.Model[PlaylistLike] = (&database.Model[PlaylistLike]{
	TableName: "playlist_likes",
}).Init()
