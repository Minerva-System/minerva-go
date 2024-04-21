package schema

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

type NewProduct struct {
	Description string          `json:"description" validate:"required,max=200"`
	Unit        string          `json:"unit" validate:"required,len=2"`
	Price       decimal.Decimal `json:"price" validate:"required"`
}

func (n *NewProduct) ToMessage(companyId string) rpc.Product {
	return rpc.Product{
		CompanyId:   companyId,
		Description: n.Description,
		Unit:        n.Unit,
		Price:       n.Price.String(),
	}
}

type UpdatedProduct struct {
	Description string          `json:"description" validate:"max=200"`
	Unit        string          `json:"unit" validate:"len=2"`
	Price       decimal.Decimal `json:"price" validate:""`
}

func (n *UpdatedProduct) ToMessage(companyId string, id string) rpc.Product {
	return rpc.Product{
		Id:          &id,
		CompanyId:   companyId,
		Description: n.Description,
		Unit:        n.Unit,
		Price:       n.Price.String(),
	}
}
