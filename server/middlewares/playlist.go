package middlewares

import (
	"github.com/bplaat/bassiemusic/models"
)

func IsPlaylistViewer(playlist *models.Playlist, authUser *models.User) bool {
	if playlist.Public {
		return true
	}

	if authUser.Role == models.UserRoleAdmin {
		return true
	}

	if models.PlaylistUserModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).Count() > 0 {
		return true
	}

	return false
}

func IsPlaylistEditor(playlist *models.Playlist, authUser *models.User) bool {
	if authUser.Role == models.UserRoleAdmin {
		return true
	}

	playlistUser := models.PlaylistUserModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).First()
	if playlistUser != nil && (playlistUser.Role == models.PlaylistUserRoleEditor || playlistUser.Role == models.PlaylistUserRoleOwner) {
		return true
	}

	return false
}

func IsPlaylistOwner(playlist *models.Playlist, authUser *models.User) bool {
	if authUser.Role == models.UserRoleAdmin {
		return true
	}

	playlistUser := models.PlaylistUserModel.Where("playlist_id", playlist.ID).Where("user_id", authUser.ID).First()
	if playlistUser != nil && playlistUser.Role == models.PlaylistUserRoleOwner {
		return true
	}

	return false
}
