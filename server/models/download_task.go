package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type DownloadTaskData struct {
	DeezerID  int64  `json:"deezer_id"`
	YoutubeID string `json:"youtube_id"`
	TrackID   string `json:"track_id"`
	ArtistID  string `json:"artist_id"`
	AlbumID   string `json:"album_id"`
}

type DownloadTask struct {
	ID           uuid.Uuid          `column:"id" json:"id"`
	Type         DownloadTaskType   `column:"type" json:"-"`
	TypeString   string             `json:"type"`
	Data         string             `column:"data" json:"-"`
	DisplayName  string             `column:"display_name" json:"display_name"`
	Status       DownloadTaskStatus `column:"status" json:"-"`
	StatusString string             `json:"status"`
	Progress     float32            `column:"progress" json:"progress"`
	CreatedAt    time.Time          `column:"created_at" json:"created_at"`
	DeezerID     *int64             `json:"deezer_id,omitempty"`
	YoutubeID    *string            `json:"youtube_id,omitempty"`
	TrackID      *uuid.Uuid         `json:"track_id,omitempty"`
	ArtistID     *uuid.Uuid         `json:"artist_id,omitempty"`
	AlbumID      *uuid.Uuid         `json:"album_id,omitempty"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1
const DownloadTaskTypeYoutubeTrack DownloadTaskType = 2
const DownloadTaskTypeUpdateDeezerArtist DownloadTaskType = 3
const DownloadTaskTypeUpdateDeezerAlbum DownloadTaskType = 4

type DownloadTaskStatus int

const DownloadTaskStatusPending DownloadTaskStatus = 0
const DownloadTaskStatusDownloading DownloadTaskStatus = 1

var DownloadTaskModel *database.Model[DownloadTask] = (&database.Model[DownloadTask]{
	TableName: "download_tasks",
	Process: func(downloadTask *DownloadTask) {
		if downloadTask.Type == DownloadTaskTypeDeezerArtist {
			downloadTask.TypeString = "deezer_artist"
		}
		if downloadTask.Type == DownloadTaskTypeDeezerAlbum {
			downloadTask.TypeString = "deezer_album"
		}
		if downloadTask.Type == DownloadTaskTypeYoutubeTrack {
			downloadTask.TypeString = "youtube_track"
		}
		if downloadTask.Type == DownloadTaskTypeUpdateDeezerArtist {
			downloadTask.TypeString = "update_deezer_artist"
		}
		if downloadTask.Type == DownloadTaskTypeUpdateDeezerAlbum {
			downloadTask.TypeString = "update_deezer_album"
		}
		if downloadTask.Status == DownloadTaskStatusPending {
			downloadTask.StatusString = "pending"
		}
		if downloadTask.Status == DownloadTaskStatusDownloading {
			downloadTask.StatusString = "downloading"
		}

		// Parse download task data
		var data DownloadTaskData
		if err := json.Unmarshal([]byte(downloadTask.Data), &data); err != nil {
			log.Fatalln(err)
		}
		if data.DeezerID != 0 {
			downloadTask.DeezerID = &data.DeezerID
		}
		if data.YoutubeID != "" {
			downloadTask.YoutubeID = &data.YoutubeID
		}
		if data.TrackID != "" {
			trackID, _ := uuid.Parse(data.TrackID)
			downloadTask.TrackID = &trackID
		}
		if data.TrackID != "" {
			ArtistID, _ := uuid.Parse(data.ArtistID)
			downloadTask.ArtistID = &ArtistID
		}
		if data.AlbumID != "" {
			AlbumID, _ := uuid.Parse(data.AlbumID)
			downloadTask.AlbumID = &AlbumID
		}
	},
}).Init()
