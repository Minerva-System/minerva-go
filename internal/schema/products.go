package schema

import (
	"github.com/shopspring/decimal"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

type NewProduct struct {
	Description string          `json:"description"`
	Unit        string          `json:"unit"`
	Price       decimal.Decimal `json:"price"`
}

func (n *NewProduct) ToMessage() rpc.Product {
	return rpc.Product{
		Description: n.Description,
		Unit:        n.Unit,
		Price:       n.Price.String(),
	}
}
