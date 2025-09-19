package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nkchakradhari780/catalogServices/internal/repository/storage"
	"github.com/nkchakradhari780/catalogServices/internal/utils/response"
)

func AddToCart(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.PathValue("user_id")
		productIDStr := r.PathValue("product_id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid user_id")))
			return
		}

		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid product_id")))
			return
		}

		var body struct {
			Quantity int `json:"quantity"`
			Discount int `json:"discount"` 
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Quantity <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid input")))
			return
		}

		cartItemID, err := storage.AddToCart(userID, productID, body.Quantity, body.Discount)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"message":      "Item added to cart successfully",
			"cart_item_id": cartItemID,
		})
	}
}

func RemoveFromCart(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.PathValue("user_id")
		productIDStr := r.PathValue("product_id")

		userId, err := strconv.Atoi(userIDStr)
		
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid user id")))
			return 
		}
		
		productId, err := strconv.Atoi(productIDStr)
		
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid product id")))
			return 
		}

		if err = storage.RemoveFromCart(userId, productId); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "item removed from successfully!", "result": "success"})
	}
}