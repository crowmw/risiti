package model

import "time"

type Receipt struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"filename"`
	Filename    string    `json:"date"`
	Date        time.Time `json:"string"`
	CreatedBy   uint64    `json:"created_by"`
}
