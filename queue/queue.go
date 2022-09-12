package queue

import (
	"os"
	"telegram-message-microservice/connections"
	"telegram-message-microservice/util"

	"github.com/streadway/amqp"
)

func QueueMessage(body []byte) bool {

	exchange := os.Getenv("RABBITMQ_EXCHANGE_NAME")
	RoutingKey := os.Getenv("RABBITMQ_QUEUE_ROUTING_KEY")

	// Abre conexão com o broker
	conn := connections.ConnectToRabbitMQ()

	// Fecha conexão aberta
	defer conn.Close()

	// Abre canal com o broker
	ch, err := conn.Channel()
	util.FailOnError(err, "Falha ao abrir canal!!")

	// Fecha canal aberto
	defer ch.Close()

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

	if err != nil {
		util.FailOnError(err, "Erro ao publicar mensagem!!")
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

	util.FailOnError(err, "Falha ao declarar exchange!!")
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

	util.FailOnError(err, "Falha ao declarar fila!!")

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

	util.FailOnError(err, "Falha ao realizar o bind da exchange com a fila!!")
}
