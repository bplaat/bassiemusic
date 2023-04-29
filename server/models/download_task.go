package models

import (
	"time"

	"github.com/bplaat/bassiemusic/core/database"
)

type DownloadTask struct {
	ID               string             `column:"id,uuid" json:"id"`
	Type             DownloadTaskType   `column:"type,int" json:"-"`
	TypeString       string             `json:"type"`
	DeezerID         int64              `column:"deezer_id,bigint" json:"deezer_id"`
	DisplayName      string             `column:"display_name,string" json:"display_name"`
	Status           DownloadTaskStatus `column:"status,int" json:"-"`
	StatusString     string             `json:"status"`
	DownloadedTracks int                `column:"downloaded_tracks,int" json:"downloaded_tracks"`
	TotalTracks      int                `column:"total_tracks, int" json:"total_tracks"`
	CreatedAt        time.Time          `column:"created_at,timestamp" json:"created_at"`
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
