package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

func GetUser(db *gorm.DB, id uuid.UUID) (model.User, error) {
	var user model.User
	result := db.First(&user, "id = ?", id)
	return user, result.Error
}

func ListUsers(db *gorm.DB, offset int, limit int) ([]model.User, error) {
	var users []model.User
	result := db.Find(&users).
		Limit(limit).
		Offset(offset)
	return users, result.Error
}

func CreateUser(db *gorm.DB, data model.User) (model.User, error) {
	result := db.Create(&data)
	if result.Error != nil {
		log.Error("Unable to create user: %v", result.Error)
		return model.User{}, result.Error
	}

	return data, nil
}

func DeleteUser(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&model.User{}, "id = ?", id).Error
}

func ExistsUser(db *gorm.DB, id uuid.UUID) (bool, error) {
	var exists bool = false
	result := db.Model(&model.User{}).
		Select("COUNT(*) > 0").
		Where("ID = ?", id).
		Find(&exists)
	return exists, result.Error
}

func UpdateUser(db *gorm.DB, data model.User) (model.User, error) {
	if result := db.Model(&data).
		Updates(model.User{
			ID: data.ID,
			Name: data.Name,
			Pwhash: data.Pwhash,
		}); result.Error != nil {
		return model.User{}, result.Error
	}
	// return GetUser(db, data.ID) // TODO: Is this really necessary?
	return data, nil
}
