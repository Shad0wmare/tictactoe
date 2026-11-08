package di

import (
	"tictactoe/cmd/server"
	"tictactoe/internal/auth"
	"tictactoe/internal/database"
	"tictactoe/internal/handler"
	"tictactoe/internal/repository"
	"tictactoe/internal/service"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	app := fx.New(
		fx.Provide(
			database.NewDBPool,
			repository.NewGameRepository,
			repository.NewUserRepository,
			auth.NewJwtProvider,
			service.NewGameService,
			service.NewUserService,
			service.NewAuthService,
			handler.NewGameHandler,
			handler.NewAuthHandler,
			handler.NewUserAuthenticator,
			handler.NewUserHandler,
		),
		fx.Invoke(server.StartHTTPServer),
	)
	return app
}
