package main

import (
	"os"
	"telegram-message-microservice/connections"
	"telegram-message-microservice/routes"
	"telegram-message-microservice/util"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load("../.env")

	if err != nil {
		util.FailOnError(err, "Falha ao carregar .env")
	}

	connections.ConnectToRabbitMQ()
}

func main() {

	defer connections.RabbitConn.Close()

	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":" + os.Getenv("API_HTTP_PORT"))
}
