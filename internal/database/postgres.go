package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

const DbConnString = "postgres://postgres:1232@localhost:5432/tictactoe?sslmode=disable"

func NewDBPool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	connString := DbConnString
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
