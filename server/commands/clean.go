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
				if err := os.Remove(fmt.Sprintf("storage/avatars/%s.jpg", avatarID)); err != nil {
					log.Fatalln(err)
				}
				log.Printf("Removed track %s music", avatarID)
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
				if err := os.Remove(fmt.Sprintf("storage/artists/small/%s.jpg", artistID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/artists/medium/%s.jpg", artistID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/artists/large/%s.jpg", artistID)); err != nil {
					log.Fatalln(err)
				}
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
				if err := os.Remove(fmt.Sprintf("storage/albums/small/%s.jpg", albumID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/albums/medium/%s.jpg", albumID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/albums/large/%s.jpg", albumID)); err != nil {
					log.Fatalln(err)
				}
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
				if err := os.Remove(fmt.Sprintf("storage/genres/small/%s.jpg", genreID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/genres/medium/%s.jpg", genreID)); err != nil {
					log.Fatalln(err)
				}
				if err := os.Remove(fmt.Sprintf("storage/genres/large/%s.jpg", genreID)); err != nil {
					log.Fatalln(err)
				}
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
				if err := os.Remove(fmt.Sprintf("storage/tracks/%s.m4a", trackID)); err != nil {
					log.Fatalln(err)
				}
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
	cleanArtistImages()
	cleanAlbumImages()
	cleanGenreImages()
	cleanTrackMusic()
}
