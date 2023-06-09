package minerva_user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Login string
	Name string
	Email *string
	Pwhash []byte
}
