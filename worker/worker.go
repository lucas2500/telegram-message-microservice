package worker

import (
	"telegram-message-microservice/queue"
	"telegram-message-microservice/telegram"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(QueueName string) {

	message := make(chan amqp.Delivery)

	go func() {
		queue.DequeueMessage(QueueName, message)
	}()

	for msg := range message {
		telegram.SendMessageToTelegram(msg.Body)
		msg.Ack(true)
	}
}
