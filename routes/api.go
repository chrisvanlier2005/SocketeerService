package routes

import (
	"JobHiraMicroservice/controllers"
	controllersapi "JobHiraMicroservice/controllers/api"
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

	api.Get("/connections", func(ctx *fiber.Ctx) error {
		return controllersapi.ConnectionController{Ctx: ctx}.Index()
	})

	api.Get("/connections/:connection", func(ctx *fiber.Ctx) error {
		return controllersapi.ConnectionController{Ctx: ctx}.Show()
	})

	api.Get("/connections/:connection/remove", func(ctx *fiber.Ctx) error {
		return controllersapi.ConnectionController{Ctx: ctx}.Delete()
	})

	api.Get("/users/:user", controllers.UserController{}.Show)
}
