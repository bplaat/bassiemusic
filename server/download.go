package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"syscall"

	"github.com/satori/go.uuid"
)

const TRACK_DURATION_SLACK int = 5

func createArtist(name string) {
	artists, err := db.Query("SELECT * FROM `artists` WHERE `name` = ?", name)
	if err != nil {
		log.Fatalln(err)
	}
	defer artists.Close()
	if !artists.Next() {
		var artistSearch DeezerArtistSearch
		fetchJson(fmt.Sprintf("https://api.deezer.com/search/artist?q=%s", url.QueryEscape(name)), &artistSearch)

		artistId := uuid.NewV4()
		db.Exec("INSERT INTO `artists` (`id`, `name`) VALUES (UUID_TO_BIN(?), ?)", artistId.String(), name)

		fetchFile(artistSearch.Data[0].PictureXl, fmt.Sprintf("storage/artists/%s.jpg", artistId.String()))
	}
}

func downloadAlbum(id int) {
	var album DeezerAlbum
	fetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", id), &album)

	// Create album row
	albumId := uuid.NewV4()
	db.Exec("INSERT INTO `albums` (`id`, `title`, `released_at`) VALUES (UUID_TO_BIN(?), ?, ?)",
		albumId.String(), album.Title, album.ReleaseDate)
	fmt.Printf("%s by %s\n", album.Title, album.Artist.Name)

	// Create album artists bindings
	for _, artist := range album.Contributors {
		createArtist(artist.Name)

		artists, err := db.Query("SELECT BIN_TO_UUID(`id`) `id` FROM `artists` WHERE `name` = ?", artist.Name)
		if err != nil {
			log.Fatalln(err)
		}
		defer artists.Close()
		if artists.Next() {
			var artistId string
			err := artists.Scan(&artistId)
			if err != nil {
				log.Fatalln(err)
			}

			albumArtistId := uuid.NewV4()
			db.Exec("INSERT INTO `album_artist` (`id`, `album_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", albumArtistId.String(), albumId.String(), artistId)
		}
	}

	fetchFile(album.CoverXl, fmt.Sprintf("storage/albums/%s.jpg", albumId.String()))

	// Create album tracks
	for _, incompleteTrack := range album.Tracks.Data {
		var track DeezerTrack
		fetchJson(fmt.Sprintf("https://api.deezer.com/track/%d", incompleteTrack.ID), &track)

		trackId := uuid.NewV4()
		db.Exec("INSERT INTO `tracks` (`id`, `album_id`, `title`, `disk`, `position`, `duration`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?)",
			trackId.String(), albumId.String(), track.Title, track.DiskNumber, track.TrackPosition, track.Duration)

		// Create track artists bindings
		for _, artist := range track.Contributors {
			createArtist(artist.Name)

			artists, err := db.Query("SELECT BIN_TO_UUID(`id`) `id` FROM `artists` WHERE `name` = ?", artist.Name)
			if err != nil {
				log.Fatalln(err)
			}
			defer artists.Close()
			if artists.Next() {
				var artistId string
				err := artists.Scan(&artistId)
				if err != nil {
					log.Fatalln(err)
				}

				trackArtistId := uuid.NewV4()
				db.Exec("INSERT INTO `track_artist` (`id`, `track_id`, `artist_id`) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?))", trackArtistId.String(), trackId.String(), artistId)
			}
		}

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

func downloadTracks() {
	var albumSearch DeezerAlbumSearch
	fetchJson(fmt.Sprintf("https://api.deezer.com/search/album?q=%s", url.QueryEscape(os.Args[2])), &albumSearch)

	// Create missing folders
	if _, err := os.Stat("storage/"); os.IsNotExist(err) {
		os.Mkdir("storage", 0755)
	}
	if _, err := os.Stat("storage/artists"); os.IsNotExist(err) {
		os.Mkdir("storage/artists", 0755)
	}
	if _, err := os.Stat("storage/albums"); os.IsNotExist(err) {
		os.Mkdir("storage/albums", 0755)
	}
	if _, err := os.Stat("storage/tracks"); os.IsNotExist(err) {
		os.Mkdir("storage/tracks", 0755)
	}

	if len(albumSearch.Data) == 0 {
		fmt.Println("Can't find any album with that title")
		return
	}

	downloadAlbum(albumSearch.Data[0].ID)
}
