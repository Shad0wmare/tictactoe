package service

import (
	"context"
	"errors"
	"tictactoe/internal/repository"

	"tictactoe/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService interface {
	CreateUser(ctx context.Context, login, password string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserInfo, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, login, password string) (*model.User, error) {
	if login == "" || password == "" {
		return nil, errors.New("логин и пароль не могут быть пустыми")
	}

	existing, err := s.repo.GetByLogin(ctx, login)
	if existing != nil {
		return nil, errors.New("пользователь с таким логином уже существует")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	user := model.User{
		UserId:   uuid.New(),
		Login:    login,
		Password: password,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *userService) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.UserInfo, error) {
	user, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, err
}
