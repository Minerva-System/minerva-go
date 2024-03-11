package model

import (
	"time"
	
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID         `gorm:"type:uuid;default:UUID()" json:"id"`
	Login string         `json:"login" gorm:"unique"`
	Name string          `json:"name" gorm:"not null"`
	Email *string        `json:"email,omitempty"`
	Pwhash []byte        `json:"-" gorm:"not null"`
	CreatedAt time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"not null"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}
