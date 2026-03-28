package service

import (
	"errors"
	"tictactoe/internal/model"
)

type GameService struct{}

func (s *GameService) IsGameOver(board [3][3]int) bool {
	if s.Evaluate(board) != 0 || !s.hasEmptyCells(board) {
		return true
	}
	return false
}

func (s *GameService) PlayerMove(savedGameState *model.Game, playerMoveState [3][3]int) (*model.Game, error) {
	if s.IsGameOver(savedGameState.Board) {
		return nil, errors.New("невозможно сделать ход, игра уже завершена")
	}

	if err := s.validatePlayerMove(savedGameState.Board, playerMoveState); err != nil {
		return nil, err
	}

	savedGameState.Board = playerMoveState

	if s.IsGameOver(savedGameState.Board) {
		return savedGameState, nil
	}

	s.computerMove(savedGameState)

	return savedGameState, nil
}

func (s *GameService) validatePlayerMove(savedGameState, playerMoveState [3][3]int) error {
	boardChangeCount := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if savedGameState[i][j] != playerMoveState[i][j] {
				if savedGameState[i][j] != model.Empty {
					return errors.New("нельзя изменять сделанные ходы")
				}
				if playerMoveState[i][j] != model.X {
					return errors.New("неправильный ход, можно ставить только крестики")
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

func (s *GameService) computerMove(g *model.Game) {
	bestScore := -1000
	var bestMove [2]int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == model.Empty {
				g.Board[i][j] = model.O
				score := s.minimax(g.Board, 0, false)
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

func (s *GameService) minimax(board [3][3]int, depth int, isMax bool) int {
	score := s.Evaluate(board)
	if score == 10 {
		return score - depth
	}
	if score == -10 {
		return score + depth
	}
	if !s.hasEmptyCells(board) {
		return 0
	}

	if isMax {
		best := -1000
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == model.Empty {
					board[i][j] = model.O
					score := s.minimax(board, depth+1, false)
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
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == model.Empty {
				board[i][j] = model.X
				score := s.minimax(board, depth+1, true)
				if score < best {
					best = score
				}
				board[i][j] = model.Empty
			}
		}
	}
	return best
}

func (s *GameService) Evaluate(b [3][3]int) int {
	for i := 0; i < 3; i++ {
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

func (s *GameService) hasEmptyCells(b [3][3]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == model.Empty {
				return true
			}
		}
	}
	return false
}
