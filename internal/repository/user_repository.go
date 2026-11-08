package repository

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserInfo, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (user_id, login, password) VALUES ($1, $2, $3);`
	_, err := r.pool.Exec(ctx, query, user.UserId, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("ошибка записи пользователя в БД: %w", err)
	}
	return err
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `SELECT user_id, login, password FROM users WHERE login = $1;`

	var user model.User
	err := r.pool.QueryRow(ctx, query, login).Scan(&user.UserId, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserInfo, error) {
	query := `SELECT user_id, login FROM users WHERE user_id = $1;`

	var user model.UserInfo
	err := r.pool.QueryRow(ctx, query, userId).Scan(&user.UserId, &user.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с userId '%s' не найден", userId)
		}
		return nil, fmt.Errorf("ошибка чтения пользователя из БД: %w", err)
	}

	return &user, err
}
