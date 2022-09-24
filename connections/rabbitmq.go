package connections

import (
	"fmt"
	"os"
	"telegram-message-microservice/util"
	"time"

	"github.com/streadway/amqp"
)

func ConnectToRabbitMQ() *amqp.Connection {

	var (
		counter int = 0
		dsn     string
		conn    *amqp.Connection
		err     error
	)

	for {

		counter++

		fmt.Println("Tentativa", counter, "de conexão!!")

		dsn = "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
		conn, err = amqp.Dial(dsn)

		if err == nil || counter == 10 {
			break
		}

		time.Sleep(5 * time.Second)
	}

	if err != nil {
		util.FailOnError(err, "Falha ao se conectar ao RBMQ")
	}

	return conn
}
