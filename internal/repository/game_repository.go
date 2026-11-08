package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tictactoe/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GameReadyToJoin = "READY_TO_JOIN"
	Player1Won      = "PLAYER_1_WON"
	Player2Won      = "PLAYER_2_WON"
	Draw            = "DRAW"
)

type GameRepository interface {
	Save(ctx context.Context, game *model.Game) error
	Load(ctx context.Context, id uuid.UUID) (*model.Game, error)
	GetGamesToJoin(ctx context.Context) ([]uuid.UUID, error)
	GetPlayerGames(ctx context.Context, userId uuid.UUID) ([]model.Game, error)
	GetLeaders(ctx context.Context, depth int64) ([]model.PlayerStats, error)
}

type gameRepository struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) GameRepository {
	return &gameRepository{pool: pool}
}

func (r *gameRepository) Save(ctx context.Context, game *model.Game) error {
	query := `
		INSERT INTO games (game_id, board, player_1, player_2, is_computer, status, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (game_id) 
		DO UPDATE SET 
					player_2 = EXCLUDED.player_2,
					status = EXCLUDED.status,
					board = EXCLUDED.board;
	`

	_, err := r.pool.Exec(ctx, query, game.ID, game.Board, game.Player1, game.Player2, game.IsComputer, game.Status, game.Created)
	if err != nil {
		return fmt.Errorf("ошибка сохранения игры в БД: %w", err)
	}

	return err
}

func (r *gameRepository) Load(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	query := `SELECT game_id, board, player_1, player_2, is_computer, status, created
				FROM games
				WHERE game_id = $1;`

	var game model.Game
	var board []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(&game.ID, &board, &game.Player1, &game.Player2, &game.IsComputer, &game.Status, &game.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("игра с id '%s' не найдена", id)
		}
		return nil, fmt.Errorf("ошибка чтения игры из БД: %w", err)
	}

	if err = json.Unmarshal(board, &game.Board); err != nil {
		return nil, fmt.Errorf("ошибка получения поля board: %w", err)
	}

	return &game, nil
}

func (r *gameRepository) GetGamesToJoin(ctx context.Context) ([]uuid.UUID, error) {
	query := `SELECT game_id
				FROM games
				WHERE status = $1;`

	var games []uuid.UUID
	var game uuid.UUID
	rows, err := r.pool.Query(ctx, query, GameReadyToJoin)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка доступных к присоединению игр из БД: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&game)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, err
}

func (r *gameRepository) GetPlayerGames(ctx context.Context, userId uuid.UUID) ([]model.Game, error) {
	query := `SELECT game_id, board, player_1, player_2, is_computer, status, created
				FROM games
				WHERE (player_1 = $1 or player_2 = $1) and
				      (status = $2 or status = $3 or status = $4) 
				order by created desc;`

	var games []model.Game
	var game model.Game
	var board []byte
	rows, err := r.pool.Query(ctx, query, userId, Player1Won, Player2Won, Draw)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка игр пользователя %s: %w", userId, err)
	}
	for rows.Next() {
		err = rows.Scan(&game.ID, &board, &game.Player1, &game.Player2, &game.IsComputer, &game.Status, &game.Created)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(board, &game.Board); err != nil {
			return nil, fmt.Errorf("ошибка получения поля board: %w", err)
		}
		games = append(games, game)
	}
	return games, err
}

func (r *gameRepository) GetLeaders(ctx context.Context, depth int64) ([]model.PlayerStats, error) {
	query := `WITH player_stats AS (
					SELECT 
						u.user_id AS user_id,
						u.login as login,
						COUNT(CASE 
							WHEN g.player_1 = u.user_id AND g.status = 'PLAYER_1_WON' THEN 1
							WHEN g.player_2 = u.user_id AND g.status = 'PLAYER_2_WON' THEN 1
						END) AS wins,
						COUNT(CASE 
							WHEN g.player_1 = u.user_id AND g.status = 'PLAYER_2_WON' THEN 1
							WHEN g.player_2 = u.user_id AND g.status = 'PLAYER_1_WON' THEN 1
						END) AS losses,
						COUNT(CASE 
							WHEN g.status = 'DRAW' THEN 1 
						END) AS draws
					FROM users u
					JOIN games g ON (g.player_1 = u.user_id OR g.player_2 = u.user_id)
					WHERE g.status IN ('PLAYER_1_WON', 'PLAYER_2_WON', 'DRAW')
					GROUP BY u.user_id, u.login
			)
			SELECT 
				user_id,
				login,
				ROUND(
					CASE 
						WHEN (losses + draws) = 0 AND wins > 0 THEN wins::numeric
						WHEN (losses + draws) = 0 AND wins = 0 THEN 0.0
						ELSE (wins::numeric / (losses + draws)::numeric)
					END, 
					2
				)::double precision AS win_rate
			FROM player_stats
			ORDER BY win_rate DESC
			LIMIT $1;`

	var playerStats model.PlayerStats
	var leaderboard []model.PlayerStats
	rows, err := r.pool.Query(ctx, query, depth)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка доступных к присоединению игр из БД: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&playerStats.UserId, &playerStats.Login, &playerStats.WinRatio)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, playerStats)
	}

	return leaderboard, err
}
