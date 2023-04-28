package websocket

import (
	"github.com/bplaat/bassiemusic/models"
)

func HandleDownloadTasksMessages(connection *Connection, messageType string, messageBytes []byte) error {
	if messageType == "download_tasks.init" {
		// Send current download tasks
		downloadTasks := models.DownloadTaskModel.OrderBy("created_at").Get()
		if len(downloadTasks) > 0 {
			if err := connection.Send("download_tasks.init.response", downloadTasks); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}
