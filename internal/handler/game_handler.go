package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"tictactoe/internal/model"
	"tictactoe/internal/service"
	"tictactoe/internal/storage"

	"github.com/google/uuid"
)

type GameHandler struct {
	Service *service.GameService
	Repo    *storage.GameRepository
}

type MoveRequest struct {
	Board [3][3]int `json:"board"`
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	game := model.NewGame()
	var err error

	err = h.Repo.Save(game)
	if err != nil {
		log.Printf("Ошибка при сохранении игры: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		log.Printf("Ошибка при отправке JSON: %v", err)
		return
	}
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	game, err := h.Repo.Load(id)
	if err != nil {
		http.Error(w, "Игра не найдена", http.StatusNotFound)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	playerMove, err := h.Service.PlayerMove(game, req.Board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.Repo.Save(playerMove)
	if err != nil {
		log.Printf("Ошибка при сохранении игры: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		log.Printf("Ошибка при отправке JSON: %v", err)
		return
	}
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	game, err := h.Repo.Load(id)
	if err != nil {
		http.Error(w, "Игра не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		log.Printf("Ошибка при отправке JSON: %v", err)
		return
	}
}
