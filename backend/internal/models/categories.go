package models

import "time"

type Category struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "income" or "expense"
	CreatedAt time.Time `json:"created_at"`
}
