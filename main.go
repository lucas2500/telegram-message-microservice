package main

import (
	"log"
	"telegram-message-microservice/queue"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err, "Falha ao carregar .env!!")
	}
}

func main() {
	conn := queue.Connect()
	queue.StartConsumer(conn)
}
