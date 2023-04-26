package websocket

import (
	"encoding/json"
	"log"

	"github.com/bplaat/bassiemusic/models"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Map map[string]any

type Message struct {
	Type string `json:"type"`
}

// Websocket upgrader
var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

// Open connections list
type Connection struct {
	Conn     *websocket.Conn
	AuthUser *models.User
}

func (connection *Connection) Send(messageType string, data any) error {
	jsonMessage, _ := json.Marshal(Map{
		"type": messageType,
		"data": data,
	})
	if err := connection.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		return err
	}
	return nil
}

var connections []*Connection

// Websocket fiber handler
func ServerHandle(c *fiber.Ctx) error {
	upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) { //nolint
		connection := Connection{
			Conn:     conn,
			AuthUser: nil,
		}
		connectionAdded := false
		for {
			// Parse incoming messages
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}

				// Remove closed connections
				for i, connection := range connections {
					if connection.Conn == conn {
						connections = append(connections[:i], connections[i+1:]...)
						break
					}
				}
				return
			}

			// Add open connection
			if !connectionAdded {
				connectionAdded = true
				connections = append(connections, &connection)
			}

			// Parse message type
			var message Message
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				log.Println(err)
				return
			}

			// Handle guest messages
			if err := HandleAuthMessages(&connection, message.Type, messageBytes); err != nil {
				log.Println(err)
			}

			// Handle authed messages
			if connection.AuthUser != nil {
				if err := HandleTracksMessages(&connection, message.Type, messageBytes); err != nil {
					log.Println(err)
				}

				// Handle admin messages
				if connection.AuthUser.Role == models.UserRoleAdmin {
					if err := HandleDownloadTasksMessages(&connection, message.Type, messageBytes); err != nil {
						log.Println(err)
					}
				}
			}
		}
	})
	return c.SendString("")
}

// Broad cast message functions
func broadcast(minimalUserRole models.UserRole, messageType string, data any) error {
	// Send json message to all open connections that have minimal user role
	for _, connection := range connections {
		if connection.AuthUser != nil {
			if connection.AuthUser.Role >= minimalUserRole {
				if err := connection.Send(messageType, data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Broadcast(messageType string, data any) error {
	return broadcast(models.UserRoleNormal, messageType, data)
}

func BroadcastAdmin(messageType string, data any) error {
	return broadcast(models.UserRoleAdmin, messageType, data)
}
