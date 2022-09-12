package controllers

import (
	"telegram-message-microservice/entities"
	"telegram-message-microservice/queue"

	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {

	response := make(map[string]string)
	response["Ping"] = "Pong"

	return c.JSON(response)
}

func SendMessage(c *fiber.Ctx) error {

	response := make(map[string]string)

	message := new(entities.Message)

	if err := c.BodyParser(message); err != nil {
		response["result"] = "Erro no parsing do json!!"
		return c.Status(400).JSON(response)
	}

	if !queue.QueueMessage(c.Body()) {
		response["result"] = "Houve eum erro ao inserir a mensagem na fila!!"
		return c.Status(500).JSON(response)
	}

	response["result"] = "Mensagem incluida na fila com sucesso!!"
	return c.Status(201).JSON(response)
}
