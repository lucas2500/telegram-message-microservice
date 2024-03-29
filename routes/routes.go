package routes

import (
	"telegram-message-microservice/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", controllers.PingPong)
	app.Post("api/SendMessage", controllers.SendMessage)
}
