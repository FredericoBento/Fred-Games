package pong

import (
	"encoding/json"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/models"
)

type Pixel int

type Direction int

const (
	DirectionLeft  = 0
	DirectionRight = 1
)

type Ball struct {
	position  models.Vector `json:"position"`
	speed     float32       `json:"speed"`
	direction Direction     `json:"direction"`
	isRunning bool          `json:"isRunning"`
	radius    Pixel         `json:"radius"`
}

type Player struct {
	paddle      Paddle `json:"paddle"`
	isConnected bool   `json:"isConnected"`
	username    string `json:"username"`
}

type Paddle struct {
	position models.Vector `json:"position"`
	length   Pixel         `json:"length"`
	width    Pixel         `json:"width"`
}

type PongGameState struct {
	player1    *Player `json:"player1"`
	player2    *Player `json:"player2"`
	ball       *Ball   `json:"ball"`
	hasStarted bool    `json:"hasStarted"`
}

func NewPlayer(username string, paddle Paddle) *Player {
	return &Player{
		paddle:      paddle,
		isConnected: false,
		username:    username,
	}
}

func NewPaddle(position models.Vector, length Pixel, width Pixel) Paddle {
	return Paddle{
		position: position,
		length:   length,
		width:    length,
	}
}

func NewBall(position models.Vector, radius Pixel, speed float32, direction Direction, isRunning bool) *Ball {
	return &Ball{
		position:  position,
		radius:    radius,
		speed:     speed,
		direction: direction,
		isRunning: isRunning,
	}
}

func NewDefaultBall() *Ball {
	return NewBall(
		models.Vector{X: 250, Y: 250},
		3,
		3.5,
		DirectionLeft,
		false,
	)
}

func NewDefaultPaddle(position models.Vector) Paddle {
	return NewPaddle(position, 10, 2)
}

func NewPongGameState(p1 *Player, p2 *Player, ball *Ball) *PongGameState {
	return &PongGameState{
		player1:    p1,
		player2:    p2,
		ball:       ball,
		hasStarted: false,
	}
}

// func (state *PongGameState) encodeProto() ([]byte, error) {
// data, err := proto.Marshal(state)
// if err != nil {
// slog.Error("could not enconde pong game state to protobuf")
// return nil, err
// }
// return data, nil
// }

func (state *PongGameState) encodeJSON() ([]byte, error) {
	data, err := json.Marshal(state)
	if err != nil {
		slog.Error("could not enconde pong game state to JSON")
		return nil, err
	}
	return data, nil
}
