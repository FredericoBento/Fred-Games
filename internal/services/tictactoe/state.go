package tictactoe

import (
	"errors"
	"log/slog"
	"strconv"
)

type Player struct {
	Username  string `json:"username"`
	Connected bool   `json:"connected"`
	Wins      int    `json:"wins"`
}

type GameState struct {
	Code    string     `json:"code"`
	Board   [3][3]int  `json:"board"`
	Player1 *Player    `json:"player1"`
	Player2 *Player    `json:"player2"`
	Turn    int        `json:"current_turn"` //odd number player2 even player1
	Status  GameStatus `json:"status"`
	Ties    int        `json:"ties"`
	Winner  int        `json:"winner"`
}

type GameStatus int

const (
	game_status_paused   = 0
	game_status_running  = 1
	game_status_finished = 2
)

var (
	ErrInvalidClient   = errors.New("Invalid client")
	ErrNoPlayerRemoved = errors.New("No player was removed")

	ErrPlayerNotFound = errors.New("Player was not found")
	ErrInvalidCell    = errors.New("Invalid Cell to make play")
)

func NewGameState(code string) *GameState {
	return &GameState{
		Code:    code,
		Player1: nil,
		Player2: nil,
		Turn:    0,
		Status:  game_status_paused,
		Winner:  0,
	}
}

func NewPlayer(username string) *Player {
	return &Player{
		Username:  username,
		Connected: true,
	}
}

func (state *GameState) AddPlayer(username string) error {
	if state.Player1 == nil {
		state.Player1 = NewPlayer(username)
	} else {
		state.Player2 = NewPlayer(username)
	}

	return nil
}

func (state *GameState) RemovePlayer(username string) error {
	if state.Player1 != nil {
		if state.Player1.Username == username {
			state.Player1 = nil
			return nil
		}
	}

	if state.Player2 != nil {
		if state.Player2.Username == username {
			state.Player2 = nil
			return nil
		}
	}

	return ErrNoPlayerRemoved
}

func (state *GameState) MakePlay(player_num int, row int, col int) error {
	var player *Player
	if player_num == 1 {
		player = state.Player1
	} else {
		player = state.Player2
	}

	if player == nil {
		return ErrPlayerNotFound
	}
	if state.Board[row][col] == 0 {
		state.Board[row][col] = player_num
	} else {
		slog.Info("Row: "+strconv.Itoa(row)+" Col: "+strconv.Itoa(col), state.Board[row][col])
		return ErrInvalidCell
	}

	return nil
}

func (state *GameState) CheckWin() (int, int, int) {
	// Check rows
	for row := 0; row < 3; row++ {
		if state.Board[row][0] != 0 && state.Board[row][0] == state.Board[row][1] && state.Board[row][1] == state.Board[row][2] {
			state.Status = game_status_finished
			state.Winner = state.Board[row][0]
			return row, 0, state.Board[row][0] // winning row, column, and value
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if state.Board[0][col] != 0 && state.Board[0][col] == state.Board[1][col] && state.Board[1][col] == state.Board[2][col] {
			state.Status = game_status_finished
			state.Winner = state.Board[0][col]
			return 0, col, state.Board[0][col] // winning row, column, and value
		}
	}

	// Check diagonal (top-left to bottom-right)
	if state.Board[0][0] != 0 && state.Board[0][0] == state.Board[1][1] && state.Board[1][1] == state.Board[2][2] {
		state.Status = game_status_finished
		state.Winner = state.Board[0][0]
		return 0, 0, state.Board[0][0] // winning row, column, and value
	}

	// Check diagonal (top-right to bottom-left)
	if state.Board[0][2] != 0 && state.Board[0][2] == state.Board[1][1] && state.Board[1][1] == state.Board[2][0] {
		state.Status = game_status_finished
		state.Winner = state.Board[0][2]
		return 0, 2, state.Board[0][2] // winning row, column, and value
	}

	// No winner, check for tie
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if state.Board[row][col] == 0 {
				// If there is an empty cell, the game is still ongoing
				return -1, -1, 0
			}
		}
	}

	state.Status = game_status_finished
	state.Winner = 0
	state.Ties += 1
	return -1, -1, 0
}

func (state *GameState) ClearBoard() {
	for row := range state.Board {
		for col := range state.Board[row] {
			state.Board[row][col] = 0
		}
	}
}

func (state *GameState) Restart(resetScoreBoard bool) {
	state.ClearBoard()
	state.Winner = 0
	state.Turn = 0
	if resetScoreBoard {
		if state.Player1 != nil {
			state.Player1.Wins = 0
		}
		if state.Player2 != nil {
			state.Player2.Wins = 0
		}

		state.Ties = 0
	}
	state.Status = game_status_running
}
