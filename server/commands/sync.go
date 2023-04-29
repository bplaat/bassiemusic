package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
)

func Sync() {
	total := models.ArtistModel.Where("sync", true).Count()
	index := 0
	models.ArtistModel.With("albums").Where("sync", true).Chunk(50, func(artists []models.Artist) {
		for _, artist := range artists {
			// Fetch deezer artist albums
			var artistAlbums structs.DeezerArtistAlbums
			if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d/albums", artist.DeezerID), &artistAlbums); err != nil {
				log.Fatalln(err)
			}

			// Download aritst albums when different (albums that already exists will be skipped)
			if len(*artist.Albums) != len(artistAlbums.Data) {
				for _, album := range artistAlbums.Data {
					if strings.Contains(album.Title, "Deezer") {
						continue
					}

					var deezerAlbum structs.DeezerAlbum
					if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/album/%d", album.ID), &deezerAlbum); err != nil {
						log.Fatalln(err)
					}

					tasks.DownloadAlbum(deezerAlbum, nil)
				}
			}

			// Print progress
			log.Printf("Synced artists %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}
