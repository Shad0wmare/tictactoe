package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"tictactoe/internal/handler"

	"go.uber.org/fx"
)

const (
	Url  = "http://localhost"
	Port = ":8080"
)

func StartHTTPServer(lc fx.Lifecycle, g *handler.GameHandler, a *handler.AuthHandler, r *handler.UserAuthenticator, u *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", a.RegisterHandler)
	mux.HandleFunc("POST /login", a.LoginHandler)
	mux.HandleFunc("POST /auth/refreshAccessToken", a.RefreshAccessTokenHandler)
	mux.HandleFunc("POST /auth/refreshTokens", a.RefreshTokensHandler)

	mux.HandleFunc("GET /user/{id}", u.GetByUserId)
	mux.HandleFunc("POST /game", g.CreateGame)
	mux.HandleFunc("GET /game/{id}", g.GetGame)
	mux.HandleFunc("GET /games", g.GetGamesToJoin)
	mux.HandleFunc("GET /leaderboard/{depth}", g.GetLeaders)
	mux.HandleFunc("GET /games/my", g.GetPlayerGames)
	mux.HandleFunc("POST /game/{id}/join", g.JoinGame)
	mux.HandleFunc("POST /game/{id}/move", g.MakeMove)

	protected := r.CheckRoute(mux)

	srv := &http.Server{
		Addr:    Port,
		Handler: protected,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Сервер запускается на %s%s\n", Url, Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("Ошибка сервера: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Остановка HTTP-сервера...")
			return srv.Shutdown(ctx)
		},
	})

	return mux
}
