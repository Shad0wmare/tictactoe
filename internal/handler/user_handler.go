package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/service"

	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) GetByUserId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	userInfo, err := u.service.GetByUserId(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		http.Error(w, "Ошибка отправки информации о пользователе", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
