package models

import "time"

type Author struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
