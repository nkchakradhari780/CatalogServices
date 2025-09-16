package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nkchakradhari780/catalogServices/internal/repository/storage"
	"github.com/nkchakradhari780/catalogServices/internal/utils/response"
)

func AddToWishList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id_str := r.PathValue("user_id")

		user_id, err := strconv.Atoi(user_id_str)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		product_id_str := r.PathValue("product_id")

		product_id, err := strconv.Atoi(product_id_str)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		wishListId, err := storage.AddToWishList(user_id, product_id)

		if err != nil {
			if err.Error() == "product already added to wish list"{
				response.WriteJson(w, http.StatusInternalServerError, map[string]string{"message": "product already added to wish list"})
				return 
			}
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Item with", slog.String("Id", fmt.Sprint(wishListId)))
		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Item Added to wishlist Successfully"})

	}
}
