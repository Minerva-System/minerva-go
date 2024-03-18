package minerva_connection

import (
	"os"
	"fmt"
	"errors"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func brokerConnect() (*rabbitmq.Conn, error) {
	var user, pw, server string
	var exists bool

	if user, exists = os.LookupEnv("BROKER_SERVICE_USER"); !exists {
		log.Error("Message broker user not defined")
		return nil, errors.New("BROKER_SERVICE_USER not defined")
	}

	if pw, exists = os.LookupEnv("BROKER_SERVICE_PASSWORD"); !exists {
		log.Error("Message broker password not defined")
		return nil, errors.New("BROKER_SERVICE_PASSWORD not defined")
	}

	if server, exists = os.LookupEnv("BROKER_SERVICE_SERVER"); !exists {
		log.Error("Message broker hostname not defined")
		return nil, errors.New("BROKER_SERVICE_SERVER not defined")
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s/%s", user, pw, server, "%2f")

	log.Info("Connecting to RabbitMQ on %s...", server)
	return rabbitmq.NewConn(dsn, rabbitmq.WithConnectionOptionsLogging)
}
