package routes

import (
	"JobHiraMicroservice/models"
	"JobHiraMicroservice/websockets"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func WebSockets(ws fiber.Router) {
	ws.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	// connect to websocket
	ws.Get("/", websocket.New(func(c *websocket.Conn) {

		defer func() {
			websockets.Unregister <- c
			err := c.Close()
			if err != nil {
				return
			}
		}()

		// get client key
		clientKey := c.Query(
			"key",
		)
		if clientKey == "" || !ClientKeyExists(clientKey) {
			log.Println("client key not found")
			return
		}

		// validate key exist
		clientToRegister := websockets.IncomingConnection{
			Key:        clientKey,
			Connection: c,
		}

		// client register
		websockets.RegisterWithKey <- clientToRegister

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				websockets.Broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}

	}))

}

func ClientKeyExists(key string) bool {
	Application := models.Application{}
	if err := models.DB.Where("client_key = ?", key).First(&Application).Error; err != nil {
		return false
	}
	return true
}
