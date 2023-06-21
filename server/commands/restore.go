package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
)

func removeOldUserAvatarIDs() {
	total := models.UserModel.WhereNotNull("avatar").Count()
	index := 0
	models.UserModel.WhereNotNull("avatar").Chunk(50, func(users []models.User) {
		for _, user := range users {
			if user.AvatarID.Valid {
				if _, err := os.Stat(fmt.Sprintf("storage/avatars/original/%s", user.AvatarID.Uuid)); os.IsNotExist(err) {
					models.UserModel.Where("id", user.ID).Update(database.Map{
						"avatar": nil,
					})
				}
			}
			log.Printf("User avatars %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func removeOldPlaylistImageIDs() {
	total := models.PlaylistModel.WhereNotNull("image").Count()
	index := 0
	models.PlaylistModel.WhereNotNull("image").Chunk(50, func(playlists []models.Playlist) {
		for _, playlist := range playlists {
			if playlist.ImageID.Valid {
				if _, err := os.Stat(fmt.Sprintf("storage/playlists/original/%s", playlist.ImageID.Uuid)); os.IsNotExist(err) {
					models.PlaylistModel.Where("id", playlist.ID).Update(database.Map{
						"image": nil,
					})
				}
			}
			log.Printf("Playlist images %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoreArtistImages() {
	total := models.ArtistModel.Count()
	index := 0
	models.ArtistModel.Chunk(50, func(artists []models.Artist) {
		for _, artist := range artists {
			if _, err := os.Stat(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID)); os.IsNotExist(err) {
				var deezerArtist structs.DeezerArtist
				if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", artist.DeezerID), &deezerArtist); err != nil {
					log.Fatalln(err)
				}
				if deezerArtist.PictureMedium != "https://e-cdns-images.dzcdn.net/images/artist//250x250-000000-80-0-0.jpg" {
					utils.FetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID))
					utils.FetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artist.ID))
					utils.FetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artist.ID))
				}
			}
			log.Printf("Artist images %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoreGenreImages() {
	total := models.GenreModel.Count()
	index := 0
	models.GenreModel.Chunk(50, func(genres []models.Genre) {
		for _, genre := range genres {
			if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); os.IsNotExist(err) {
				var deezerGenre structs.DeezerGenre
				if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", genre.DeezerID), &deezerGenre); err != nil {
					log.Fatalln(err)
				}
				if deezerGenre.PictureMedium != "https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg" {
					utils.FetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
					utils.FetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
					utils.FetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
				}
			}
			log.Printf("Genre images %.2f%%\n", float32(index+1)/float32(total)*100.0)
			index++
		}
	})
}

func restoresAlbumCovers() {
	total := models.AlbumModel.Count()
	index := 0
	models.AlbumModel.Chunk(50, func(albums []models.Album) {
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
	total := models.TrackModel.WhereNotNull("youtube_id").Count()
	index := 0
	models.TrackModel.WhereNotNull("youtube_id").Chunk(50, func(tracks []models.Track) {
		for _, track := range tracks {
			if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); os.IsNotExist(err) {
				downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]",
					fmt.Sprintf("https://www.youtube.com/watch?v=%s", track.YoutubeID.String),
					"-o", fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
				log.Println(downloadCommand.String())
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
	removeOldPlaylistImageIDs()
	restoreArtistImages()
	restoreGenreImages()
	restoresAlbumCovers()
	restoreTrackMusic()
}
