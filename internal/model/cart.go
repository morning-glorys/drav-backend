package model

type Cart struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

type CartItem struct {
	ID           int    `json:"id"`
	CartID       int    `json:"cart_id"`
	ProductID    int    `json:"product_id"`
	Quantity     int    `json:"quantity"`
	ProductName  string `json:"product_name"`
	ProductPrice int    `json:"product_price"`
}

type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}
