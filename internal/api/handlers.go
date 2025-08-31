package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nkchakradhari780/catalogServices/internal/modules"
	"github.com/nkchakradhari780/catalogServices/internal/repository/storage"
	"github.com/nkchakradhari780/catalogServices/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var product modules.Product

		err := json.NewDecoder(r.Body).Decode(&product)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(product); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateProduct(product.Name, product.Price, product.Stock, product.CategoryID, product.Brand, product.Images)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Creating New Product", slog.String("productId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})

	}
}
