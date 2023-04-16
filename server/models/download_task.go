package models

import (
	"time"

	"github.com/bplaat/bassiemusic/database"
)

type DownloadTask struct {
	ID         string           `column:"id,uuid" json:"id"`
	Type       DownloadTaskType `column:"type,int" json:"-"`
	TypeString string           `json:"type"`
	DeezerID   int64            `column:"deezer_id,bigint" json:"deezer_id"`
	Status     int              `column:"status,int" json:"status"`
	Progress   int              `column:"progress,int" json:"progress"`
	CreatedAt  time.Time        `column:"created_at,timestamp" json:"created_at"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1

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
		},
	}).Init()
}
