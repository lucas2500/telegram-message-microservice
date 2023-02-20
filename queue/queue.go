package queue

import (
	"context"
	"log"
	"telegram-message-microservice/connections"
	"telegram-message-microservice/entities"
	"telegram-message-microservice/util"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func QueueMessage(body []byte, qp entities.QueueProperties) bool {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	defer cancel()

	// Obtem conexão aberta com o RabbitMQ
	conn := connections.RabbitConn

	// Abre canal com o broker
	ch, err := conn.Channel()
	util.FailOnError(err, "Falha ao abrir canal!!")

	// Fecha canal aberto
	defer ch.Close()

	// Declara exchange
	SetExchange(ch, qp.Exchange)

	// Declara fila
	SetQueue(ch, qp.Queue, qp.DLX)

	// Realiza o bind da exchange com a fila
	SetQueueBind(ch, qp.Queue, qp.Exchange, qp.RoutingKey)

	// Publica a mensagem
	err = ch.PublishWithContext(ctx,
		qp.Exchange,
		qp.RoutingKey,
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

func SetQueue(ch *amqp.Channel, queue string, DLX map[string]interface{}) string {

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		DLX,
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

func DequeueMessage(queue string, message chan amqp.Delivery, WorkerId int) {

	// Obtem conexão aberta com o RabbitMQ
	conn := connections.RabbitConn

	ch, ErrChan := conn.Channel()

	util.FailOnError(ErrChan, "Falha ao abrir canal!!")

	defer ch.Close()

	_, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao declarar fila!!")

	// Define a quantidade de mensagens que serão enviadas para o worker
	// antes da confirmação do recebimento
	ch.Qos(
		1,
		0,
		true,
	)

	msgs, err := ch.Consume(
		queue,
		"telegram-consumer",
		false,
		false,
		false,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao registrar consumer!!")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println("Worker", WorkerId, "consumindo mensagem")
			message <- d
		}
	}()

	<-forever
}
