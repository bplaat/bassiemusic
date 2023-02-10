package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"

	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/utils"
	uuid "github.com/satori/go.uuid"
)

const TRACK_DURATION_SLACK int = 5

func createArtist(id int, name string) string {
	artists := database.Query("SELECT BIN_TO_UUID(`id`) FROM `artists` WHERE `name` = ?", name)
	defer artists.Close()

	if artists.Next() {
		var artistID string
		_ = artists.Scan(&artistID)
		return artistID
	}

	var artist DeezerArtist
	utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", id), &artist)

	artistID := uuid.NewV4()
	database.Exec("INSERT INTO `artists` (`id`, `name`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?)", artistID.String(), name, id)
	utils.FetchFile(artist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artistID.String()))
	utils.FetchFile(artist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artistID.String()))
	utils.FetchFile(artist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artistID.String()))
	return artistID.String()
}

func createGenre(id int, name string) string {
	genres := database.Query("SELECT BIN_TO_UUID(`id`) FROM `genres` WHERE `name` = ?", name)
	defer genres.Close()

	if genres.Next() {
		var genreID string
		_ = genres.Scan(&genreID)
		return genreID
	}

	var genre DeezerGenre
	utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", id), &genre)

	genreID := uuid.NewV4()
	database.Exec("INSERT INTO `genres` (`id`, `name`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?)", genreID.String(), name, id)
	if genre.PictureMedium != "" {
		utils.FetchFile(genre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genreID.String()))
		utils.FetchFile(genre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genreID.String()))
		utils.FetchFile(genre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genreID.String()))
	} else {
		utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/small/%s.jpg", genreID.String()))
		utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//500x500-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/medium/%s.jpg", genreID.String()))
		utils.FetchFile("https://e-cdns-images.dzcdn.net/images/misc//1000x1000-000000-80-0-0.jpg", fmt.Sprintf("storage/genres/large/%s.jpg", genreID.String()))
	}
	return genreID.String()
}

func downloadAlbum(id int) {
	var album DeezerAlbum
	utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", id), &album)

	// Check if album already exists
	albums := database.Query("SELECT BIN_TO_UUID(`id`) FROM `albums` WHERE `title` = ?", album.Title)
	defer albums.Close()
	if albums.Next() {
		return
	}

	// Create album row
	albumType := models.AlbumTypeAlbum
	if album.RecordType == "ep" {
		albumType = models.AlbumTypeEP
	}
	if album.RecordType == "single" {
		albumType = models.AlbumTypeSingle
	}
	albumID := uuid.NewV4()
	database.Exec("INSERT INTO `albums` (`id`, `type`, `title`, `released_at`, `explicit`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		albumID.String(), albumType, album.Title, album.ReleaseDate, album.ExplicitLyrics, album.ID)
	fmt.Printf("Downloading %s by %s\n", album.Title, album.Artist.Name)

	// Create album genres
	for _, genre := range album.Genres.Data {
		genreID := createGenre(genre.ID, genre.Name)
		albumGenreID := uuid.NewV4()
		database.Exec("INSERT INTO `album_genre` (`id`, `album_id`, `genre_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", albumGenreID.String(), albumID.String(), genreID)
	}

	// Create album artists bindings
	for _, artist := range album.Contributors {
		artistID := createArtist(artist.ID, artist.Name)
		albumArtistID := uuid.NewV4()
		database.Exec("INSERT INTO `album_artist` (`id`, `album_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", albumArtistID.String(), albumID.String(), artistID)
	}

	utils.FetchFile(album.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", albumID.String()))
	utils.FetchFile(album.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", albumID.String()))
	utils.FetchFile(album.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", albumID.String()))

	// Create album tracks
	for _, incompleteTrack := range album.Tracks.Data {
		var track DeezerTrack
		utils.FetchJson(fmt.Sprintf("https://api.deezer.com/track/%d", incompleteTrack.ID), &track)

		// Search for youtube video
		searchCommand := exec.Command("yt-dlp", "--dump-json", fmt.Sprintf("ytsearch25:%s - %s - %s", track.Contributors[0].Name, track.Album.Title, track.Title))
		stdout, err := searchCommand.StdoutPipe()
		if err != nil {
			log.Fatalln(err)
		}
		if err := searchCommand.Start(); err != nil {
			log.Fatalln(err)
		}
		for {
			var video YoutubeVideo
			if err := json.NewDecoder(stdout).Decode(&video); err != nil {
				break
			}

			if track.Duration >= video.Duration-TRACK_DURATION_SLACK && track.Duration <= video.Duration+TRACK_DURATION_SLACK {
				if err := searchCommand.Process.Signal(syscall.SIGTERM); err != nil {
					log.Fatalln(err)
				}

				trackID := uuid.NewV4()
				database.Exec("INSERT INTO `tracks` (`id`, `album_id`, `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?)",
					trackID.String(), albumID.String(), track.Title, track.DiskNumber, track.TrackPosition, video.Duration, track.ExplicitLyrics, track.ID, video.ID)

				// Create track artists bindings
				for _, artist := range track.Contributors {
					artistID := createArtist(artist.ID, artist.Name)
					trackArtistID := uuid.NewV4()
					database.Exec("INSERT INTO `track_artist` (`id`, `track_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", trackArtistID.String(), trackID.String(), artistID)
				}

				downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID),
					"-o", fmt.Sprintf("storage/tracks/%s.m4a", trackID.String()))
				if err := downloadCommand.Start(); err != nil {
					log.Fatalln(err)
				}

				fmt.Printf("%d. %s\n", track.TrackPosition, track.Title)

				break
			}
		}
	}
}

func DownloadTask() {
	for {
		time.Sleep(time.Second)

		// Get first download task
		downloadTask := models.DownloadTaskModel().First()
		if downloadTask == nil {
			continue
		}

		// Do download task
		if downloadTask.Type == "deezer_artist" {
			var artistAlbums DeezerArtistAlbums
			utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d/albums", downloadTask.DeezerID), &artistAlbums)
			for _, album := range artistAlbums.Data {
				if !downloadTask.Singles && album.RecordType == "single" {
					continue
				}
				downloadAlbum(album.ID)
			}
		}
		if downloadTask.Type == "deezer_album" {
			downloadAlbum(int(downloadTask.DeezerID))
		}

		// Delete download task when done
		models.DownloadTaskModel().Where("id", downloadTask.ID).Delete()
	}
}
