package service

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/model"
	"tictactoe/internal/repository"
	"time"

	"github.com/google/uuid"
)

const (
	readyToJoin        = "READY_TO_JOIN"
	player1Turn        = "PLAYER_1_TURN"
	player2Turn        = "PLAYER_2_TURN"
	opponentIsComputer = "computer"
	Player1Won         = "PLAYER_1_WON"
	Player2Won         = "PLAYER_2_WON"
	Draw               = "DRAW"
	noWinner           = ""
)

type GameService interface {
	IsGameOver(gameStatus string) bool
	PlayerMove(savedGameState *model.Game, userId uuid.UUID, playerMoveState [3][3]int) (*model.Game, error)
	ValidatePlayerMove(move string, savedGameState, playerMoveState [3][3]int) error
	ComputerMove(g *model.Game)
	Minimax(board [3][3]int, depth int, isMax bool) int
	Evaluate(b [3][3]int) int
	winCheck(b [3][3]int) string
	HasEmptyCells(b [3][3]int) bool
	CreateGame(gameCreator uuid.UUID, opponent string) (*model.Game, error)
	SaveGame(ctx context.Context, game *model.Game) error
	LoadGame(ctx context.Context, id uuid.UUID) (*model.Game, error)
	GetGamesToJoin(ctx context.Context) ([]uuid.UUID, error)
	JoinGame(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) error
	GetPlayerGames(ctx context.Context, userId uuid.UUID) ([]model.Game, error)
	GetLeaders(ctx context.Context, depth int64) ([]model.PlayerStats, error)
}
type gameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) GameService {
	return &gameService{repo: repo}
}

func (s *gameService) IsGameOver(gameStatus string) bool {
	if gameStatus == Player1Won || gameStatus == Player2Won || gameStatus == Draw {
		return true
	}
	return false
}

func (s *gameService) PlayerMove(savedGameState *model.Game, userId uuid.UUID, playerMoveState [3][3]int) (*model.Game, error) {
	if s.IsGameOver(savedGameState.Status) {
		return nil, fmt.Errorf("невозможно сделать ход, статус игры: %s", savedGameState.Status)
	}

	if savedGameState.Status == player2Turn && savedGameState.Player1 == userId {
		return nil, errors.New("сейчас не ваш ход")
	}
	if savedGameState.Status == player1Turn && !savedGameState.IsComputer && *savedGameState.Player2 == userId {
		return nil, errors.New("сейчас не ваш ход")
	}

	err := s.ValidatePlayerMove(savedGameState.Status, savedGameState.Board, playerMoveState)
	if err != nil {
		return nil, err
	}

	savedGameState.Board = playerMoveState
	status := s.winCheck(savedGameState.Board)
	if status != noWinner {
		savedGameState.Status = status
	}

	if savedGameState.IsComputer {
		s.ComputerMove(savedGameState)
		status = s.winCheck(savedGameState.Board)
		if status != noWinner {
			savedGameState.Status = status
		}
		return savedGameState, nil
	}

	if savedGameState.Status == player1Turn {
		savedGameState.Status = player2Turn
	} else if savedGameState.Status == player2Turn {
		savedGameState.Status = player1Turn
	}

	return savedGameState, nil
}

func (s *gameService) ValidatePlayerMove(move string, savedGameState, playerMoveState [3][3]int) error {
	boardChangeCount := 0
	for i := range 3 {
		for j := range 3 {
			if savedGameState[i][j] != playerMoveState[i][j] {
				if savedGameState[i][j] != model.Empty {
					return errors.New("нельзя изменять сделанные ходы")
				}
				if move == player1Turn && playerMoveState[i][j] != model.X {
					return errors.New("неправильный ход, можно ставить только крестики")
				}
				if move == player2Turn && playerMoveState[i][j] != model.O {
					return errors.New("неправильный ход, можно ставить только нолики")
				}
				boardChangeCount++
			}
		}
	}
	if boardChangeCount > 1 {
		return errors.New("разрешено сделать только один ход")
	}
	if boardChangeCount == 0 {
		return errors.New("необходимо сделать ход")
	}
	return nil
}

func (s *gameService) ComputerMove(g *model.Game) {
	bestScore := -1000
	var bestMove [2]int

	for i := range 3 {
		for j := range 3 {
			if g.Board[i][j] == model.Empty {
				g.Board[i][j] = model.O
				score := s.Minimax(g.Board, 0, false)
				g.Board[i][j] = model.Empty
				if score > bestScore {
					bestScore = score
					bestMove = [2]int{i, j}
				}
			}
		}
	}
	g.Board[bestMove[0]][bestMove[1]] = model.O
}

