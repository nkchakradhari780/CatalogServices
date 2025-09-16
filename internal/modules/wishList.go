package modules

type WishList struct {
	WishListId int    `json:"wishList_id,omitempty" `
	ProductId  int    `json:"product_id,omitempty" validate:"required"`
	UserId     int    `json:"user_id,omitempty" validate:"required"`
	AddedAt    string `json:"added_at"`
}
