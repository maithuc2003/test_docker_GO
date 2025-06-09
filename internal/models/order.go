package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	UserID    int       `json:"user_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	OrderedAt time.Time `json:"ordered_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
