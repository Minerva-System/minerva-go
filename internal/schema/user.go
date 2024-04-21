package schema

import (
	_ "github.com/go-playground/validator/v10"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

type NewUser struct {
	Login    string `json:"login" validate:"required,min=5,max=25"`
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (n *NewUser) ToMessage(companyId string) rpc.User {
	var email *string = nil
	if n.Email != "" {
		email = &n.Email
	}

	return rpc.User{
		CompanyId: companyId,
		Login:     n.Login,
		Name:      n.Name,
		Email:     email,
		Password:  &n.Password,
	}
}
