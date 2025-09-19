package modules

import "time"

type Cart struct {
	CartId     int `json:"cart_id,omitempty"`
	CartItemID int `json:"cartItemId" validate:"required"`
}

type CartItem struct {
	CartItemId  int       `json:"cart_item_id,omitempty" `
	CartId      int       `json:"cart_id" validate:"required"`
	ProductId   int       `json:"product_id" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required"`
	PriceAtTime float64   `json:"price_at_time" validate:"required"`
	Discount    float64   `json:"discount" validate:"required"`
	Subtotal    float64   `json:"subtotal"`
	AddedAt     time.Time `json:"added_at"`
}
