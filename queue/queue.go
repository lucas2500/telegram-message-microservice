package queue

import (
	"log"
	"os"
	"telegram-message-microservice/telegram"
	"time"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Connection {

	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	conn, err := amqp.Dial(dsn)

	FailOnError(err, "Falha ao se conectar ao RBMQ!!")

	return conn
}

func StartConsumer(conn *amqp.Connection) {

	ch, ErrChan := conn.Channel()
	FailOnError(ErrChan, "Falha ao abrir canal!!")

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

	defer ch.Close()

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

func QueueMessage(conn *amqp.Connection, body []byte) bool {

	ch, ErrChan := conn.Channel()
	FailOnError(ErrChan, "Falha ao abrir canal!!")

	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION"),
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	defer ch.Close()

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
