package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/model"
	"tictactoe/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректная модель пользователя", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := model.SignUpResponse{
		UserId: &user.UserId,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Ошибка отправки JWT ответа", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректная модель JWT запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Authorize(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Ошибка отправки JWT ответа", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректная модель запроса JWT Refresh", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Ошибка отправки JWT ответа", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректная модель запроса JWT Refresh", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Ошибка отправки JWT ответа", http.StatusInternalServerError)
		return
	}
}
