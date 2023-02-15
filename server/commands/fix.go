package commands

import (
	"fmt"
	"log"
	"math"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
)

// Try to redownload missing album tracks
func downloadAlbumMissingTracks(album *models.Album, deezerAlbum *structs.DeezerAlbum) {
	log.Printf("Fixing album %s by %s\n", album.Title, album.Artists[0].Name)
	for _, incompleteTrack := range deezerAlbum.Tracks.Data {
		// Get Deezer track info
		var deezerTrack structs.DeezerTrack
		if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/track/%d", incompleteTrack.ID), &deezerTrack); err != nil {
			log.Fatalln(err)
		}

		// When track is missing redownload it
		track := models.TrackModel(nil).Where("album_id", album.ID).Where("disk", deezerTrack.DiskNumber).Where("position", deezerTrack.TrackPosition).First()
		if track == nil {
			tasks.DownloadTrack(album.ID, deezerTrack.ID)
		}
	}
}

func Fix() {
	// Loop over all albums paginated
	totalAlbums := int(models.AlbumModel(nil).Count())
	for page := 1; page <= int(math.Ceil(float64(totalAlbums)/20.0)); page++ {
		albums := models.AlbumModel(nil).With("artists", "tracks").Paginate(page, 20).Data
		for _, album := range albums {
			// Fetch deezer album
			var deezerAlbum structs.DeezerAlbum
			if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", album.DeezerID), &deezerAlbum); err != nil {
				log.Fatalln(err)
			}

			// Check if track counts are different
			if len(album.Tracks) != deezerAlbum.NbTracks {
				downloadAlbumMissingTracks(&album, &deezerAlbum)
			}
		}
	}
}
