package controllers

import (
	"os"
	"telegram-message-microservice/entities"
	"telegram-message-microservice/queue"

	"github.com/gofiber/fiber/v2"
)

func PingPong(c *fiber.Ctx) error {

	response := map[string]string{"Ping": "Pong"}
	return c.JSON(response)
}

func SendMessage(c *fiber.Ctx) error {

	message := new(entities.Message)

	if err := c.BodyParser(message); err != nil {
		response := map[string]string{"result": "Erro no parsing do json!!"}
		return c.Status(400).JSON(response)
	}

	QueueProps := queue.QueueProperties{
		Exchange:   os.Getenv("RABBITMQ_MESSAGE_EXCHANGE"),
		RoutingKey: os.Getenv("RABBITMQ_MESSAGE_QUEUE_ROUTING_KEY"),
		Queue:      os.Getenv("RABBITMQ_MESSAGE_QUEUE"),
		Body:       c.Body(),
	}

	if !QueueProps.QueueMessage() {
		response := map[string]string{"result": "Houve eum erro ao inserir a mensagem na fila!!"}
		return c.Status(500).JSON(response)
	}

	response := map[string]string{"result": "Mensagem incluida na fila com sucesso!!"}
	return c.Status(201).JSON(response)
}
