package schema

import (
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	// log "github.com/Minerva-System/minerva-go/pkg/log"
	// model "github.com/Minerva-System/minerva-go/internal/model"
)

type NewUser struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (n *NewUser) ToMessage() rpc.User {
	var email *string = nil
	if n.Email != "" {
		email = &n.Email
	}

	return rpc.User{
		Login:    n.Login,
		Name:     n.Name,
		Email:    email,
		Password: &n.Password,
	}
}
