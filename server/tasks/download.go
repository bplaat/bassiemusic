package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bplaat/bassiemusic/consts"
	"github.com/bplaat/bassiemusic/database"
	"github.com/bplaat/bassiemusic/models"
	"github.com/bplaat/bassiemusic/structs"
	"github.com/bplaat/bassiemusic/utils"
	"github.com/bplaat/bassiemusic/utils/uuid"
)

func createArtist(deezerID int, name string, synced bool) string {
	// Check if artist already exists
	artist := models.ArtistModel(nil).Where("name", name).First()
	if artist != nil {
		if artist.Synced != synced {
			models.ArtistModel(nil).Where("id", artist.ID).Update(database.Map{
				"synced": synced,
			})
		}
		return artist.ID
	}

	// Get Deezer artist info
	var deezerArtist structs.DeezerArtist
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", deezerID), &deezerArtist); err != nil {
		log.Fatalln(err)
	}

	// Create artist
	artistID := uuid.New()
	models.ArtistModel(nil).Create(database.Map{
		"id":        artistID.String(),
		"name":      name,
		"deezer_id": deezerID,
		"synced":    synced,
	})
	if deezerArtist.PictureMedium != "https://e-cdns-images.dzcdn.net/images/artist//250x250-000000-80-0-0.jpg" {
		utils.FetchFile(deezerArtist.PictureMedium, fmt.Sprintf("storage/artists/small/%s.jpg", artistID.String()))
		utils.FetchFile(deezerArtist.PictureBig, fmt.Sprintf("storage/artists/medium/%s.jpg", artistID.String()))
		utils.FetchFile(deezerArtist.PictureXl, fmt.Sprintf("storage/artists/large/%s.jpg", artistID.String()))
	}
	return artistID.String()
}

func createGenre(deezerID int, name string) string {
	// Check if genre already exists
	genre := models.GenreModel(nil).Where("name", name).First()
	if genre != nil {
		return genre.ID
	}

	// Get Deezer genre info
	var deezerGenre structs.DeezerGenre
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/genre/%d", deezerID), &deezerGenre); err != nil {
		log.Fatalln(err)
	}

	// Create genre
	genreID := uuid.New()
	models.GenreModel(nil).Create(database.Map{
		"id":        genreID.String(),
		"name":      name,
		"deezer_id": deezerID,
	})
	if deezerGenre.PictureMedium != "https://e-cdns-images.dzcdn.net/images/misc//250x250-000000-80-0-0.jpg" {
		utils.FetchFile(deezerGenre.PictureMedium, fmt.Sprintf("storage/genres/small/%s.jpg", genreID.String()))
		utils.FetchFile(deezerGenre.PictureBig, fmt.Sprintf("storage/genres/medium/%s.jpg", genreID.String()))
		utils.FetchFile(deezerGenre.PictureXl, fmt.Sprintf("storage/genres/large/%s.jpg", genreID.String()))
	}
	return genreID.String()
}

func CreateTrack(albumID string, deezerID int) {
	// Get Deezer track info
	var deezerTrack structs.DeezerTrack
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/track/%d", deezerID), &deezerTrack); err != nil {
		log.Fatalln(err)
	}

	// Create track
	trackID := uuid.New()
	models.TrackModel(nil).Create(database.Map{
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
		models.TrackArtistModel().Create(database.Map{
			"track_id":  trackID.String(),
			"artist_id": artistID,
		})
	}
}

