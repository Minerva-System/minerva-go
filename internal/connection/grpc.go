package minerva_connection

import (
	"errors"
	"fmt"
	"os"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"

	log "github.com/Minerva-System/minerva-go/pkg/log"
)

const (
	GrpcClientKindUser     string = "USER"
	GrpcClientKindSession  string = "SESSION"
	GrpcClientKindProducts string = "PRODUCTS"
	GrpcClientKindTenant   string = "TENANT"
)

func newGrpcClientPool(clientKind string) (*grpcpool.Pool, error) {
	var host string
	var exists bool

	varname := fmt.Sprintf("MINERVA_%s_HOST", clientKind)
	log.Debug("Host env variable: %s", varname)
	if host, exists = os.LookupEnv(varname); !exists {
		log.Error("Host for %s service not defined", clientKind)
		return nil, errors.New(fmt.Sprintf("%s not defined", host))
	}

	var factory grpcpool.Factory = func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error("Unable to connect to %s (%s): %v", clientKind, host, err)
			return nil, err
		}
		log.Info(
			"Connection to %s (%s) created",
			clientKind,
			host,
		)
		return conn, err
	}

	return grpcpool.New(factory, 1, 5, time.Second)
}
