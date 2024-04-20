package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

func GetCompany(db *gorm.DB, id uuid.UUID) (model.Company, error) {
	var company model.Company
	result := db.First(&company, "id = ?", id)
	return company, result.Error
}

func GetCompanyBySlug(db *gorm.DB, slug string) (model.Company, error) {
	var company model.Company
	result := db.First(&company, "slug = ?", slug)
	return company, result.Error
}

func ListCompanies(db *gorm.DB, offset int, limit int) ([]model.Company, error) {
	var companies []model.Company
	result := db.Find(&companies).
		Limit(limit).
		Offset(offset)
	return companies, result.Error
}

func CreateCompany(db *gorm.DB, data model.Company) (model.Company, error) {
	result := db.Create(&data)
	if result.Error != nil {
		log.Error("Unable to create company: %v", result.Error)
		return model.Company{}, result.Error
	}
	return data, nil
}

func DisableCompany(db *gorm.DB, id uuid.UUID) error {
	now := time.Now()
	result := db.Model(&model.Company{}).
		Where("id = ?", id).
		Updates(model.Company{
			DeletedAt: &now,
		})
	return result.Error
}

func ExistsCompany(db *gorm.DB, id uuid.UUID) (bool, error) {
	var exists bool = false
	result := db.Model(&model.Company{}).
		Select("COUNT(*) > 0").
		Where("id = ?", id).
		Find(&exists)
	return exists, result.Error
}

func ExistsCompanyBySlug(db *gorm.DB, slug string) (bool, error) {
	var exists bool = false
	result := db.Model(&model.Company{}).
		Select("COUNT(*) > 0").
		Where("slug = ?", slug).
		Find(&exists)
	return exists, result.Error
}


func UpdateCompany(db *gorm.DB, data model.Company) (model.Company, error) {
	result := db.Model(&data).
		Updates(model.Company{
			ID: data.ID,
			Slug: data.Slug,
			CompanyName: data.CompanyName,
			TradingName: data.TradingName,
		})
	if result.Error != nil {
		return model.Company{}, result.Error
	}
	return GetCompany(db, data.ID)
}
