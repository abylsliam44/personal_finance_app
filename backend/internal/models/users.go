package models

import "time"

type User struct {
    ID               int       `json:"id"`
    Name             string    `json:"name"`
    Email            string    `json:"email"`
    PasswordHash     string    `json:"password_hash"`
    PreferredCurrency string   `json:"preferred_currency"`
    CreatedAt        time.Time `json:"created_at"`
}
