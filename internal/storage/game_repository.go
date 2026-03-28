package storage

import (
	"errors"
	"sync"
	"tictactoe/internal/model"

	"github.com/google/uuid"
)

type GameRepository struct {
	games sync.Map
}

func NewGameRepository() *GameRepository {
	return &GameRepository{}
}

func (r *GameRepository) Save(game *model.Game) error {
	r.games.Store(game.ID, game)
	return nil
}

func (r *GameRepository) Load(id uuid.UUID) (*model.Game, error) {
	value, ok := r.games.Load(id)
	if !ok {
		return nil, errors.New("game not found")
	}
	return value.(*model.Game), nil
}
