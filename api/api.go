package main

import (
	"log"
	"os"
	"telegram-message-microservice/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal(err, "Falha ao carregar .env!!")
	}

}

func main() {

	app := fiber.New()
	routes.SetupRoutes(app)

	app.Listen(":" + os.Getenv("API_HTTP_PORT"))
}
