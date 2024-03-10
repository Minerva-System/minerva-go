package minerva_connection

import (
	"log/slog"
	
	"gorm.io/gorm"
	rabbitmq "github.com/wagslane/go-rabbitmq"
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

	slog.Info("Initializing connection collection...")

	if(options.WithDatabase) {
		slog.Info("Connecting to database...")
		col.DB, err = databaseConnect()

		if err != nil {
			slog.Error("Error while connecting to database: %v", err)
			return Collection{}, err
		}
	}

	if(options.WithMessageBroker) {
		slog.Info("Connecting to message broker...")
		col.Broker, err = brokerConnect()

		if err != nil {
			slog.Error("Error while connecting to message broker: %v", err)
			return Collection{}, err
		}
	}

	slog.Info("Requested connections established.")
	return col, nil
}
