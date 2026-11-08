package handler

import (
	"context"
	"net/http"
	"strings"
	"tictactoe/internal/auth"

	"github.com/google/uuid"
)

type UserAuthenticator struct {
	jwtProvider auth.JwtProvider
}

func NewUserAuthenticator(jwtProvider auth.JwtProvider) *UserAuthenticator {
	return &UserAuthenticator{jwtProvider: jwtProvider}
}

func (a *UserAuthenticator) CheckRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/register" ||
			path == "/login" ||
			path == "/auth/refreshAccessToken" ||
			path == "/auth/refreshTokens" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "некорректный формат токена", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		userClaims, err := a.jwtProvider.ValidateAccessToken(tokenStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		signedReq := a.Sign(r, userClaims.UserId)
		next.ServeHTTP(w, signedReq)
	})
}

func (a *UserAuthenticator) Sign(r *http.Request, userId uuid.UUID) *http.Request {
	ctx := context.WithValue(r.Context(), "userId", userId)
	return r.WithContext(ctx)
}
