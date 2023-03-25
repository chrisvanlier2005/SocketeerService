package main

import (
	"JobHiraMicroservice/routes"
	"JobHiraMicroservice/websockets"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	/*err := models.DB.AutoMigrate(models.User{}, models.Application{})
	if err != nil {
		return
	}*/
	app := fiber.New()

	routes.WebRoutes(app.Group("/"))
	routes.WebSockets(app.Group("/ws"))
	routes.ApiRoutes(app.Group("/api"))

	go websockets.RunHub()

	log.Fatal(app.Listen(":8081"))
}
