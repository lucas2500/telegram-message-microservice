package main

import (
	"os"
	"sync"
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

	var wg sync.WaitGroup

	conn := connections.ConnectToRabbitMQ()

	defer conn.Close()

	Queue := os.Getenv("RABBITMQ_MESSAGE_QUEUE")

	wg.Add(1)

	go worker.StartConsumer(Queue)

	wg.Wait()
}
