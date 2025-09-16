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
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateNewUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user modules.Users

		err := json.NewDecoder(r.Body).Decode(&user)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(user); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		hashPass, err := HashPassword(user.Password)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("error hashing password")))
			return
		}

		userId, err := storage.CreateUser(user.Name, user.Email, hashPass, user.Phone, user.Role, user.Address)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Creating New User", slog.String("UserId", fmt.Sprint(userId)))
		response.WriteJson(w, http.StatusOK, map[string]string{"message": "User Created Successfully"})
	}
}