func SearchAndDownloadTrackMusic(track *models.Track) error {
	// Search for youtube video
	searchCommand := exec.Command("yt-dlp", "--dump-json", fmt.Sprintf("ytsearch25:%s - %s - %s", (*track.Artists)[0].Name, track.Album.Title, track.Title))
	log.Println(searchCommand.String())
	stdout, err := searchCommand.StdoutPipe()
	if err != nil {
		return err
	}
	if err := searchCommand.Start(); err != nil {
		return err
	}
	for {
		var youtubeVideo structs.YoutubeVideo
		if err := json.NewDecoder(stdout).Decode(&youtubeVideo); err != nil {
			return err
		}

		// When video duration is in slack download it
		if track.Duration >= float32(youtubeVideo.Duration-consts.TRACK_DURATION_SLACK) &&
			track.Duration <= float32(youtubeVideo.Duration+consts.TRACK_DURATION_SLACK) {
			if err := searchCommand.Process.Kill(); err != nil {
				log.Fatalln(err)
			}

			// Download right youtube video
			downloadCommand := exec.Command("yt-dlp", "-f", "bestaudio[ext=m4a]", fmt.Sprintf("https://www.youtube.com/watch?v=%s", youtubeVideo.ID),
				"-o", fmt.Sprintf("storage/tracks/%s.m4a", track.ID))
			log.Println(downloadCommand.String())

			if err := downloadCommand.Run(); err != nil {
				log.Fatalln(err)
			}

			// Update track
			models.TrackModel(nil).Where("id", track.ID).Update(database.Map{
				"duration":   youtubeVideo.Duration,
				"youtube_id": youtubeVideo.ID,
			})
			return nil
		}
	}
}

func DownloadAlbum(deezerID int) {
	// Get Deezer album info
	var deezerAlbum structs.DeezerAlbum
	if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/album/%d", deezerID), &deezerAlbum); err != nil {
		log.Fatalln(err)
	}

	// Check if album already exists
	if models.AlbumModel(nil).Where("title", deezerAlbum.Title).First() != nil {
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
	models.AlbumModel(nil).Create(database.Map{
		"id":          albumID.String(),
		"type":        albumType,
		"title":       deezerAlbum.Title,
		"released_at": deezerAlbum.ReleaseDate,
		"explicit":    deezerAlbum.ExplicitLyrics,
		"deezer_id":   deezerID,
	})
	utils.FetchFile(deezerAlbum.CoverMedium, fmt.Sprintf("storage/albums/small/%s.jpg", albumID.String()))
	utils.FetchFile(deezerAlbum.CoverBig, fmt.Sprintf("storage/albums/medium/%s.jpg", albumID.String()))
	utils.FetchFile(deezerAlbum.CoverXl, fmt.Sprintf("storage/albums/large/%s.jpg", albumID.String()))

	// Create album genre bindings
	for _, genre := range deezerAlbum.Genres.Data {
		genreID := createGenre(genre.ID, genre.Name)
		models.AlbumGenreModel().Create(database.Map{
			"album_id": albumID.String(),
			"genre_id": genreID,
		})
	}

	// Create album artist bindings
	for _, artist := range deezerAlbum.Contributors {
		artistID := createArtist(artist.ID, artist.Name, false)
		models.AlbumArtistModel().Create(database.Map{
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
		track := models.TrackModel(nil).With("album", "artists").Where("album_id", albumID.String()).Where("title", deezerTrack.Title).First()
		if err := SearchAndDownloadTrackMusic(track); err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		log.Printf("[DOWNLOAD] %s - %d-%d - %s\n", deezerAlbum.Title, track.Disk, track.Position, track.Title)
	}

	log.Printf("[DOWNLOAD] Done downloading album\n")
}

func DownloadTask() {
	for {
		// Wait a little while
		time.Sleep(5 * time.Second)

		// Get first download task
		downloadTask := models.DownloadTaskModel().First()
		if downloadTask == nil {
			continue
		}

		// Do download task
		if downloadTask.Type == "deezer_artist" {
			// Create artist
			var deezerArtist structs.DeezerArtist
			if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d", downloadTask.DeezerID), &deezerArtist); err != nil {
				log.Fatalln(err)
			}
			createArtist(deezerArtist.ID, deezerArtist.Name, true)

			// Download artist albums
			var artistAlbums structs.DeezerArtistAlbums
			if err := utils.FetchJson(fmt.Sprintf("https://api.deezer.com/artist/%d/albums", downloadTask.DeezerID), &artistAlbums); err != nil {
				log.Fatalln(err)
			}
			for _, album := range artistAlbums.Data {
				if !strings.Contains(album.Title, "Deezer") {
					continue
				}
				DownloadAlbum(album.ID)
			}
		}
		if downloadTask.Type == "deezer_album" {
			DownloadAlbum(int(downloadTask.DeezerID))
		}

		// Delete download task when done
		models.DownloadTaskModel().Where("id", downloadTask.ID).Delete()
	}
}
