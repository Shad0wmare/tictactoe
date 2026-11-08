package model

import "github.com/google/uuid"

type PlayerStats struct {
	UserId   uuid.UUID `json:"user_id"`
	Login    string    `json:"login"`
	WinRatio float64   `json:"win_ratio"`
}
