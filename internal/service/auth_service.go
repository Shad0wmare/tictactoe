package service

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/auth"

	"tictactoe/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, credentials model.SignUpRequest) (*model.User, error)
	Authorize(ctx context.Context, req model.JwtRequest) (*model.JwtResponse, error)
	RefreshAccessToken(ctx context.Context, refreshTokenStr string) (*model.JwtResponse, error)
	RefreshTokens(ctx context.Context, refreshTokenStr string) (*model.JwtResponse, error)
}
type authService struct {
	userService UserService
	jwtProvider auth.JwtProvider
}

func NewAuthService(userService UserService, jwtProvider auth.JwtProvider) AuthService {
	return &authService{userService: userService, jwtProvider: jwtProvider}
}

func (s *authService) Register(ctx context.Context, credentials model.SignUpRequest) (*model.User, error) {
	createdUser, err := s.userService.CreateUser(ctx, credentials.Login, credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("не удалось зарегистрировать пользователя: %w", err)
	}
	return createdUser, err
}

func (s *authService) Authorize(ctx context.Context, req model.JwtRequest) (*model.JwtResponse, error) {
	user, err := s.userService.GetByLogin(ctx, req.Login)
	if err != nil || user.Password != req.Password {
		return nil, errors.New("неверный логин или пароль")
	}

	userInfo := &model.UserInfo{
		UserId: user.UserId,
		Login:  user.Login,
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userInfo)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(userInfo)
	if err != nil {
		return nil, err
	}

	response := model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, err
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (*model.JwtResponse, error) {
	claims, err := s.jwtProvider.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("невалидный refresh token: %w", err)
	}

	user, err := s.userService.GetByUserId(ctx, claims.UserId)
	if err != nil {
		return nil, errors.New("пользователь из токена не найден")
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshTokenStr string) (*model.JwtResponse, error) {
	claims, err := s.jwtProvider.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("невалидный refresh токен: %w", err)
	}

	user, err := s.userService.GetByUserId(ctx, claims.UserId)
	if err != nil {
		return nil, errors.New("пользователь из токена не найден")
	}

	newAccessToken, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
