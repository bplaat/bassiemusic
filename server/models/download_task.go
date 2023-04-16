package models

import (
	"time"

	"github.com/bplaat/bassiemusic/database"
)

type DownloadTask struct {
	ID           string           `column:"id,uuid" json:"id"`
	Type         DownloadTaskType `column:"type,int" json:"-"`
	TypeString   string           `json:"type"`
	DeezerID     int64            `column:"deezer_id,bigint" json:"deezer_id"`
	DisplayName  string           `column:"display_name,string" json:"display_name"`
	Status       StatusType       `column:"status,int" json:"-"`
	StatusString string           `json:"status"`
	Progress     int              `column:"progress,int" json:"progress"`
	CreatedAt    time.Time        `column:"created_at,timestamp" json:"created_at"`
}

type DownloadTaskType int
type StatusType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1

const DownloadStatusTypePending StatusType = 0
const DownloadStatusTypeDownloading StatusType = 1

func DownloadTaskModel() *database.Model[DownloadTask] {
	return (&database.Model[DownloadTask]{
		TableName: "download_tasks",
		Process: func(downloadTask *DownloadTask) {
			if downloadTask.Type == DownloadTaskTypeDeezerArtist {
				downloadTask.TypeString = "deezer_artist"
			}
			if downloadTask.Type == DownloadTaskTypeDeezerAlbum {
				downloadTask.TypeString = "deezer_album"
			}

			if downloadTask.Status == DownloadStatusTypePending {
				downloadTask.StatusString = "pending"
			}
			if downloadTask.Status == DownloadStatusTypeDownloading {
				downloadTask.StatusString = "downloading"
			}
		},
	}).Init()
}
