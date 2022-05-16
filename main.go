package main

import (
	"log"
	"net/http"
	"os"
	"telegram-message-microservice/queue"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err, "Falha ao carregar .env!!")
	}

	http.Get(os.Getenv("TELEGRAM_BASE_URL"))
}

func main() {
	ch := queue.Connect()
	queue.StartConsumer(ch)
}
