package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})

	//run database
	configs.ConnectDB()

	//routes
	routes.AlunosRoute(app)

	app.Listen(":6000")
}
