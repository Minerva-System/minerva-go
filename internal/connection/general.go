package minerva_connection

import (
	"gorm.io/gorm"
	rabbitmq "github.com/wagslane/go-rabbitmq"

	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type CollectionOptions struct {
	WithDatabase bool
	WithMessageBroker bool
}

type Collection struct {
	DB *gorm.DB
	Broker *rabbitmq.Conn
}

func NewCollection(options CollectionOptions) (Collection, error) {
	var err error
	col := Collection{}

	log.Info("Initializing connection collection...")

	if(options.WithDatabase) {
		log.Info("Connecting to database...")
		col.DB, err = databaseConnect()

		if err != nil {
			log.Error("Error while connecting to database: %v", err)
			return Collection{}, err
		}
	}

	if(options.WithMessageBroker) {
		log.Info("Connecting to message broker...")
		col.Broker, err = brokerConnect()

		if err != nil {
			log.Error("Error while connecting to message broker: %v", err)
			return Collection{}, err
		}
	}

	log.Info("Requested connections established.")
	return col, nil
}
