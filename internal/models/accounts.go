package models

import "time"

type Account struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
