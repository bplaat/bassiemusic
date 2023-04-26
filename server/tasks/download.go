package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/bplaat/bassiemusic/consts"
	"github.com/bplaat/bassiemusic/controllers/websocket"
	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
)

func createArtist(deezerID int, name string, sync bool) string {
	// Check if artist already exists
	artist := models.ArtistModel.Where("name", name).First()
	if artist != nil {
		if sync {
			models.ArtistModel.Where("id", artist.ID).Update(database.Map{
				"sync": true,
			})
		}
		return artist.ID
	}

	// Get Deezer artist info
	var deezerArtist structs.DeezerArtist
	if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/artist/%d", deezerID), &deezerArtist); err != nil {
		log.Fatalln(err)
	}

	// Create artist
	artistID := uuid.New()
	models.ArtistModel.Create(database.Map{
		"id":        artistID.String(),
		"name":      name,
		"deezer_id": deezerID,
		"sync":      sync,
	})

	if deezerArtist.PictureMedium != "https://e-cdns-images.dzcdn.net/images/artist//250x250-000000-80-0-0.jpg" {
		utils.DeezerFetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artistID.String()))
		utils.DeezerFetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artistID.String()))
		utils.DeezerFetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artistID.String()))
	}
	return artistID.String()
}

func createGenre(deezerID int, name string) string {
	// Check if genre already exists
	genre := models.GenreModel.Where("name", name).First()
	if genre != nil {
		return genre.ID
	}

	// Get Deezer genre info
	var deezerGenre structs.DeezerGenre
	if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/genre/%d", deezerID), &deezerGenre); err != nil {
		log.Fatalln(err)
	}

	// Create genre
	genreID := uuid.New()
	models.GenreModel.Create(database.Map{
		"id":        genreID.String(),
		"name":      name,
		"deezer_id": deezerID,
	})
	if deezerGenre.PictureMedium != "https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg" {
		utils.DeezerFetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genreID.String()))
		utils.DeezerFetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genreID.String()))
		utils.DeezerFetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genreID.String()))
	}
	return genreID.String()
}

func CreateTrack(albumID string, deezerID int64) {
	// Get Deezer track info
	var deezerTrack structs.DeezerTrack
	if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/track/%d", deezerID), &deezerTrack); err != nil {
		log.Fatalln(err)
	}

	// Create track
	trackID := uuid.New()
	models.TrackModel.Create(database.Map{
		"id":         trackID.String(),
		"album_id":   albumID,
		"title":      deezerTrack.Title,
		"disk":       deezerTrack.DiskNumber,
		"position":   deezerTrack.TrackPosition,
		"duration":   deezerTrack.Duration,
		"explicit":   deezerTrack.ExplicitLyrics,
		"deezer_id":  deezerTrack.ID,
		"youtube_id": nil,
		"plays":      0,
	})

	// Create track artists bindings
	for _, artist := range deezerTrack.Contributors {
		artistID := createArtist(artist.ID, artist.Name, false)
		models.TrackArtistModel.Create(database.Map{
			"track_id":  trackID.String(),
			"artist_id": artistID,
		})
	}
}

func SearchAndDownloadTrackMusic(track *models.Track) error {
	// Search for youtube video
	var searchQuery string
	if len(*track.Artists) > 0 {
		searchQuery = fmt.Sprintf("%s - %s", (*track.Artists)[0].Name, track.Title)
	} else {
		searchQuery = fmt.Sprintf("%s - %s", track.Album.Title, track.Title)
	}
	searchCommand := exec.Command("yt-dlp", "--dump-json", "ytsearch10:"+searchQuery)

	log.Println(searchCommand.String())

	stdout, err := searchCommand.StdoutPipe()

	if err != nil {
		return err
	}

	if err := searchCommand.Start(); err != nil {
		return err
	}

	fallPoins := 0

	lowestScore := 1000000.0
	youtubeID := ""
	youtubeDuration := 0

	for {
		var youtubeVideo structs.YoutubeVideo
		if err := json.NewDecoder(stdout).Decode(&youtubeVideo); err != nil {
			break
		}

		score := float32(youtubeVideo.Duration) - track.Duration
		if score < 0 {
			score = score - score - score
		}
		if score+float32(fallPoins) < float32(lowestScore) {
			lowestScore = float64(score) + float64(fallPoins)
			youtubeID = youtubeVideo.ID
			youtubeDuration = youtubeVideo.Duration
		}

		fallPoins += consts.PUNISHMENT_POINTS
	}

	if youtubeID != "" {
		// Download right youtube video
		downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", youtubeID),
			"-o", fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
		log.Println(downloadCommand.String())

		if err := downloadCommand.Run(); err != nil {
			log.Fatalln(err)
		}

		// Update track
		models.TrackModel.Where("id", track.ID).Update(database.Map{
			"duration":   youtubeDuration,
			"youtube_id": youtubeID,
		})
	}
	return nil
}

