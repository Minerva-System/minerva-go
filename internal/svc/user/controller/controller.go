package controller

import (
	"errors"
	"encoding/json"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	repository "github.com/Minerva-System/minerva-go/internal/svc/user/repository"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

const PAGESIZE = 100


func UserIndex(db *gorm.DB, companyId string, page int64) (*rpc.UserList, error) {
	if page < 0 {
		log.Error("Invalid page index: %d", page)
		return nil, status.Error(codes.OutOfRange, "Invalid user list page")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return nil, status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}

	list, err := repository.ListUsers(db, parsedCompanyId, int(page) * PAGESIZE, PAGESIZE)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error acessing database: %v", err)
	}

	return model.User{}.ListToMessage(list), nil
}

func GetUser(db *gorm.DB, companyId string, id string) (*rpc.User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return nil, status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return nil, status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}
	
	usr, err := repository.GetUser(db, parsedCompanyId, parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("User %s not found", id)
			return nil, status.Error(codes.NotFound, "User not found")
		}
		
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := usr.ToMessage()
	return &msg, nil
}

func CreateUser(db *gorm.DB, data *rpc.User) (*rpc.User, error) {
	if data.Id != nil {
		log.Error("New user must not have a predefined id")
		return nil, status.Error(codes.InvalidArgument, "New user must not have a predefined id")
	}
	
	if data.Password == nil {
		log.Error("Password for new user cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "Password for new user cannot be empty")
	}
	
	m, err := model.User{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping user data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping user data: %v", err)
	}

	log.Debug("Accessing database...")
	created, err := repository.CreateUser(db, m)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	go func() {
		if json, err := json.Marshal(created); err == nil {
			log.Debug("User created: %s", json)
		}
	}()

	msg := created.ToMessage()
	return &msg, nil
}

func DeleteUser(db *gorm.DB, companyId string, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}

	exists, err := repository.ExistsUser(db, parsedCompanyId, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	if !exists {
		log.Error("User %s not found", id)
		return status.Error(codes.NotFound, "User not found")
	}
	
	err = repository.DeleteUser(db, parsedCompanyId, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}
	return nil
}

func UpdateUser(db *gorm.DB, data *rpc.User) (*rpc.User, error) {
	if data.Id == nil {
		log.Error("User id is missing")
		return nil, status.Error(codes.InvalidArgument, "User id is missing")
	}

	if data.Login != "" {
		log.Error("User login cannot be changed")
		return nil, status.Error(codes.InvalidArgument, "User login cannot be changed")
	}
	
	d, err := model.User{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping user data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping user data: %v", err)
	}

	result, err := repository.UpdateUser(db, d)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := result.ToMessage()
	return &msg, nil
}
