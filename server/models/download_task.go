package models

import (
	"database/sql"
	"time"
)

type DownloadTask struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	DeezerID  int64     `json:"deezer_id"`
	Singles   bool      `json:"singles"`
	CreatedAt time.Time `json:"created_at"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1

func DownloadTaskScan(downloadTaskQuery *sql.Rows) DownloadTask {
	var downloadTask DownloadTask
	var downloadTaskType DownloadTaskType
	downloadTaskQuery.Scan(&downloadTask.ID, &downloadTaskType, &downloadTask.DeezerID, &downloadTask.Singles, &downloadTask.CreatedAt)
	if downloadTaskType == DownloadTaskTypeDeezerArtist {
		downloadTask.Type = "deezer_artist"
	}
	if downloadTaskType == DownloadTaskTypeDeezerAlbum {
		downloadTask.Type = "deezer_album"
	}
	return downloadTask
}
