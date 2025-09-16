package modules

type Cart struct {
	ID         int `json:"cart_id,omitempty"`
	CartItemID int `json:"cartItemId" validate:"required"`
}

type CartItem struct {
	ID         int      `json:"cart_item_id,omitempty" `
	Name       string   `json:"name" validate:"required"`
	Price      int      `json:"price" validate:"required"`
	Stock      int      `json:"stock" validate:"required"`
	CategoryID string   `json:"category_id" validate:"required"`
	Brand      string   `json:"brand" validate:"required"`
	Images     []string `json:"images,omitempty"`
}
