package model

import "github.com/google/uuid"

const (
	Empty = 0
	X     = 1
	O     = 2
)

type Game struct {
	ID    uuid.UUID `json:"id"`
	Board [3][3]int `json:"board"`
}

func NewGame() *Game {
	return &Game{
		ID:    uuid.New(),
		Board: [3][3]int{},
	}
}
