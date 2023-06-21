package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type Data struct {
	DeezerID  int64  `json:"deezer_id"`
	YoutubeID string `json:"youtube_id"`
	TrackID   string `json:"track_id"`
}

type DownloadTask struct {
	ID           uuid.Uuid          `column:"id" json:"id"`
	Type         DownloadTaskType   `column:"type" json:"-"`
	TypeString   string             `json:"type"`
	JsonData     string             `column:"data" json:"-"`
	DeezerID     int64              `json:"deezer_id"`
	YoutubeID    string             `json:"youtube_id"`
	TrackID      uuid.Uuid          `json:"track_id"`
	DisplayName  string             `column:"display_name" json:"display_name"`
	Status       DownloadTaskStatus `column:"status" json:"-"`
	StatusString string             `json:"status"`
	Progress     float32            `column:"progress" json:"progress"`
	CreatedAt    time.Time          `column:"created_at" json:"created_at"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1
const DownloadTaskTypeYoutubeTrack DownloadTaskType = 2

type DownloadTaskStatus int

const DownloadTaskStatusPending DownloadTaskStatus = 0
const DownloadTaskStatusDownloading DownloadTaskStatus = 1

var DownloadTaskModel *database.Model[DownloadTask] = (&database.Model[DownloadTask]{
	TableName: "download_tasks",
	Process: func(downloadTask *DownloadTask) {
		var data Data
		if err := json.Unmarshal([]byte(downloadTask.JsonData), &data); err != nil {
			log.Println(err)
			return
		}

		downloadTask.YoutubeID = data.YoutubeID
		downloadTask.DeezerID = data.DeezerID

		trackID, _ := uuid.Parse(data.TrackID)
		downloadTask.TrackID = trackID

		if downloadTask.Type == DownloadTaskTypeDeezerArtist {
			downloadTask.TypeString = "deezer_artist"
		}
		if downloadTask.Type == DownloadTaskTypeDeezerAlbum {
			downloadTask.TypeString = "deezer_album"
		}
		if downloadTask.Type == DownloadTaskTypeYoutubeTrack {
			downloadTask.TypeString = "youtube_track"
		}

		if downloadTask.Status == DownloadTaskStatusPending {
			downloadTask.StatusString = "pending"
		}
		if downloadTask.Status == DownloadTaskStatusDownloading {
			downloadTask.StatusString = "downloading"
		}
	},
}).Init()
