package queue

import (
	"context"
	"log"
	"telegram-message-microservice/connections"
	"telegram-message-microservice/util"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueProperties struct {
	Exchange   string
	RoutingKey string
	Queue      string
	DLX        map[string]interface{}
	Body       []byte
}

type QueueConsumer struct {
	Queue          string
	MessageChannel chan amqp.Delivery
}

func (q QueueProperties) QueueMessage() bool {

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
	q.SetExchange(ch)

	// Declara fila
	q.SetQueue(ch)

	// Realiza o bind da exchange com a fila
	q.SetQueueBind(ch)

	// Publica a mensagem
	err = ch.PublishWithContext(ctx,
		q.Exchange,
		q.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        q.Body,
		},
	)

	if err != nil {
		util.FailOnError(err, "Erro ao publicar mensagem!!")
		return false
	}

	return true
}

func (q QueueProperties) SetExchange(ch *amqp.Channel) {

	err := ch.ExchangeDeclare(
		q.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao declarar exchange!!")
}

func (q QueueProperties) SetQueue(ch *amqp.Channel) string {

	queue, err := ch.QueueDeclare(
		q.Queue,
		true,
		false,
		false,
		false,
		q.DLX,
	)

	util.FailOnError(err, "Falha ao declarar fila!!")

	return queue.Name
}

func (q QueueProperties) SetQueueBind(ch *amqp.Channel) {

	err := ch.QueueBind(
		q.Queue,
		q.RoutingKey,
		q.Exchange,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao realizar o bind da exchange com a fila!!")
}

func (q QueueConsumer) DequeueMessage(WorkerId int) {

	// Obtem conexão aberta com o RabbitMQ
	conn := connections.RabbitConn

	ch, ErrChan := conn.Channel()

	util.FailOnError(ErrChan, "Falha ao abrir canal!!")

	defer ch.Close()

	_, err := ch.QueueDeclare(
		q.Queue,
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
		q.Queue,
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
			q.MessageChannel <- d
		}
	}()

	<-forever
}
