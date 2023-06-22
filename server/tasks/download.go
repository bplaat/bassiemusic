package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
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

func createArtist(deezerID int, name string, sync bool) uuid.Uuid {
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
		"id":        artistID,
		"name":      name,
		"deezer_id": deezerID,
		"sync":      sync,
	})

	if deezerArtist.PictureMedium != "https://e-cdns-images.dzcdn.net/images/artist//250x250-000000-80-0-0.jpg" {
		utils.DeezerFetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artistID))
		utils.DeezerFetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artistID))
		utils.DeezerFetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artistID))
	}
	return artistID
}

func createGenre(deezerID int, name string) uuid.Uuid {
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
		"id":        genreID,
		"name":      name,
		"deezer_id": deezerID,
	})
	if deezerGenre.PictureMedium != "https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg" {
		utils.DeezerFetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genreID))
		utils.DeezerFetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genreID))
		utils.DeezerFetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genreID))
	}
	return genreID
}

func CreateTrack(albumID uuid.Uuid, deezerID int64) {
	// Get Deezer track info
	var deezerTrack structs.DeezerTrack
	if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/track/%d", deezerID), &deezerTrack); err != nil {
		log.Fatalln(err)
	}

	// Create track
	trackID := uuid.New()
	models.TrackModel.Create(database.Map{
		"id":         trackID,
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
	for index, artist := range deezerTrack.Contributors {
		artistID := createArtist(artist.ID, artist.Name, false)
		models.TrackArtistModel.Create(database.Map{
			"track_id":  trackID,
			"artist_id": artistID,
			"position":  index + 1,
		})
	}
}

func SearchAndDownloadTrackMusic(track *models.Track, youtubeID string, findYoutubeID bool) error {
	youtubeDuration := 0

	if findYoutubeID {
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

		fallPoints := 0
		lowestScore := 1000000.0

		for {
			var youtubeVideo structs.YoutubeVideo
			if err := json.NewDecoder(stdout).Decode(&youtubeVideo); err != nil {
				break
			}

			score := float32(youtubeVideo.Duration) - track.Duration
			if score < 0 {
				score = score - score - score
			}
			if score+float32(fallPoints) < float32(lowestScore) {
				lowestScore = float64(score) + float64(fallPoints)
				youtubeID = youtubeVideo.ID
				youtubeDuration = youtubeVideo.Duration
			}

			fallPoints += consts.PUNISHMENT_POINTS
		}
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

func DownloadAlbum(deezerAlbum structs.DeezerAlbum, downloadTask *models.DownloadTask, downloadedTracks *int, totalTracks *int) int {
	// Check if album already exists
	if models.AlbumModel.Where("title", deezerAlbum.Title).First() != nil {
		return len(deezerAlbum.Tracks.Data)
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
		"id":          albumID,
		"type":        albumType,
		"title":       deezerAlbum.Title,
		"released_at": deezerAlbum.ReleaseDate,
		"explicit":    deezerAlbum.ExplicitLyrics,
		"deezer_id":   deezerAlbum.ID,
	})

	utils.DeezerFetchFile(deezerAlbum.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", albumID))
	utils.DeezerFetchFile(deezerAlbum.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", albumID))
	utils.DeezerFetchFile(deezerAlbum.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", albumID))

	// Create album genre bindings
	for _, genre := range deezerAlbum.Genres.Data {
		genreID := createGenre(genre.ID, genre.Name)
		models.AlbumGenreModel.Create(database.Map{
			"album_id": albumID,
			"genre_id": genreID,
		})
	}

	// Create album artist bindings
	for index, artist := range deezerAlbum.Contributors {
		artistID := createArtist(artist.ID, artist.Name, false)
		models.AlbumArtistModel.Create(database.Map{
			"album_id":  albumID,
			"artist_id": artistID,
			"position":  index + 1,
		})
	}

	// Create album tracks
	for _, incompleteTrack := range deezerAlbum.Tracks.Data {
		CreateTrack(albumID, incompleteTrack.ID)
	}

	// Download album tracks music
	for _, deezerTrack := range deezerAlbum.Tracks.Data {
		track := models.TrackModel.With("album", "artists").Where("album_id", albumID).Where("title", deezerTrack.Title).First()
		if err := SearchAndDownloadTrackMusic(track, "", true); err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		log.Printf("[DOWNLOAD] %s - %d-%d - %s\n", deezerAlbum.Title, track.Disk, track.Position, track.Title)

		// Log when a track haves been downloaded
		if downloadTask != nil && downloadedTracks != nil && totalTracks != nil {
			*downloadedTracks += 1

			// Update download task progress
			downloadTask.Progress = float32(math.Round((float64(*downloadedTracks) / float64(*totalTracks)) * 100))
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"progress": downloadTask.Progress,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}
		}
	}

	log.Printf("[DOWNLOAD] Done downloading album\n")
	return *downloadedTracks
}

func fetchAlbums(DeezerID int64) ([]structs.DeezerAlbum, int) {
	nextUrl := fmt.Sprintf("https://api.deezer.com/artist/%d/albums", DeezerID)
	var fullAlbums []structs.DeezerAlbum
	totalTracks := 0
	for {
		var artistAlbums structs.DeezerArtistAlbums
		if err := utils.DeezerFetch(nextUrl, &artistAlbums); err != nil {
			log.Fatalln(err)
		}

		for _, album := range artistAlbums.Data {
			var deezerAlbum structs.DeezerAlbum
			if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/album/%d", album.ID), &deezerAlbum); err != nil {
				log.Fatalln(err)
			}
			totalTracks += len(deezerAlbum.Tracks.Data)
			fullAlbums = append(fullAlbums, deezerAlbum)
		}

		if artistAlbums.Next == "" {
			break
		} else {
			nextUrl = artistAlbums.Next
		}
	}
	return fullAlbums, totalTracks
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
		if downloadTask.Type == models.DownloadTaskTypeYoutubeTrack {
			// Update download task status
			downloadTask.StatusString = "downloading"
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			track := models.TrackModel.Find(*downloadTask.TrackID)
			os.Remove(fmt.Sprintf("storage/tracks/%s.m4a", *downloadTask.TrackID))

			if err := SearchAndDownloadTrackMusic(track, *downloadTask.YoutubeID, false); err != nil && err != io.EOF {
				log.Fatalln(err)
			}
		}

		if downloadTask.Type == models.DownloadTaskTypeUpdateDeezerArtist {
			// Update download task status
			downloadTask.StatusString = "downloading"
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}

			// Check if all tracks are downloaded
			artistAlbums, totalTracks := fetchAlbums(*downloadTask.DeezerID)
			downloadedTracks := 0
			for _, album := range artistAlbums {
				if strings.Contains(album.Title, "Deezer") {
					continue
				}

				// Check for each track if it exist and it already downloaded
				fullAlbum := models.AlbumModel.WhereRaw("`deezer_id` LIKE ?", "%"+strconv.FormatInt(int64(album.ID), 10)+"%").First()
				if fullAlbum == nil {
					// Download missing album
					downloadedTracks = DownloadAlbum(album, downloadTask, &downloadedTracks, &totalTracks)
				} else {
					// Walk over each track to see if a track is not downloaded
					for _, deezerTrack := range album.Tracks.Data {
						fullTrack := models.TrackModel.With("album", "artists").Where("album_id", fullAlbum.ID).Where("title", deezerTrack.Title).First()

						if fullTrack == nil {
							// Todo: Download track with new database entry
						} else {
							var music string
							if _, err := os.Stat(fmt.Sprintf("storage/tracks/%s.m4a", fullTrack.ID)); err == nil {
								music = fmt.Sprintf("%s/tracks/%s.m4a", os.Getenv("STORAGE_URL"), fullTrack.ID)
							}
							if music == "" {
								SearchAndDownloadTrackMusic(fullTrack, "", true)
							}
							downloadedTracks++
						}
					}
				}

				// Update download task progress
				downloadTask.Progress = float32(math.Round((float64(downloadedTracks) / float64(totalTracks)) * 100))
				models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
					"progress": downloadTask.Progress,
				})

				// Broadcast download tasks update message to admins
				if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
					log.Println(err)
				}
			}
		}

		if downloadTask.Type == models.DownloadTaskTypeDeezerArtist {
			// Update download task status
			downloadTask.Status = models.DownloadTaskStatusDownloading
			downloadTask.StatusString = "downloading"
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}

			// Create artist
			var deezerArtist structs.DeezerArtist
			if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/artist/%d", *downloadTask.DeezerID), &deezerArtist); err != nil {
				log.Fatalln(err)
			}
			createArtist(deezerArtist.ID, deezerArtist.Name, true)

			// Download artist albums
			artistAlbums, totalTracks := fetchAlbums(*downloadTask.DeezerID)
			downloadedTracks := 0

			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}
			for _, album := range artistAlbums {
				if strings.Contains(album.Title, "Deezer") {
					continue
				}
				downloadedTracks = DownloadAlbum(album, downloadTask, &downloadedTracks, &totalTracks)

				// Update download task progress
				downloadTask.Progress = float32(math.Round((float64(downloadedTracks) / float64(totalTracks)) * 100))
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
			// Fetch album data
			var deezerAlbum structs.DeezerAlbum
			if err := utils.DeezerFetch(fmt.Sprintf("https://api.deezer.com/album/%d", *downloadTask.DeezerID), &deezerAlbum); err != nil {
				log.Fatalln(err)
			}

			// Update download task status
			downloadTask.Status = models.DownloadTaskStatusDownloading
			downloadTask.StatusString = "downloading"
			models.DownloadTaskModel.Where("id", downloadTask.ID).Update(database.Map{
				"status": downloadTask.Status,
			})

			// Broadcast download tasks update message to admins
			if err := websocket.BroadcastAdmin("download_tasks.update", downloadTask); err != nil {
				log.Println(err)
			}

			// Download album
			downloadedTracks := 0
			totalTracks := len(deezerAlbum.Tracks.Data)
			downloadedTracks = DownloadAlbum(deezerAlbum, downloadTask, &downloadedTracks, &totalTracks)
		}

		// Broadcast download tasks delete message to admins
		if err := websocket.BroadcastAdmin("download_tasks.delete", downloadTask); err != nil {
			log.Println(err)
		}

		// Delete download task when done
		models.DownloadTaskModel.Where("id", downloadTask.ID).Delete()
	}
}
