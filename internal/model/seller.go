package model

import "time"

type Seller struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	StoreName  string    `json:"store_name"`
	IsVerified bool      `json:"is_verified"`
	Rating     float64   `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
}
