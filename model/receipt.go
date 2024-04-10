package model

import (
	"time"

	"gorm.io/gorm"
)

type Receipt struct {
	gorm.Model
	Name        string    `gorm:"name" json:"name" validate:"nonzero"`
	Description string    `gorm:"description" json:"filename"`
	Filename    string    `gorm:"filename" json:"date" validate:"nonzero"`
	Date        time.Time `gorm:"date" json:"string" validate:"nonzero"`
}
