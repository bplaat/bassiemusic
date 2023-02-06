package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/tasks"
	"github.com/bplaat/bassiemusic/utils"
)

func restore() {
	// Redownload all artist images
	artistsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `artists`")
	defer artistsQuery.Close()
	artists := models.ArtistsScan(nil, artistsQuery, false, false)
	for index, artist := range artists {
		if _, err := os.Stat(fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID)); !os.IsNotExist(err) {
			continue
		}
		var deezerArtist tasks.DeezerArtist
		utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", artist.DeezerID), &deezerArtist)
		utils.FetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artist.ID))
		utils.FetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artist.ID))
		utils.FetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artist.ID))
		fmt.Printf("Artist images %.2f%%\n", float32(index+1)/float32(len(artists))*100.0)
	}

	// Redownload all genre images
	genresQuery := database.Query("SELECT BIN_TO_UUID(`id`), `name`, `deezer_id`, `created_at` FROM `genres`")
	defer genresQuery.Close()
	genres := models.GenresScan(nil, genresQuery, false)
	for index, genre := range genres {
		if _, err := os.Stat(fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID)); !os.IsNotExist(err) {
			continue
		}
		var deezerGenre tasks.DeezerGenre
		utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", genre.DeezerID), &deezerGenre)
		if deezerGenre.PictureMedium != "" {
			utils.FetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
			utils.FetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
			utils.FetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
		} else {
			utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/small/%s.jpg", genre.ID))
			utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//500x500-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/medium/%s.jpg", genre.ID))
			utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//1000x1000-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/large/%s.jpg", genre.ID))
		}
		fmt.Printf("Genre images %.2f%%\n", float32(index+1)/float32(len(genres))*100.0)
	}

	// Redownload all album covers
	albumsQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `title`, `released_at`, `explicit`, `deezer_id`, `created_at` FROM `albums`")
	defer albumsQuery.Close()
	albums := models.AlbumsScan(nil, albumsQuery, false, false, false)
	for index, album := range albums {
		if _, err := os.Stat(fmt.Sprintf("storage/albums/small/%s.jpg", album.ID)); !os.IsNotExist(err) {
			continue
		}
		var deezerAlbum tasks.DeezerAlbum
		utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", album.DeezerID), &deezerAlbum)
		utils.FetchFile(deezerAlbum.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", album.ID))
		utils.FetchFile(deezerAlbum.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", album.ID))
		utils.FetchFile(deezerAlbum.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", album.ID))
		fmt.Printf("Album covers %.2f%%\n", float32(index+1)/float32(len(albums))*100.0)
	}

	// Redownload all tracks
	tracksQuery := database.Query("SELECT BIN_TO_UUID(`id`), BIN_TO_UUID(`album_id`), `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`, `plays`, `created_at` FROM `tracks`")
	defer tracksQuery.Close()
	tracks := models.TracksScan(nil, tracksQuery, false, false)
	for index, track := range tracks {
		if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", track.ID)); !os.IsNotExist(err) {
			continue
		}
		downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", track.YoutubeID),
			"-o", fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
		if err := downloadCommand.Run(); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Track music %.2f%%\n", float32(index+1)/float32(len(tracks))*100.0)
	}
}
