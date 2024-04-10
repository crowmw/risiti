package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"email" json:"email" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password string `gorm:"password" json:"password" validate:"min=8"`
}
