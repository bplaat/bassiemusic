package websocket

import (
	"encoding/json"
	"time"

	"github.com/bplaat/bassiemusic/models"
)

type AuthLoginMessage struct {
	Type string `json:"type"`
	Data struct {
		Token string `json:"token"`
	}
}

func HandleAuthMessages(connection *Connection, messageType string, messageBytes []byte) error {
	if messageType == "auth.validate" {
		// Parse websocket JSON message
		var message AuthLoginMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			return err
		}

		// Update tracks counter
		session := models.SessionModel.With("user").Where("token", message.Data.Token).WhereRaw("`expires_at` > ?", time.Now()).First()
		if session == nil {
			if err := connection.Send("auth.login.response", Map{"success": false}); err != nil {
				return err
			}
			return nil
		}

		// Set current auth user to session user
		connection.AuthUser = session.User

		// Send success response
		if err := connection.Send("auth.login.response", Map{"success": true}); err != nil {
			return err
		}
		return nil
	}

	return nil
}
