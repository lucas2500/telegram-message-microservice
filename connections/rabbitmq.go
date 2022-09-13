package connections

import (
	"os"
	"telegram-message-microservice/util"

	"github.com/streadway/amqp"
)

func ConnectToRabbitMQ() *amqp.Connection {

	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	conn, err := amqp.Dial(dsn)

	if err != nil {
		util.FailOnError(err, "Falha ao se conectar ao RBMQ")
	}

	return conn
}
