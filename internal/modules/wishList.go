package modules

type WishList struct {
	ID         int      `json:"wishList_id,omitempty" `
	ProductId  int      `json:"product_id,omitempty" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Price      int      `json:"price" validate:"required"`
	Stock      int      `json:"stock" validate:"required"`
	CategoryID string   `json:"category_id" validate:"required"`
	Brand      string   `json:"brand" validate:"required"`
	Images     []string `json:"images,omitempty"`
}
