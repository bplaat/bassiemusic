package controllers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bplaat/bassiemusic/models"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

type Message struct {
	Type     string  `json:"type"`
	Token    string  `json:"token"`
	TrackID  string  `json:"track_id"`
	Position float32 `json:"position"`
}

// Maintain a list of WebSocket connections
var connections []*websocket.Conn

func Websocket(c *fiber.Ctx) error {
	upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) { //nolint
		var authUser *models.User
		var addedConnection = false

		for {
			// Parse incoming json messages
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}

				// Remove the connection from the connections list
				for i, c := range connections {
					if c == conn {
						connections = append(connections[:i], connections[i+1:]...)
						break
					}
				}

				return
			}
			var message Message
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				log.Println(err)
				return
			}

			// Guest messages
			if message.Type == "auth" {
				session := models.SessionModel().With("user").Where("token", message.Token).WhereRaw("`expires_at` > ?", time.Now()).First()
				if session == nil {
					response := "{\"success\":false}"
					if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
						log.Println(err)
						return
					}
					continue
				}
				authUser = session.User

				if !addedConnection && authUser.Role == 1 {
					// Append new admin connection to the connections list
					connections = append(connections, conn)
					addedConnection = true

					// Send current active download tasks
					downloadTasks := models.DownloadTaskModel().Get()

					if len(downloadTasks) > 0 {
						response, _ := json.Marshal(fiber.Map{
							"success":       true,
							"type":          "allTasks",
							"downloadTasks": downloadTasks,
						})
						if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
							log.Println(err)
							return
						}
					}
				}

				continue
			}

			// Authed messages
			if authUser == nil {
				response := "{\"success\":false}"
				if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
					log.Println(err)
					return
				}
				continue
			}

			if message.Type == "track_play" {
				models.HandleTrackPlay(authUser, message.TrackID, message.Position)
				continue
			}
		}
	})
	return c.SendString("")
}

// Send a message to all connected admin WebSocket devices
func SendMessageToAll(message []byte) error {
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}
	}
	return nil
}
