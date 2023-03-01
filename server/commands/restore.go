package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
)

func removeOldUserAvatarIDs() {
	total := models.UserModel().Count()
	index := 0
	models.UserModel().Chunk(50, func(users []models.User) {
		for _, user := range users {
			if _, err := os.Stat(fmt.Sprintf("storage/avatars/%s.jpg", user.Avatar)); os.IsNotExist(err) {
				models.UserModel().Where("id", user.ID).Update(database.Map{
					"avatar": nil,
				})
			}
			log.Printf("User avatars %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoreArtistImages() {
	total := models.ArtistModel(nil).Count()
	index := 0
	models.ArtistModel(nil).Chunk(50, func(artists []models.Artist) {
		for _, artist := range artists {
			if _, err := os.Stat(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID)); os.IsNotExist(err) {
				var deezerArtist structs.DeezerArtist
				if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", artist.DeezerID), &deezerArtist); err != nil {
					log.Fatalln(err)
				}
				utils.FetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID))
				utils.FetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artist.ID))
				utils.FetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artist.ID))
			}
			log.Printf("Artist images %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoreGenreImages() {
	total := models.GenreModel(nil).Count()
	index := 0
	models.GenreModel(nil).Chunk(50, func(genres []models.Genre) {
		for _, genre := range genres {
			if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); os.IsNotExist(err) {
				var deezerGenre structs.DeezerGenre
				if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", genre.DeezerID), &deezerGenre); err != nil {
					log.Fatalln(err)
				}
				if deezerGenre.PictureMedium != "" {
					utils.FetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
					utils.FetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
					utils.FetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
				} else {
					utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
					utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//500x500-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
					utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//1000x1000-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
				}
			}
			log.Printf("Genre images %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoresAlbumCovers() {
	total := models.AlbumModel(nil).Count()
	index := 0
	models.AlbumModel(nil).Chunk(50, func(albums []models.Album) {
		for _, album := range albums {
			if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); os.IsNotExist(err) {
				var deezerAlbum structs.DeezerAlbum
				if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", album.DeezerID), &deezerAlbum); err != nil {
					log.Fatalln(err)
				}
				utils.FetchFile(deezerAlbum.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", album.ID))
				utils.FetchFile(deezerAlbum.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", album.ID))
				utils.FetchFile(deezerAlbum.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", album.ID))
			}
			log.Printf("Album covers %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoreTrackMusic() {
	total := models.TrackModel(nil).WhereNotNull("youtube_id").Count()
	index := 0
	models.TrackModel(nil).WhereNotNull("youtube_id").Chunk(50, func(tracks []models.Track) {
		for _, track := range tracks {
			if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); os.IsNotExist(err) {
				downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", *track.YoutubeID),
					"-o", fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
				if err := downloadCommand.Run(); err != nil {
					log.Fatalln(err)
				}
			}
			log.Printf("Track music %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func Restore() {
	removeOldUserAvatarIDs()
	restoreArtistImages()
	restoreGenreImages()
	restoresAlbumCovers()
	restoreTrackMusic()
}
