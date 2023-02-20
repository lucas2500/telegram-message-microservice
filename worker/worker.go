package worker

import (
	"log"
	"os"
	"strconv"
	"telegram-message-microservice/queue"
	"telegram-message-microservice/telegram"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(QueueName string) {

	ConsumerProps := queue.QueueConsumer{
		Queue:          QueueName,
		MessageChannel: make(chan amqp.Delivery),
	}

	WorkersNumber, err := strconv.Atoi(os.Getenv("WORKERS_NUMBER"))

	if err != nil {
		log.Fatal("Erro ao carregar configuração do worker")
	}

	if WorkersNumber <= 0 {
		WorkersNumber = 1
	}

	for i := 0; i < WorkersNumber; i++ {

		WorkerId := i

		log.Println("Worker", i, "up and running - Queue", QueueName)

		go func() {
			ConsumerProps.DequeueMessage(WorkerId)
		}()
	}

	for msg := range ConsumerProps.MessageChannel {
		telegram.SendMessageToTelegram(msg.Body)
		msg.Ack(true)
	}
}
