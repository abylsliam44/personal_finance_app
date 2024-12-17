package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	AccountID   int       `json:"account_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` //"income" or "expense"
	CategoryID  int       `json:"category"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
