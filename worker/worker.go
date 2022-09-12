package worker

import (
	"log"
	"os"
	"telegram-message-microservice/telegram"
	"telegram-message-microservice/util"
	"time"

	"github.com/streadway/amqp"
)

func StartConsumer(conn *amqp.Connection) {

	ch, ErrChan := conn.Channel()
	util.FailOnError(ErrChan, "Falha ao abrir canal!!")

	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao declarar fila!!")

	msgs, err := ch.Consume(
		q.Name,
		"telegram-consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	defer ch.Close()

	util.FailOnError(err, "Falha ao registrar consumer!!")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Nova mensagem recebida: %s", d.Body)

			// Envia mensagem ao Telegram!!
			// Aguarda dois segundos entre cada requisição para evitar HTTP 429 (Too many requests)
			telegram.SendMessageToTelegram(d.Body)
			time.Sleep(2 * time.Second)
		}
	}()

	log.Printf(" [*] Aguardando novas mensagens...")
	<-forever
}
