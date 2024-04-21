package minerva_connection

import (
	"gorm.io/gorm"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	grpcpool "github.com/processout/grpc-go-pool"

	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type CollectionOptions struct {
	WithDatabase bool
	WithMessageBroker bool
	WithUserService bool
	WithSessionService bool
	WithProductsService bool
	WithTenantService bool
}

type Collection struct {
	DB *gorm.DB
	Broker *rabbitmq.Conn
	UserSvc *grpcpool.Pool
	SessionSvc *grpcpool.Pool
	ProductsSvc *grpcpool.Pool
	TenantSvc *grpcpool.Pool
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

	if(options.WithUserService) {
		log.Info("Connecting to user service...")
		col.UserSvc, err = newGrpcClientPool(GrpcClientKindUser)
		if err != nil {
			log.Error("Error while creating user service pool: %v", err)
			return Collection{}, err
		}
	}

	if(options.WithSessionService) {
		log.Info("Connecting to session service...")
		col.SessionSvc, err = newGrpcClientPool(GrpcClientKindSession)
		if err != nil {
			log.Error("Error while creating session service pool: %v", err)
			return Collection{}, err
		}
	}

	if(options.WithProductsService) {
		log.Info("Connecting to products service...")
		col.ProductsSvc, err = newGrpcClientPool(GrpcClientKindProducts)
		if err != nil {
			log.Error("Error while creating products service pool: %v", err)
			return Collection{}, err
		}
	}

	if(options.WithTenantService) {
		log.Info("Connecting to tenant service...")
		col.TenantSvc, err = newGrpcClientPool(GrpcClientKindTenant)
		if err != nil {
			log.Error("Error while creating tenant service pool: %v", err)
			return Collection{}, err
		}
	}

	log.Info("Requested connections established.")
	return col, nil
}
