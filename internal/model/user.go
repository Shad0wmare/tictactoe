package model

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID
	Login    string `json:"login" database:"login"`
	Password string `json:"password" database:"password"`
}

type UserInfo struct {
	UserId uuid.UUID `json:"id"`
	Login  string    `json:"login"`
}
