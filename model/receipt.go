package model

import "time"

type Receipt struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name" validate:"nonzero"`
	Description string    `json:"filename"`
	Filename    string    `json:"date" validate:"nonzero"`
	Date        time.Time `json:"string" validate:"nonzero"`
	CreatedBy   uint64    `json:"created_by"`
}
