package main

import (
	"telegram-message-microservice/connections"
	"telegram-message-microservice/util"
	"telegram-message-microservice/worker"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		util.FailOnError(err, "Falha ao carregar .env")
	}
}

func main() {
	conn := connections.ConnectToRabbitMQ()
	worker.StartConsumer(conn)
}
