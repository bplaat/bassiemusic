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
	total := models.ArtistModel(nil).Where("sync", true).Count()
	index := 0
	models.ArtistModel(nil).Where("sync", true).With("albums").Chunk(50, func(artists []models.Artist) {
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
					tasks.DownloadAlbum(album.ID)
				}
			}

			// Print progress
			log.Printf("Synced artists %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}