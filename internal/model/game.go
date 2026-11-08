package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	Empty = 0
	X     = 1
	O     = 2
)

type Game struct {
	ID         uuid.UUID  `json:"id" database:"game_id"`
	Board      [3][3]int  `json:"board" database:"board"`
	Player1    uuid.UUID  `json:"player_1" database:"player_1"`
	Player2    *uuid.UUID `json:"player_2,omitempty" database:"player_2"`
	IsComputer bool       `json:"is_computer" database:"is_computer"`
	Status     string     `json:"status" database:"status"`
	Created    time.Time  `json:"created" database:"created"`
}

type Move struct {
	Board [3][3]int `json:"board"`
}

type NewGameRequest struct {
	Opponent string `json:"opponent"`
}
