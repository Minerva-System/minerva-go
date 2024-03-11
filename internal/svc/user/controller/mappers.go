package controller

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	model "github.com/Minerva-System/minerva-go/internal/model"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func MapModelToMessage(m model.User) rpc.User {
	id := m.ID.String()
	
	return rpc.User{
		Id: &id,
		Login: m.Login,
		Name: m.Name,
		Email: m.Email,
		Password: nil, // Never give back a password hash
	}
}

func MapModelListToMessage(l []model.User) *rpc.UserList {
	var result rpc.UserList
	for _, m := range l {
		msg := MapModelToMessage(m)
		result.Users = append(result.Users, &msg)
	}
	return &result
}

func MapMessageToModel(m *rpc.User) (model.User, error) {
	id := uuid.UUID{}
	pwhash := make([]byte, 0)
	var err error
	
	if m.Id != nil {
		id, err = uuid.FromBytes([]byte(*m.Id))
		if err != nil {
			log.Error("Unable to parse UUID from gRPC message to User model: %v", err)
			return model.User{}, err
		}
	}
	
	if m.Password != nil {
		pwhash, err = bcrypt.GenerateFromPassword([]byte(*m.Password), 8)
		if err != nil {
			log.Error("Unable to generate password hash: %v", err)
			return model.User{}, err
		}
	}
	
	return model.User{
		ID: id,
		Login: m.Login,
		Name: m.Name,
		Email: m.Email,
		Pwhash: pwhash,
	}, nil
}
