package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	util "github.com/Minerva-System/minerva-go/pkg/util"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:UUID();primaryKey" json:"id"`
	CompanyID   uuid.UUID       `gorm:"type:uuid;primaryKey" json:"-"`
	Description string          `gorm:"type:varchar(200);not null" json:"description"`
	Unit        string          `gorm:"type:char(2);not null" json:"unit"`
	Price       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"price"`
	CreatedAt   time.Time       `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time       `gorm:"not null" json:"updatedAt"`
	Company     Company         `gorm:"foreignKey:CompanyID;references:ID" json:"-"`
}

func (p *Product) ToMessage() rpc.Product {
	id := p.ID.String()

	return rpc.Product{
		Id:          &id,
		CompanyId:   p.CompanyID.String(),
		Description: p.Description,
		Unit:        util.StringToUnit(p.Unit),
		Price:       p.Price.String(),
	}
}

func (Product) ListToMessage(l []Product) *rpc.ProductList {
	var result rpc.ProductList
	for _, p := range l {
		msg := p.ToMessage()
		result.Products = append(result.Products, &msg)
	}
	return &result
}

func (Product) FromMessage(p *rpc.Product) (Product, error) {
	id := uuid.UUID{}
	var err error

	if p.Id != nil {
		log.Debug("Parsing UUID: %s", *p.Id)
		id, err = uuid.Parse(*p.Id)
		if err != nil {
			log.Error("Unable to parse UUID from gRPC message to Product model: %v", err)
			return Product{}, err
		}
	}

	log.Debug("Parsing company UUID: %s", p.CompanyId)
	companyId, err := uuid.Parse(p.CompanyId)
	if err != nil {
		log.Error("Unable to parse company UUID from gRPC message to Product model: %v", err)
		return Product{}, err
	}

	log.Debug("Parsing price: %s", p.Price)
	price, err := decimal.NewFromString(p.Price)
	if err != nil {
		log.Error("Unable to parse decimal %s: %v", p.Price, err)
		return Product{}, err
	}

	return Product{
		ID:          id,
		CompanyID:   companyId,
		Description: p.Description,
		Unit:        util.StringToUnit(p.Unit),
		Price:       price,
	}, nil
}

func (Product) FromListMessage(xs *rpc.ProductList) ([]Product, error) {
	result := make([]Product, 0)
	if xs != nil {
		for _, x := range xs.Products {
			p, err := Product{}.FromMessage(x)
			if err != nil {
				return make([]Product, 0), err
			}
			result = append(result, p)
		}
	}
	return result, nil
}
