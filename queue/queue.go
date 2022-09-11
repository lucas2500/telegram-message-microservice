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
		os.Getenv("RABBITMQ_QUEUE_NAME"),
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

			// Envia mensagem ao Telegram!!
			// Aguarda dois segundos entre cada requisição para evitar HTTP 429 (Too many requests)
			telegram.SendMessageToTelegram(d.Body)
			time.Sleep(2 * time.Second)
		}
	}()

	log.Printf(" [*] Aguardando novas mensagens...")
	<-forever
}

func QueueMessage(conn *amqp.Connection, body []byte) bool {

	exchange := os.Getenv("RABBITMQ_EXCHANGE_NAME")
	RoutingKey := os.Getenv("RABBITMQ_QUEUE_ROUTING_KEY")

	// Abre canal com o broker
	ch, err := conn.Channel()
	FailOnError(err, "Falha ao abrir canal!!")

	// Declara exchange
	SetExchange(ch, exchange)

	// Declara fila
	queue := SetQueue(ch)

	// Realiza o bind da exchange com a fila
	SetQueueBind(ch, queue, exchange, RoutingKey)

	// Publica a mensagem
	err = ch.Publish(
		exchange,
		RoutingKey,
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

func SetExchange(ch *amqp.Channel, exchange string) {

	err := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Falha ao declarar exchange!!")
}

func SetQueue(ch *amqp.Channel) string {

	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Falha ao declarar fila!!")

	return q.Name
}

func SetQueueBind(ch *amqp.Channel, queue string, exchange string, RoutingKey string) {

	err := ch.QueueBind(
		queue,
		RoutingKey,
		exchange,
		false,
		nil,
	)

	FailOnError(err, "Falha ao realizar o bind da exchange com a fila!!")
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
