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

func Websocket(c *fiber.Ctx) error {
	upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) {
		var authUser *models.User
		for {
			// Parse incoming json messages
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
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
