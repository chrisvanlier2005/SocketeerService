package controllers

import (
	"JobHiraMicroservice/models"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (u UserController) Index(c *fiber.Ctx) error {
	var users []models.User
	models.DB.Find(&users)
	return c.JSON(users)
}

func (u UserController) Show(c *fiber.Ctx) error {
	id := c.Params("user")
	var user models.User
	models.DB.Find(&user, id)
	return c.JSON(user)
}
