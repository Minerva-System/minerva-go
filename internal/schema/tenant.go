package schema

import (
	_ "github.com/go-playground/validator/v10"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

type NewCompany struct {
	Slug string          `json:"slug" validate:"required,min=3,max=30"`
	CompanyName string   `json:"companyName" validate:"required,min=3,max=255"`
	TradingName string   `json:"tradingName" validate:"required,min=3,max=255"`
}

func (n *NewCompany) ToMessage() rpc.Company {
	return rpc.Company{
		Slug: n.Slug,
		CompanyName: n.CompanyName,
		TradingName: n.TradingName,
	}
}

type UpdatedCompany struct {
	Slug string          `json:"slug" validate:"max=30"`
	CompanyName string   `json:"companyName" validate:"max=255"`
	TradingName string   `json:"tradingName" validate:"max=255"`
}

func (n *UpdatedCompany) ToMessage(id string) rpc.Company {
	return rpc.Company{
		Id: &id,
		Slug: n.Slug,
		CompanyName: n.CompanyName,
		TradingName: n.TradingName,
	}
}
