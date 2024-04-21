package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

func GetProduct(db *gorm.DB, companyId uuid.UUID, id uuid.UUID) (model.Product, error) {
	var product model.Product
	result := db.
		Where("company_id = ?", companyId).
		First(&product, "id = ?", id)
	return product, result.Error
}

func ListProducts(db *gorm.DB, companyId uuid.UUID, offset int, limit int) ([]model.Product, error) {
	var products []model.Product
	result := db.Find(&products).
		Where("company_id = ?", companyId).
		Limit(limit).
		Offset(offset)
	return products, result.Error
}

func CreateProduct(db *gorm.DB, data model.Product) (model.Product, error) {
	result := db.Create(&data)
	if result.Error != nil {
		log.Error("Unable to create product: %v", result.Error)
		return model.Product{},result.Error
	}

	return data, nil
}

func DeleteProduct(db *gorm.DB, companyId uuid.UUID, id uuid.UUID) error {
	return db.Delete(&model.Product{}, "id = ? AND company_id = ?", id, companyId).Error
}

func ExistsProduct(db *gorm.DB, companyId uuid.UUID, id uuid.UUID) (bool, error) {
	var exists bool = false
	result := db.Model(&model.Product{}).
		Select("COUNT(*) > 0").
		Where("id = ? AND company_id = ?", id, companyId).
		Find(&exists)
	return exists, result.Error
}

// func UpdateProduct(db *gorm.DB, data model.Product) (model.Product, error)

