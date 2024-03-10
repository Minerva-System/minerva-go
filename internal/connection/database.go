package minerva_connection

import (
	"os"
	"fmt"
	"errors"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	slogGorm "github.com/orandin/slog-gorm"

	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func databaseConnect() (*gorm.DB, error) {
	var user, pw, server, dbname string
	var exists bool

	if user, exists = os.LookupEnv("DATABASE_SERVICE_USER"); !exists {
		log.Error("Database user not defined")
		return nil, errors.New("DATABASE_SERVICE_USER not defined")
	}

	if pw, exists = os.LookupEnv("DATABASE_SERVICE_PASSWORD"); !exists {
		log.Error("Database password not defined")
		return nil, errors.New("DATABASE_SERVICE_PASSWORD not defined")
	}

	if server, exists = os.LookupEnv("DATABASE_SERVICE_SERVER"); !exists {
		log.Error("Database hostname not defined")
		return nil, errors.New("DATABASE_SERVICE_SERVER not defined")
	}

	if dbname, exists = os.LookupEnv("DATABASE_SERVICE_DBNAME"); !exists {
		log.Error("Database name not defined")
		return nil, errors.New("DATABASE_SERVICE_DBNAME not defined")
	}
	
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		pw,
		server,
		dbname,
	)

	log.Info("Connecting to MySQL on %s...", server)
	return gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
			DefaultStringSize: 256,
			DontSupportRenameIndex: true,
			DontSupportRenameColumn: true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{
			Logger: slogGorm.New(),
		})
}
