package models

import "time"

type Deposit struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    AccountID   int       `json:"account_id"`
    InitialAmount float64 `json:"initial_amount"`
    InterestRate float64  `json:"interest_rate"`
    CreatedAt   time.Time `json:"created_at"`
    EndDate     time.Time `json:"end_date"`
}
