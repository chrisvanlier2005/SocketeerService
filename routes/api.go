package routes

import (
	"JobHiraMicroservice/controllers"
	"JobHiraMicroservice/models"
	"JobHiraMicroservice/websockets"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(api fiber.Router) {
	api.Get("/broadcast", func(c *fiber.Ctx) error {
		serverKey := c.Query("key")
		message := c.Query("message")
		if serverKey == "" || message == "" {
			return c.JSON(map[string]string{
				"code":     "error",
				"message:": "server key & message is required",
			})
		}

		err := websockets.BroadcastMessage(websockets.BroadcastMessageType{
			Message:   message,
			ServerKey: serverKey,
			ClientKey: "test-client-key",
		}, websockets.Channel{
			ChannelName: "test-channel",
		})
		if err != nil {
			return c.JSON(map[string]string{
				"code":     "error",
				"message:": "server key not found",
			})
		}

		return c.JSON(map[string]string{
			"code":       "success",
			"server-key": serverKey,
		})
	})

	api.Get("/connections", func(c *fiber.Ctx) error {
		// get the server key
		serverKey := c.Query("key")
		if serverKey == "" {
			return c.JSON(map[string]string{
				"code":     "error",
				"message:": "server key and message is required",
			})
		}

		application := models.Application{}
		if err := models.DB.Where("server_key = ?", serverKey).First(&application).Error; err != nil {
			return c.JSON(map[string]string{
				"code":    "error",
				"message": "application with key " + serverKey + " Does not exist",
			})
		}

		connections := make(map[string]*websockets.Client)
		for _, value := range websockets.Clients {
			if value.ClientKey == application.ClientKey {
				connections[value.Id] = value
			}
		}
		return c.JSON(connections)
	})

	api.Get("/users/:user", controllers.UserController{}.Show)
}
