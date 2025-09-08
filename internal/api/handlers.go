package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/nkchakradhari780/catalogServices/internal/modules"
	"github.com/nkchakradhari780/catalogServices/internal/repository/storage"
	"github.com/nkchakradhari780/catalogServices/internal/utils/response"
)

func CreateNewProduct(storage storage.Storage) http.HandlerFunc {
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

func GetProductById(storage storage.Storage) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		slog.Info("Fetching Product", slog.String("productId", idStr))


		id, err := strconv.Atoi(idStr)
        if err != nil {
            slog.Error("Invalid product id", slog.String("productId", idStr), slog.String("error", err.Error()))
            response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid product id")))
            return
        }
		product, err := storage.GetProductById(id) 
		if err != nil {
			slog.Error("Error fetching product", slog.String("productId", idStr), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, product)
	}
}

func GetProducts(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Fetching all products")

		products, err := storage.GetProducts()
		if err != nil {
			fmt.Println("Error fetching products:", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, products)
	}
}

func GetDefaultProducts(storate storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Fetching Default Products")

		products, err := storate.GetDefaultProducts()
		if err != nil {
			fmt. Println("Error Fetching products: ", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError((err)))
			return 
		}

		response.WriteJson(w, http.StatusOK, products)
	}
}

func GetFilteredProducts(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Fetching Filtered products")
		filters := r.URL.Query()

		products, err := storage.GetFilteredProducts(filters)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, products)
	}
}

func UpdateProductById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Updating Product")

		idStr:= r.PathValue("id")
		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		var product modules.Product
		err = json.NewDecoder(r.Body).Decode(&product)
		
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		if err := validator.New().Struct(product); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 
		}

		updatedProduct, err := storage.UpdateProductById(id, product.Name, product.Price, product.Stock, product.CategoryID, product.Brand, product.Images)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError((err)))
			return 
		}

		response.WriteJson(w, http.StatusOK, updatedProduct)

	}
}

func DeleteProductById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting Product")

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr) 
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}
		if err = storage.DeleteProductById(id); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"result": "success"})
		
	}
}