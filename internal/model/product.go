package model

import "time"

type Product struct {
	ID          int       `json:"id"`
	SellerID    int       `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductListQuery struct {
	Page     int
	Limit    int
	Search   string
	MinPrice *int
	MaxPrice *int
}
