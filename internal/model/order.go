package model

import "time"

type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	TotalPrice int         `json:"total_price"`
	Status     string      `json:"status"`
	Address    string      `json:"address"`
	CreatedAt  time.Time   `json:"created_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}

type CheckoutRequest struct {
	Address string `json:"address" binding:"required,min=10"`
}
