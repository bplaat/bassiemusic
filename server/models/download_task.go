package models

import (
	"time"

	"github.com/bplaat/bassiemusic/database"
)

type DownloadTask struct {
	ID        string           `column:"id,uuid" json:"id"`
	TypeInt   DownloadTaskType `column:"type,int" json:"-"`
	Type      string           `json:"type"`
	DeezerID  int64            `column:"deezer_id,bigint" json:"deezer_id"`
	CreatedAt time.Time        `column:"created_at,timestamp" json:"created_at"`
}

type DownloadTaskType int

const DownloadTaskTypeDeezerArtist DownloadTaskType = 0
const DownloadTaskTypeDeezerAlbum DownloadTaskType = 1

func DownloadTaskModel() *database.Model[DownloadTask] {
	return (&database.Model[DownloadTask]{
		TableName: "download_tasks",
		Process: func(downloadTask *DownloadTask) {
			if downloadTask.TypeInt == DownloadTaskTypeDeezerArtist {
				downloadTask.Type = "deezer_artist"
			}
			if downloadTask.TypeInt == DownloadTaskTypeDeezerAlbum {
				downloadTask.Type = "deezer_album"
			}
		},
	}).Init()
}
