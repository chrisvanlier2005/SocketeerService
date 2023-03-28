package controllersapi

import (
	"JobHiraMicroservice/models"
	"JobHiraMicroservice/websockets"
	"github.com/gofiber/fiber/v2"
)

type ConnectionController struct {
	Ctx *fiber.Ctx
}

func (controller ConnectionController) Index() error {
	serverKey := controller.Ctx.Query("key")
	if serverKey == "" {
		return controller.Ctx.JSON(map[string]string{
			"code":     "error",
			"message:": "server key and message is required",
		})
	}
	application := models.Application{}
	if err := models.DB.Where("server_key = ?", serverKey).First(&application).Error; err != nil {
		return controller.Ctx.JSON(map[string]string{
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
	return controller.Ctx.JSON(connections)
}

func (controller ConnectionController) Delete() error {
	connection := controller.Ctx.Params("connection")
	if connection == "" {
		return controller.Ctx.JSON(map[string]string{
			"code":    "error",
			"message": "connection id is required",
		})
	}

	return controller.Ctx.JSON(map[string]string{
		"code":    "success",
		"message": "Succesfully removed " + connection + " from the connected users",
	})
}

func (controller ConnectionController) Show() error {
	return controller.Ctx.JSON(map[string]string{
		"code": "not implemented",
	})
}
