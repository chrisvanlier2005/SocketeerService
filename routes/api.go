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

	api.Get("/applications", func(c *fiber.Ctx) error {
		var apps []models.Application
		models.DB.Find(&apps)
		return c.JSON(apps)
	})

	api.Get("/users/:user", controllers.UserController{}.Show)
}
