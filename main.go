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
	Queue := os.Getenv("RABBITMQ_MESSAGE_QUEUE")
	conn := connections.ConnectToRabbitMQ()

	defer conn.Close()

	wg.Add(1)
	go worker.StartConsumer(conn, Queue)
	wg.Wait()
}
