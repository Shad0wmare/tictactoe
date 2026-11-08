package model

import "github.com/google/uuid"

type SignUpRequest struct {
	Login    string `json:"login" database:"login"`
	Password string `json:"password" database:"password"`
}

type SignUpResponse struct {
	UserId *uuid.UUID
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}
