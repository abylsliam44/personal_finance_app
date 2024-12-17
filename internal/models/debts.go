package models

import "time"

type Debt struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Contact   string    `json:"contact"`
    Amount    float64   `json:"amount"`
    DueDate   time.Time `json:"due_date"`
    CreatedAt time.Time `json:"created_at"`
}
