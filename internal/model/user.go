package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:UUID()" json:"id"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Login     string    `json:"login" gorm:"unique"`
	Name      string    `json:"name" gorm:"not null"`
	Email     *string   `json:"email,omitempty"`
	Pwhash    []byte    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Company   Company   `gorm:"foreignKey:CompanyID;references:ID" json:"-"`
}

func (m *User) ToMessage() rpc.User {
	id := m.ID.String()

	return rpc.User{
		Id:       &id,
		Login:    m.Login,
		Name:     m.Name,
		Email:    m.Email,
		Password: nil, // Never give back a password hash
	}
}

func (User) ListToMessage(l []User) *rpc.UserList {
	var result rpc.UserList
	for _, m := range l {
		msg := m.ToMessage()
		result.Users = append(result.Users, &msg)
	}
	return &result
}

func (User) FromMessage(m *rpc.User) (User, error) {
	id := uuid.UUID{}
	pwhash := make([]byte, 0)
	var err error

	if m.Id != nil {
		log.Debug("Parsing UUID: %s", *m.Id)
		id, err = uuid.Parse(*m.Id)
		if err != nil {
			log.Error("Unable to parse UUID from gRPC message to User model: %v", err)
			return User{}, err
		}
	}

	if m.Password != nil {
		pwhash, err = bcrypt.GenerateFromPassword([]byte(*m.Password), 8)
		if err != nil {
			log.Error("Unable to generate password hash: %v", err)
			return User{}, err
		}
	}

	return User{
		ID:     id,
		Login:  m.Login,
		Name:   m.Name,
		Email:  m.Email,
		Pwhash: pwhash,
	}, nil
}

func (User) FromListMessage(xs *rpc.UserList) ([]User, error) {
	result := make([]User, 0)
	if xs != nil {
		for _, x := range xs.Users {
			m, err := User{}.FromMessage(x)
			if err != nil {
				return make([]User, 0), err
			}
			result = append(result, m)
		}
	}
	return result, nil
}
