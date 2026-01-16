package httpAuth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/dtoAuth"
	"github.com/seva-up/App_Products/internal/models"
)

type authDelivery struct {
	authUS auth.UserService
}

func NewAuthDelivery(authUS auth.UserService) auth.UserDelivery {
	return &authDelivery{authUS: authUS}
}

func (h *authDelivery) Register(w http.ResponseWriter, r *http.Request) {
	var in dtoAuth.InRegisters

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if in.Email == "" || in.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user := &models.User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
		Email:     in.Email,
		Role:      in.Role,
	}
	ctx := r.Context()
	createdUser, err := h.authUS.Register(ctx, user)
	if err != nil {
		// Проверяем тип ошибки для корректного HTTP статуса
		if strings.Contains(err.Error(), "уже существует") {
			http.Error(w, "пользователь уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "ошибка при регистрации", http.StatusInternalServerError)
		return
	}
	createdUser.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
