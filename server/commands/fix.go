package commands

import (
	"fmt"
	"log"
	"math"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
)

// Try to redownload missing album tracks
func fixAlbum(album models.Album) {
	log.Printf("Fix album %s by %s\n", album.Title, album.Artists[0].Name)
	// TODO
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
				fixAlbum(album)
			}
		}
	}
}
