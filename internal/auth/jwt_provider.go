package auth

import (
	"errors"
	"fmt"
	"time"

	"tictactoe/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

type JwtProvider struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJwtProvider() JwtProvider {
	return JwtProvider{
		accessSecret:  []byte("access-key"),
		refreshSecret: []byte("refresh-key"),
		accessTTL:     10 * time.Minute,
		refreshTTL:    3 * time.Hour,
	}
}

func (j *JwtProvider) GenerateAccessToken(userId *model.UserInfo) (string, error) {
	claims := UserClaims{
		UserId: userId.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId.Login,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.accessSecret)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации access token: %w", err)
	}
	return signedToken, nil
}

func (j *JwtProvider) GenerateRefreshToken(user *model.UserInfo) (string, error) {
	claims := UserClaims{
		UserId: user.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Login,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.refreshSecret)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации refresh token: %w", err)
	}
	return signedToken, nil
}

func (j *JwtProvider) ValidateAccessToken(tokenStr string) (*UserClaims, error) {
	return j.validateToken(tokenStr, j.accessSecret)
}

func (j *JwtProvider) ValidateRefreshToken(tokenStr string) (*UserClaims, error) {
	return j.validateToken(tokenStr, j.refreshSecret)
}

func (j *JwtProvider) GetUUIDFromToken(tokenStr string) (uuid.UUID, error) {
	claims, err := j.ValidateAccessToken(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserId, nil
}

func (j *JwtProvider) validateToken(tokenStr string, secret []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи токена: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("токен недействителен: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("некорректная структура claims в токене")
	}

	return claims, nil
}
