package commands

import (
	"fmt"
	"io"
	"log"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
)

// Create missing album tracks
func createMissingTracks() {
	models.AlbumModel(nil).With("artists", "tracks").Chunk(50, func(albums []models.Album) {
		for _, album := range albums {
			// Fetch deezer album
			var deezerAlbum structs.DeezerAlbum
			if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", album.DeezerID), &deezerAlbum); err != nil {
				log.Fatalln(err)
			}

			// Create missing album tracks
			if len(album.Tracks) != len(deezerAlbum.Tracks.Data) {
				log.Printf("Fix album %s by %s\n", album.Title, album.Artists[0].Name)
				for _, deezerTrack := range deezerAlbum.Tracks.Data {
					track := models.TrackModel(nil).Where("album_id", album.ID).Where("title", deezerTrack.Title).First()
					if track == nil {
						log.Printf("%s\n", deezerTrack.Title)
						tasks.CreateTrack(album.ID, deezerTrack.ID)
					}
				}
			}
		}
	})
}

// Try to search and download youtube videos for all tracks without it
func searchAndDownloadMissingTrackMusic() {
	models.TrackModel(nil).With("album", "artists").WhereNull("youtube_id").Chunk(50, func(tracks []models.Track) {
		for _, track := range tracks {
			log.Printf("Redownloading track %s - %d-%d - %s\n", track.Album.Title, track.Disk, track.Position, track.Title)
			if err := tasks.SearchAndDownloadTrackMusic(&track); err != nil && err != io.EOF {
				log.Fatalln(err)
			}
		}
	})
}

func Fix() {
	createMissingTracks()
	searchAndDownloadMissingTrackMusic()
}
