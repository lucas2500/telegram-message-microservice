package worker

import (
	"fmt"
	"log"
	"telegram-message-microservice/connections"
	"telegram-message-microservice/telegram"
	"telegram-message-microservice/util"
)

func StartConsumer(queue string) {

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

	msgs, err := ch.Consume(
		queue,
		"telegram-consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	util.FailOnError(err, "Falha ao registrar consumer!!")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Nova mensagem recebida: %s", d.Body)
			telegram.SendMessageToTelegram(d.Body)
		}
	}()

	fmt.Println("Inicializando worker da fila...", queue)
	log.Printf(" [*] Aguardando novas mensagens...")
	<-forever
}
