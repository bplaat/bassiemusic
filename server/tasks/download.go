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

func createArtist(id int, name string, image string) string {
	artists := database.Query("SELECT BIN_TO_UUID(`id`) FROM `artists` WHERE `name` = ?", name)
	defer artists.Close()

	if artists.Next() {
		var artistId string
		artists.Scan(&artistId)
		return artistId
	}

	artistId := uuid.NewV4()
	database.Exec("INSERT INTO `artists` (`id`, `name`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?)", artistId.String(), name, id)
	utils.FetchFile(image, fmt.Sprintf("storage/artists/%s.jpg", artistId.String()))
	return artistId.String()
}

func createGenre(id int, name string) string {
	genres := database.Query("SELECT BIN_TO_UUID(`id`) FROM `genres` WHERE `name` = ?", name)
	defer genres.Close()

	if genres.Next() {
		var genreId string
		genres.Scan(&genreId)
		return genreId
	}

	var genre DeezerGenre
	utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", id), &genre)

	genreId := uuid.NewV4()
	database.Exec("INSERT INTO `genres` (`id`, `name`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?)", genreId.String(), name, id)
	utils.FetchFile(genre.PictureXl, fmt.Sprintf("storage/genres/%s.jpg", genreId.String()))
	return genreId.String()
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
	albumId := uuid.NewV4()
	database.Exec("INSERT INTO `albums` (`id`, `type`, `title`, `released_at`, `explicit`, `deezer_id`) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		albumId.String(), albumType, album.Title, album.ReleaseDate, album.ExplicitLyrics, album.ID)
	fmt.Printf("Downloading %s by %s\n", album.Title, album.Artist.Name)

	// Create album genres
	for _, genre := range album.Genres.Data {
		genreId := createGenre(genre.ID, genre.Name)
		albumGenreId := uuid.NewV4()
		database.Exec("INSERT INTO `album_genre` (`id`, `album_id`, `genre_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", albumGenreId.String(), albumId.String(), genreId)
	}

	// Create album artists bindings
	for _, artist := range album.Contributors {
		artistId := createArtist(artist.ID, artist.Name, artist.PictureXl)
		albumArtistId := uuid.NewV4()
		database.Exec("INSERT INTO `album_artist` (`id`, `album_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", albumArtistId.String(), albumId.String(), artistId)
	}

	utils.FetchFile(album.CoverXl, fmt.Sprintf("storage/albums/%s.jpg", albumId.String()))

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
				searchCommand.Process.Signal(syscall.SIGTERM)

				trackId := uuid.NewV4()
				database.Exec("INSERT INTO `tracks` (`id`, `album_id`, `title`, `disk`, `position`, `duration`, `explicit`, `deezer_id`, `youtube_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?)",
					trackId.String(), albumId.String(), track.Title, track.DiskNumber, track.TrackPosition, video.Duration, track.ExplicitLyrics, track.ID, video.ID)

				// Create track artists bindings
				for _, artist := range track.Contributors {
					artistId := createArtist(artist.ID, artist.Name, artist.PictureXl)
					trackArtistId := uuid.NewV4()
					database.Exec("INSERT INTO `track_artist` (`id`, `track_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", trackArtistId.String(), trackId.String(), artistId)
				}

				downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID),
					"-o", fmt.Sprintf("storage/tracks/%s.m4a", trackId.String()))
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

		downloadTaskQuery := database.Query("SELECT BIN_TO_UUID(`id`), `type`, `deezer_id`, `singles`, `created_at` FROM `download_tasks` ORDER BY `created_at` LIMIT 2")
		defer downloadTaskQuery.Close()

		if !downloadTaskQuery.Next() {
			continue
		}

		task := models.DownloadTaskScan(downloadTaskQuery)

		if task.Type == "deezer_artist" {
			var artistAlbums DeezerArtistAlbums
			utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d/albums", task.DeezerID), &artistAlbums)
			for _, album := range artistAlbums.Data {
				if !task.Singles && album.RecordType == "single" {
					continue
				}
				downloadAlbum(album.ID)
			}
		}

		if task.Type == "deezer_album" {
			downloadAlbum(int(task.DeezerID))
		}

		// Delete download task when done
		database.Exec("DELETE FROM `download_tasks` WHERE `id` = UUID_TO_BIN(?)", task.ID)
	}
}