func (s *gameService) Minimax(board [3][3]int, depth int, isMax bool) int {
	score := s.Evaluate(board)
	if score == 10 {
		return score - depth
	}
	if score == -10 {
		return score + depth
	}
	if !s.HasEmptyCells(board) {
		return 0
	}

	if isMax {
		best := -1000
		for i := range 3 {
			for j := range 3 {
				if board[i][j] == model.Empty {
					board[i][j] = model.O
					score = s.Minimax(board, depth+1, false)
					if score > best {
						best = score
					}
					board[i][j] = model.Empty
				}
			}
		}
		return best
	}

	best := 1000
	for i := range 3 {
		for j := range 3 {
			if board[i][j] == model.Empty {
				board[i][j] = model.X
				score = s.Minimax(board, depth+1, true)
				if score < best {
					best = score
				}
				board[i][j] = model.Empty
			}
		}
	}
	return best
}

func (s *gameService) Evaluate(b [3][3]int) int {
	for i := range 3 {
		if b[i][0] != model.Empty && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			if b[i][0] == model.O {
				return 10
			}
			return -10
		}
		if b[0][i] != model.Empty && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			if b[0][i] == model.O {
				return 10
			}
			return -10
		}
	}
	if b[1][1] != model.Empty {
		if (b[0][0] == b[1][1] && b[1][1] == b[2][2]) || (b[0][2] == b[1][1] && b[1][1] == b[2][0]) {
			if b[1][1] == model.O {
				return 10
			}
			return -10
		}
	}
	return 0
}

func (s *gameService) winCheck(b [3][3]int) string {
	result := 0
	status := noWinner
	for i := range 3 {
		if b[i][0] != 0 && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			result = b[i][0]
			break
		}
		if b[0][i] != 0 && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			result = b[0][i]
			break
		}
	}
	if b[0][0] != 0 && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		result = b[0][0]
	}
	if b[0][2] != 0 && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		result = b[0][2]
	}

	if result == 1 {
		status = Player1Won
	} else if result == 2 {
		status = Player2Won
	} else if !s.HasEmptyCells(b) {
		status = Draw
	}

	return status
}

func (s *gameService) HasEmptyCells(b [3][3]int) bool {
	for i := range 3 {
		for j := range 3 {
			if b[i][j] == model.Empty {
				return true
			}
		}
	}
	return false
}

func (s *gameService) CreateGame(gameCreator uuid.UUID, opponent string) (*model.Game, error) {
	game := &model.Game{
		ID:      uuid.New(),
		Player1: gameCreator,
		Board:   [3][3]int{},
		Created: time.Now(),
	}

	if opponent == opponentIsComputer {
		game.IsComputer = true
		game.Player2 = nil
		game.Status = player1Turn
	} else {
		game.IsComputer = false
		game.Player2 = nil
		game.Status = readyToJoin
	}

	return game, nil
}

func (s *gameService) SaveGame(ctx context.Context, game *model.Game) error {
	err := s.repo.Save(ctx, game)
	return err
}

func (s *gameService) LoadGame(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	game, err := s.repo.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return game, err
}

func (s *gameService) GetGamesToJoin(ctx context.Context) ([]uuid.UUID, error) {
	gamesList, err := s.repo.GetGamesToJoin(ctx)
	if err != nil {
		return nil, err
	}
	if gamesList == nil {
		err = errors.New("нет доступных игр для присоединения")
	}
	return gamesList, err
}

func (s *gameService) JoinGame(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) error {
	game, err := s.repo.Load(ctx, gameId)
	if err != nil {
		return err
	}
	if game.Status != readyToJoin {
		return errors.New("нельзя присоединиться к этой игре, она уже идёт/завершена")
	}
	if game.Player1 == userId {
		return errors.New("нельзя играть против самого себя")
	}
	game.Player2 = &userId
	game.Status = player1Turn
	err = s.repo.Save(ctx, game)
	return err
}

func (s *gameService) GetPlayerGames(ctx context.Context, userId uuid.UUID) ([]model.Game, error) {
	finishedGames, err := s.repo.GetPlayerGames(ctx, userId)
	if err != nil {
		return nil, err
	}
	if finishedGames == nil {
		err = errors.New("у пользователя нет завершенных игр")
	}
	return finishedGames, err
}

func (s *gameService) GetLeaders(ctx context.Context, depth int64) ([]model.PlayerStats, error) {
	leaderboard, err := s.repo.GetLeaders(ctx, depth)
	if err != nil {
		return nil, err
	}
	return leaderboard, err
}
