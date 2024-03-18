package model

import "time"

type Receipt struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"filename"`
	Filename    string    `json:"date"`
	Date        time.Time `json:"string"`
}
