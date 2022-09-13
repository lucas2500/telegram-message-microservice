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
	MessageQueue := os.Getenv("RABBITMQ_QUEUE_NAME")
	MessageDeadLetterQueue := os.Getenv("RABBITMQ_DLQ_NAME")

	conn := connections.ConnectToRabbitMQ()

	defer conn.Close()

	wg.Add(2)
	go worker.StartConsumer(conn, MessageQueue)
	go worker.StartConsumer(conn, MessageDeadLetterQueue)
	wg.Wait()
}
