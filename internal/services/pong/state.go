package pong

import (
	"errors"

	"github.com/FredericoBento/HandGame/internal/models"
)

type Pixel int

type Paddle struct {
	Position     models.Vector2D `json:"position"`
	LastPosition models.Vector2D `json:"last_position", omitempty`
	Length       Pixel           `json:"position"`
	Width        Pixel           `json:"width"`
	Speed        float32         `json:"speed"`
}

type Player struct {
	Username string  `json:"username"`
	Score    int     `json:"points"`
	Paddle   *Paddle `json:"paddle"`
}

type Direction int

type Ball struct {
	Position     models.Vector2D `json:"position"`
	LastPosition models.Vector2D `json:"last_position"`
	Radius       float32         `json:"radius"`
	Speed        float32         `json:"speed"`
	Direction    Direction
	Dx           float32 `json:"dx"`
	Dy           float32 `json:"dy"`
}

type Canvas struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type GameStatus int

type GameState struct {
	Player1 *Player    `json:"player1"`
	Player2 *Player    `json:"player2"`
	Ball    *Ball      `json:"ball"`
	Canvas  Canvas     `json:"canvas"`
	Status  GameStatus `json:"status"`
}

const (
	ball_direction_left  = 0
	ball_direction_right = 1
	ball_direction_none  = 2

	game_status_paused   = 0
	game_status_running  = 1
	game_status_finished = 2

	default_playable_height = 360 - 50 // Canvas as 25 on top and bottom with information
	default_playable_width  = 640

	default_paddle_speed  = 300
	default_paddle_length = 40
	default_paddle_width  = 4
	default_paddle_x      = 30

	default_ball_speed  = 5
	default_ball_radius = 7
)

func NewPaddle(position models.Vector2D, length Pixel, width Pixel, speed float32) *Paddle {
	return &Paddle{
		Position:     position,
		LastPosition: position,
		Length:       length,
		Width:        width,
		Speed:        speed,
	}
}

func NewPlayer(username string, paddle *Paddle, isConnected bool) *Player {
	return &Player{
		Username: username,
		Score:    0,
		Paddle:   paddle,
	}
}

func NewBall(position models.Vector2D, radius float32, speed float32) *Ball {
	return &Ball{
		Position:     position,
		LastPosition: position,
		Radius:       radius,
		Speed:        speed,
		Direction:    ball_direction_none,
		Dx:           0,
		Dy:           0,
	}
}

func NewGameState(ball *Ball, width float32, height float32) *GameState {
	if width == 0 {
		width = default_playable_width
	}
	if height == 0 {
		height = default_playable_height
	}
	if ball == nil {
		ball = NewBall(
			models.Vector2D{
				X: width / 2,
				Y: (height + 50) / 2,
			},
			default_ball_radius,
			default_ball_speed,
		)
	}
	ball.Position.Y -= ball.Radius
	return &GameState{
		// RoomCode: code,
		Player1: nil,
		Player2: nil,
		Ball:    ball,
		Canvas: Canvas{
			Width:  width,
			Height: height,
		},
	}
}

func (state *GameState) AddPlayer(username string, paddle *Paddle) error {
	if state.Player1 != nil && state.Player2 != nil {
		return errors.New("game already has both players")
	}
	if paddle == nil {
		pos := models.Vector2D{
			X: default_paddle_x,
			Y: (360 / 2) - 20,
		}
		paddle = NewPaddle(pos, default_paddle_length, default_paddle_width, default_paddle_speed)
	}
	if state.Player1 == nil {
		state.Player1 = NewPlayer(username, paddle, true)
	} else {
		paddle.Position.X = state.Canvas.Width - 30
		state.Player2 = NewPlayer(username, paddle, true)
	}
	return nil
}

func (state *GameState) RemovePlayer(username string) error {
	if username == state.Player1.Username {
		state.Player1 = nil
		return nil
	} else {
		if username == state.Player2.Username {
			state.Player2 = nil
			return nil
		}
	}
	return errors.New("Invalid player to remove")
}

func (state *GameState) UpdatePlayer1Paddle(y float32) {
	if state.Player1 != nil {
		state.Player1.Paddle.LastPosition = state.Player1.Paddle.Position
		state.Player1.Paddle.Position.Y = y
	}
}

func (state *GameState) UpdatePlayer2Paddle(y float32) {
	if state.Player2 != nil {
		state.Player2.Paddle.LastPosition = state.Player1.Paddle.Position
		state.Player2.Paddle.Position.Y = y
	}
}
