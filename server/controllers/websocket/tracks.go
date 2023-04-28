package websocket

import (
	"encoding/json"

	"github.com/bplaat/bassiemusic/models"
)

type TracksPlayMessage struct {
	Type string `json:"type"`
	Data struct {
		TrackID  string  `json:"track_id"`
		Position float32 `json:"position"`
	}
}

func HandleTracksMessages(connection *Connection, messageType string, messageBytes []byte) error {
	if messageType == "tracks.play" {
		// Parse websocket JSON message
		var message TracksPlayMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			return err
		}

		// Update tracks counter
		models.HandleTrackPlay(connection.AuthUser, message.Data.TrackID, message.Data.Position)
		return nil
	}

	return nil
}