func DownloadAlbum(deezerID int) {
	// Get Deezer album info
	var deezerAlbum structs.DeezerAlbum
	if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/album/%d", deezerID), &deezerAlbum); err != nil {
		log.Fatalln(err)
	}

	// Check if album already exists
	if models.AlbumModel.Where("title", deezerAlbum.Title).First() != nil {
		return
	}

	// Create album
	log.Printf("[DOWNLOAD] Start downloading album %s by %s\n", deezerAlbum.Title, deezerAlbum.Artist.Name)
	albumType := models.AlbumTypeAlbum
	if deezerAlbum.RecordType == "ep" {
		albumType = models.AlbumTypeEP
	}
	if deezerAlbum.RecordType == "single" {
		albumType = models.AlbumTypeSingle
	}
	albumID := uuid.New()
	models.AlbumModel.Create(database.Map{
		"id":          albumID.String(),
		"type":        albumType,
		"title":       deezerAlbum.Title,
		"released_at": deezerAlbum.ReleaseDate,
		"explicit":    deezerAlbum.ExplicitLyrics,
		"deezer_id":   deezerID,
	})

	utils.DeezerFetchFile(deezerAlbum.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", albumID.String()))
	utils.DeezerFetchFile(deezerAlbum.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", albumID.String()))
	utils.DeezerFetchFile(deezerAlbum.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", albumID.String()))

	// Create album genre bindings
	for _, genre := range deezerAlbum.Genres.Data {
		genreID := createGenre(genre.ID, genre.Name)
		models.AlbumGenreModel.Create(database.Map{
			"album_id": albumID.String(),
			"genre_id": genreID,
		})
	}

	// Create album artist bindings
	for _, artist := range deezerAlbum.Contributors {
		artistID := createArtist(artist.ID, artist.Name, false)
		models.AlbumArtistModel.Create(database.Map{
			"album_id":  albumID.String(),
			"artist_id": artistID,
		})
	}

	// Create album tracks
	for _, incompleteTrack := range deezerAlbum.Tracks.Data {
		CreateTrack(albumID.String(), incompleteTrack.ID)
	}

	// Download album tracks music
	for _, deezerTrack := range deezerAlbum.Tracks.Data {
		track := models.TrackModel.With("album", "artists").Where("album_id", albumID.String()).Where("title", deezerTrack.Title).First()
		if err := SearchAndDownloadTrackMusic(track); err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		log.Printf("[DOWNLOAD] %s - %d-%d - %s\n", deezerAlbum.Title, track.Disk, track.Position, track.Title)
	}

	log.Printf("[DOWNLOAD] Done downloading album\n")
}

func fetchAlbums(DeezerID int64) []structs.DeezerArtistAlbum {
	nextUrl := fmt.Sprintf("https://api.deezer.com/artist/%d/albums", DeezerID)
	var albums []structs.DeezerArtistAlbum
	for {
		var artistAlbums structs.DeezerArtistAlbums
		if err := utils.DeezerFetch(nextUrl, &artistAlbums); err != nil {
			log.Fatalln(err)
		}
		albums = append(albums, artistAlbums.Data...)
		if artistAlbums.Next == "" {
			break
		} else {
			nextUrl = artistAlbums.Next
		}
	}
	return albums
}

func DownloadTask() {
	for {
		// Wait a little while
		time.Sleep(5 * time.Second)

		// Get oldest download task
		downloadTask := models.DownloadTaskModel.OrderBy("created_at").First()
		if downloadTask == nil {
			continue
		}

		//  Execute current download task
		if downloadTask.Type == models.DownloadTaskTypeDeezerArtist {
			// Update download task status
			downloadTask.Status = models.DownloadTaskStatusDownloading
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}

			// Create artist
			var deezerArtist structs.DeezerArtist
			if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/artist/%d", downloadTask.DeezerID), &deezerArtist); err != nil {
				log.Fatalln(err)
			}
			createArtist(deezerArtist.ID, deezerArtist.Name, true)

			// Download artist albums
			artistAlbums := fetchAlbums(downloadTask.DeezerID)
			for index, album := range artistAlbums {
				if strings.Contains(album.Title, "Deezer") {
					continue
				}
				DownloadAlbum(album.ID)

				// Update download task progress
				downloadTask.Progress = float32(math.Round((float64(index) / float64(len(artistAlbums))) * 100))
				models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
					"progress": downloadTask.Progress,
				})

				// Broadcast download tasks update message to admins
				if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
					log.Println(err)
				}
			}
		}

		if downloadTask.Type == models.DownloadTaskTypeDeezerAlbum {
			// Update download task status
			downloadTask.Status = models.DownloadTaskStatusDownloading
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}

			// Download album
			DownloadAlbum(int(downloadTask.DeezerID))
		}

		// Broadcast download tasks delete message to admins
		if err := websocket.BroadcastAdmin("download_tasks.delete", downloadTask); err != nil {
			log.Println(err)
		}

		// Delete download task when done
		models.DownloadTaskModel.Where("id", downloadTask.ID).Delete()
	}
}
