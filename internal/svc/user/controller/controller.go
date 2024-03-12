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
)

const PAGESIZE = 100


func UserIndex(db *gorm.DB, page int64) (*rpc.UserList, error) {
	if page < 0 {
		log.Error("User controller: Invalid page index: %d", page)
		return nil, status.Error(codes.OutOfRange, "Invalid user list page")
	}

	list, err := repository.ListUsers(db, int(page) * PAGESIZE, PAGESIZE)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error acessing database: %v", err)
	}

	return MapModelListToMessage(list), nil
}

func GetUser(db *gorm.DB, id string) (*rpc.User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return nil, status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}
	
	usr, err := repository.GetUser(db, parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("User %s not found", id)
			return nil, status.Error(codes.NotFound, "User not found")
		}
		
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := MapModelToMessage(usr)
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
	
	m, err := MapMessageToModel(data)
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

	msg := MapModelToMessage(created)
	return &msg, nil
}

func DeleteUser(db *gorm.DB, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	exists, err := repository.ExistsUser(db, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	if !exists {
		log.Error("User %s not found", id)
		return status.Error(codes.NotFound, "User not found")
	}
	
	err = repository.DeleteUser(db, parsedId)
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
	
	d, err := MapMessageToModel(data)
	if err != nil {
		log.Error("Error mapping user data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping user data: %v", err)
	}

	result, err := repository.UpdateUser(db, d)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := MapModelToMessage(result)
	return &msg, nil
}
