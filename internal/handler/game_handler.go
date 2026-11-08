package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tictactoe/internal/model"
	"tictactoe/internal/service"

	"github.com/google/uuid"
)

type GameHandler struct {
	service service.GameService
}

func NewGameHandler(service service.GameService) *GameHandler {
	return &GameHandler{service: service}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatorId, _ := ctx.Value("userId").(uuid.UUID)
	var req model.NewGameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректная модель новой игры", http.StatusBadRequest)
		return
	}

	if req.Opponent != "computer" && req.Opponent != "player" {
		http.Error(w, "body.opponent должно принимать значения computer или player", http.StatusBadRequest)
		return
	}

	game, err := h.service.CreateGame(creatorId, req.Opponent)
	if err != nil {
		http.Error(w, "Ошибка при создании игры", http.StatusInternalServerError)
		return
	}

	err = h.service.SaveGame(ctx, game)
	if err != nil {
		http.Error(w, "Ошибка при сохранении игры", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		http.Error(w, "Ошибка отправки игрового поля", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := ctx.Value("userId").(uuid.UUID)
	gameIdStr := r.PathValue("id")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	game, err := h.service.LoadGame(ctx, gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var req model.Move
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректная модель поля", http.StatusBadRequest)
		return
	}

	playerMove, err := h.service.PlayerMove(game, userId, req.Board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.service.SaveGame(ctx, playerMove)
	if err != nil {
		http.Error(w, "Ошибка при сохранении игры", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ход выполнен")
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID", http.StatusBadRequest)
		return
	}

	game, err := h.service.LoadGame(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		http.Error(w, "Ошибка отправки игрового поля", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) GetGamesToJoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	game, err := h.service.GetGamesToJoin(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		http.Error(w, "Ошибка отправки списка игр", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	player, _ := ctx.Value("userId").(uuid.UUID)
	gameIdStr := r.PathValue("id")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		http.Error(w, "Неверный формат UUID игры", http.StatusBadRequest)
		return
	}
	err = h.service.JoinGame(ctx, gameId, player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Присоединение успешно")
}

func (h *GameHandler) GetPlayerGames(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	player, _ := ctx.Value("userId").(uuid.UUID)
	games, err := h.service.GetPlayerGames(ctx, player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(games)
	if err != nil {
		http.Error(w, "Ошибка отправки списка игр", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) GetLeaders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	depthStr := r.PathValue("depth")
	depth, err := strconv.ParseInt(depthStr, 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	leaderboard, err := h.service.GetLeaders(ctx, depth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(leaderboard)
	if err != nil {
		http.Error(w, "Ошибка отправки списка игр", http.StatusInternalServerError)
		return
	}
}
