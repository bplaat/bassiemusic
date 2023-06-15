package models

import (
	"time"

	"github.com/bplaat/bassiemusic/core/database"
	"github.com/bplaat/bassiemusic/core/uuid"
)

type DownloadTask struct {
	ID           uuid.Uuid          `column:"id" json:"id"`
	Type         DownloadTaskType   `column:"type" json:"-"`
	TypeString   string             `json:"type"`
	DeezerID     int64              `column:"deezer_id" json:"deezer_id"`
	DisplayName  string             `column:"display_name" json:"display_name"`
	Status       DownloadTaskStatus `column:"status" json:"-"`
	StatusString string             `json:"status"`
	Progress     float32            `column:"progress" json:"progress"`
	CreatedAt    time.Time          `column:"created_at" json:"created_at"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1

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

		if downloadTask.Status == DownloadTaskStatusPending {
			downloadTask.StatusString = "pending"
		}
		if downloadTask.Status == DownloadTaskStatusDownloading {
			downloadTask.StatusString = "downloading"
		}
	},
}).Init()
