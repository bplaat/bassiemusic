package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bplaat/bassiemusic/models"
)

// Clean up all unused user avatars
func cleanUserAvatars() {
	if err := filepath.Walk("storage/avatars", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			avatarID := parts[0]
			if models.UserModel().Where("avatar_id", avatarID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/avatars/original/%s", avatarID))
				_ = os.Remove(fmt.Sprintf("storage/avatars/small/%s.jpg", avatarID))
				_ = os.Remove(fmt.Sprintf("storage/avatars/medium/%s.jpg", avatarID))
				log.Printf("Removed track %s music", avatarID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

// Clean up all unused playlist images
func cleanPlaylistImages() {
	if err := filepath.Walk("storage/playlists", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			imageID := parts[0]
			if models.PlaylistModel(nil).Where("image_id", imageID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/playlists/original/%s", imageID))
				_ = os.Remove(fmt.Sprintf("storage/playlists/small/%s.jpg", imageID))
				_ = os.Remove(fmt.Sprintf("storage/playlists/medium/%s.jpg", imageID))
				log.Printf("Removed track %s music", imageID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

// Clean up all unused artists images
func cleanArtistImages() {
	if err := filepath.Walk("storage/artists/small", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			artistID := parts[0]
			if models.ArtistModel(nil).Find(artistID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/artists/small/%s.jpg", artistID))
				_ = os.Remove(fmt.Sprintf("storage/artists/medium/%s.jpg", artistID))
				_ = os.Remove(fmt.Sprintf("storage/artists/large/%s.jpg", artistID))
				log.Printf("Removed artist %s images", artistID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

// Clean up all unused albums images
func cleanAlbumImages() {
	if err := filepath.Walk("storage/albums/small", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			albumID := parts[0]
			if models.AlbumModel(nil).Find(albumID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/albums/small/%s.jpg", albumID))
				_ = os.Remove(fmt.Sprintf("storage/albums/medium/%s.jpg", albumID))
				_ = os.Remove(fmt.Sprintf("storage/albums/large/%s.jpg", albumID))
				log.Printf("Removed album %s images", albumID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

// Clean up all unused genres images
func cleanGenreImages() {
	if err := filepath.Walk("storage/genres/small", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			genreID := parts[0]
			if models.GenreModel(nil).Find(genreID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/genres/small/%s.jpg", genreID))
				_ = os.Remove(fmt.Sprintf("storage/genres/medium/%s.jpg", genreID))
				_ = os.Remove(fmt.Sprintf("storage/genres/large/%s.jpg", genreID))
				log.Printf("Removed genre %s images", genreID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

// Clean up all unused tracks music
func cleanTrackMusic() {
	if err := filepath.Walk("storage/tracks", func(path string, _ os.FileInfo, err error) error {
		parts := strings.Split(filepath.Base(path), ".")
		if len(parts) == 2 {
			trackID := parts[0]
			if models.TrackModel(nil).WhereNotNull("youtube_id").Find(trackID) == nil {
				_ = os.Remove(fmt.Sprintf("storage/tracks/%s.m4a", trackID))
				log.Printf("Removed track %s music", trackID)
			}
		}
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

func Clean() {
	cleanUserAvatars()
	cleanPlaylistImages()
	cleanArtistImages()
	cleanAlbumImages()
	cleanGenreImages()
	cleanTrackMusic()
}
