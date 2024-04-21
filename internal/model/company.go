package model

import (
	"time"

	"github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	util "github.com/Minerva-System/minerva-go/pkg/util"
)

type Company struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:UUID()" json:"id"`
	Slug        string         `gorm:"unique" json:"slug"`
	CompanyName string         `gorm:"not null" json:"companyName"`
	TradingName string         `gorm:"not null" json:"tradingName"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (m *Company) ToMessage() rpc.Company {
	id := m.ID.String()

	return rpc.Company{
		Id:          &id,
		Slug:        m.Slug,
		CompanyName: m.CompanyName,
		TradingName: m.TradingName,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

func (Company) ListToMessage(l []Company) *rpc.CompanyList {
	var result rpc.CompanyList
	for _, m := range l {
		msg := m.ToMessage()
		result.Companies = append(result.Companies, &msg)
	}
	return &result
}

func (Company) FromMessage(m *rpc.Company) (Company, error) {
	id := uuid.UUID{}
	var err error

	if m.Id != nil {
		log.Debug("Parsing UUID: %s", *m.Id)
		id, err = uuid.Parse(*m.Id)
		if err != nil {
			log.Error("Unable to parse UUID from gRPC message to Company model: %v", err)
			return Company{}, err
		}
	}

	slug, err := util.HygienizeSlug(m.Slug)
	if err != nil {
		log.Error("Error hygienizing company slug: %v", err)
		return Company{}, err
	}

	return Company{
		ID:          id,
		Slug:        slug,
		CompanyName: m.CompanyName,
		TradingName: m.TradingName,
		CreatedAt:   m.CreatedAt.AsTime(),
		UpdatedAt:   m.UpdatedAt.AsTime(),
	}, nil
}

func (Company) FromListMessage(xs *rpc.CompanyList) ([]Company, error) {
	result := make([]Company, 0)
	if xs != nil {
		for _, x := range xs.Companies {
			m, err := Company{}.FromMessage(x)
			if err != nil {
				return make([]Company, 0), err
			}
			result = append(result, m)
		}
	}
	return result, nil
}
