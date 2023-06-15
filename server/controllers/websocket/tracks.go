package websocket

import (
	"encoding/json"

	"github.com/bplaat/bassiemusic/core/uuid"
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
		TrackID, err := uuid.Parse(message.Data.TrackID)
		if err != nil {
			return err
		}
		models.HandleTrackPlay(connection.AuthUser, TrackID, message.Data.Position)
		return nil
	}

	return nil
}
