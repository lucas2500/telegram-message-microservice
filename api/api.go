package main

import (
	"log"
	"telegram-message-microservice/entity"
	"telegram-message-microservice/queue"

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
	app.Get("/", Hello)
	SetupRoutes(app)

	app.Listen(":3001")
}

func SetupRoutes(app *fiber.App) {
	app.Post("api/SendMessage", SendMessage)
}

func Hello(c *fiber.Ctx) error {

	var res entity.Response
	res.Result = "Hello!!"

	return c.JSON(res)
}

func SendMessage(c *fiber.Ctx) error {

	ch := queue.Connect()
	var res entity.Response

	message := new(entity.Message)

	if err := c.BodyParser(message); err != nil {
		res.Result = "Erro no parsing do json!!"
		return c.Status(400).JSON(res)
	}

	if !queue.QueueMessage(ch, string(c.Body())) {
		res.Result = "Houve eum erro ao inserir a mensagem na fila!!"
		return c.Status(500).JSON(res)

	}

	res.Result = "Mensagem incluida na fila com sucesso!!"
	return c.Status(201).JSON(res)
}
