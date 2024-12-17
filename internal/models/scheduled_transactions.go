package models

import "time"

type ScheduledTransaction struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    AccountID   int       `json:"account_id"`
    Amount      float64   `json:"amount"`
    Type        string    `json:"type"` //"income" or "expense"
    Schedule    string    `json:"schedule"` //"daily", "weekly", "monthly"
    CreatedAt   time.Time `json:"created_at"`
}
