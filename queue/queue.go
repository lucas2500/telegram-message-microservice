package queue

import (
	"log"
	"os"
	"telegram-message-microservice/telegram"
	"time"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Channel {

	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	conn, err := amqp.Dial(dsn)

	FailOnError(err, "Falha ao se conectar ao RBMQ!!")

	ch, err := conn.Channel()
	FailOnError(err, "Falha ao abrir canal!!")

	return ch
}

func StartConsumer(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Falha ao declarar fila!!")

	msgs, err := ch.Consume(
		q.Name,
		"telegram-consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Falha ao registrar consumer!!")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Nova mensagem recebida: %s", d.Body)
			telegram.SendMessageToTelegram(d.Body)
			time.Sleep(time.Second * 2)
		}
	}()

	log.Printf(" [*] Aguardando novas mensagens...")
	<-forever
}

func QueueMessage(ch *amqp.Channel, body string) bool {

	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION"),
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)

	if err != nil {
		FailOnError(err, "Erro ao publicar mensagem!!")
		return false
	}

	return true
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
