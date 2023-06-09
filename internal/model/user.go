package minerva_user

import (
	"log"
	"time"
	
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	rpc "minervarpc"
)

type User struct {
	ID uuid.UUID         `gorm:"type:uuid;default:UUID()"`
	Login string
	Name string
	Email *string
	Pwhash []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (u *User) ToMessage() rpc.User {
	id := u.ID.String()
	
	return rpc.User{
		Id: &id,
		Login: u.Login,
		Name: u.Name,
		Email: u.Email,
		Password: nil,
	}
}

func UserFromMessage(u *rpc.User) (User, error) {
	id := uuid.UUID{}
	var err error
	
	if u.Id != nil {
		id, err = uuid.FromBytes([]byte(*u.Id))
		if err != nil {
			log.Printf("Error while converting User message to model: %v", err)
			return User{}, err
		}
	}

	hash := make([]byte, 0)

	if u.Password != nil {
		hash, err = bcrypt.GenerateFromPassword([]byte(*u.Password), 8)
		if err != nil {
			log.Printf("Unable to generate password hash: %v", err)
			return User{}, err
		}
	}
	
	return User{
		ID: id,
		Login: u.Login,
		Name: u.Name,
		Email: u.Email,
		Pwhash: hash,
	}, nil
}
