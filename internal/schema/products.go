package schema

import (
	"github.com/shopspring/decimal"
	_ "github.com/go-playground/validator/v10"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

type NewProduct struct {
	Description string          `json:"description" validate:"required"`
	Unit        string          `json:"unit" validate:"required,len=2"`
	Price       decimal.Decimal `json:"price" validate:"required"`
}

func (n *NewProduct) ToMessage() rpc.Product {
	return rpc.Product{
		Description: n.Description,
		Unit:        n.Unit,
		Price:       n.Price.String(),
	}
}
