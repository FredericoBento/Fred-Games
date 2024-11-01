package pong

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"time"

	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/FredericoBento/HandGame/internal/ws"
)

type EventMessage struct {
	Message string `json:"message"`
}

type EventCreateRoomData struct {
	Code string `json:"code"`
}

type EventRoomCreatedData struct {
	Code string `json:"code"`
}

type EventJoinRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

type EventPlayerJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

type EventJoinedRoomData struct {
	Code   string `json:"code"`
	Player string `json:"player"`
}

type EventPaddleMoveData struct {
	Paddle_y float32 `json:"y"`
}

const (
	EventTypeGameSettings = 0
	EventTypeMessage      = 1

	EventTypeCreateRoom  = 21
	EventTypeCreatedRoom = 22

	EventTypeJoinRoom         = 23
	EventTypeJoinedRoom       = 24
	EventTypePlayerJoinedRoom = 25

	EventTypePaddleMoved = 35

	EventTypeBallShot   = 36
	EventTypeBallUpdate = 37
	EventTypeGoal       = 38
)

var (
	ErrServerError = errors.New("Server couldnt process request")
)

func (s *PongService) HandleEventMessage(event *ws.Event, client *ws.Client) {
	message := EventMessage{}
	err := json.Unmarshal(event.Data, &message)
	if err != nil {
		slog.Error("Invalid data for EventMessage: " + err.Error())
		return
	}
	slog.Info("Got Message from " + client.Username + ": " + message.Message)
}

func (s *PongService) HandleEventCreateRoom(event *ws.Event, client *ws.Client) {
	code := utils.RandomString(4)
	_, exist := s.Hub.Rooms[code]
	for exist {
		code = utils.RandomString(4)
		_, exist = s.Hub.Rooms[code]
	}
	room := ws.NewRoom(code, 2)
	s.Hub.Rooms[code] = room

	err := room.AddClient(client)
	if err != nil {
		client.SendErrorEvent(event)
		return
	}
	createdRoomEvent := ws.NewEvent(EventTypeCreatedRoom, room.Code, "", client.Username)

	eventData := EventRoomCreatedData{
		Code: room.Code,
	}
	bytes, err := utils.EncodeJSON(eventData)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	state := NewGameState(nil, 0, 0)
	s.GameStates[code] = state
	state.AddPlayer(client.Username, nil)
	createdRoomEvent.Data = bytes
	client.SendEvent(&createdRoomEvent)
	err = room.AddClient(client)
	if err != nil {
		slog.Error("Could not add client to room")
	}
}

func (s *PongService) HandleEventJoinRoom(event *ws.Event, client *ws.Client) {
	data := EventJoinRoomData{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		slog.Error("Invalid data for JoinRoomData")
		client.SendErrorEventWithMessage(event, ErrServerError.Error())
		return
	}
	if _, ok := s.Hub.Rooms[data.Code]; !ok {
		slog.Error(err.Error())
		client.SendErrorEventWithMessage(event, ErrInvalidCode.Error())
		return
	}
	room := s.Hub.Rooms[data.Code]
	if len(room.Clients) <= 0 {
		client.SendErrorEventWithMessage(event, "Room is empty")
		return
	}
	err = room.AddClient(client)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	var otherClientUsername string
	for _, c := range room.Clients {
		if c.Username != client.Username {
			otherClientUsername = c.Username
			eventData := EventPlayerJoinedRoomData{
				Code:   data.Code,
				Player: client.Username,
			}
			slog.Info(eventData.Code)
			bytes, err := utils.EncodeJSON(eventData)
			if err != nil {
				c.SendErrorEventWithMessage(event, err.Error())
				return
			}
			event := ws.NewEvent(EventTypePlayerJoinedRoom, data.Code, "server", c.Username)
			event.Data = bytes
			c.SendEvent(&event)
		}
	}
	err = s.GameStates[room.Code].AddPlayer(client.Username, nil)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}

	eventData := EventJoinedRoomData{
		Code:   room.Code,
		Player: otherClientUsername,
	}
	bytes, err := utils.EncodeJSON(eventData)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	joinedEvent := ws.NewSimpleEvent(EventTypeJoinedRoom, client.Username)
	joinedEvent.Data = bytes
	client.SendEvent(&joinedEvent)
}

func (s *PongService) HandleEventPaddleMove(event *ws.Event, client *ws.Client) {
	data := EventPaddleMoveData{}
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		slog.Error("Invalid data for paddle pressed event")
		client.SendErrorEventWithMessage(event, ErrServerError.Error())
		return
	}
	room, ok := s.Hub.Rooms[client.RoomCode]
	if !ok {
		client.SendErrorEventWithMessage(event, "Invalid Room")
		return
	}
	state, ok := s.GameStates[room.Code]
	if !ok {
		client.SendErrorEventWithMessage(event, "State Does not exist")
		return
	}
	if state.Player1 == nil || state.Player2 == nil {
		return
	}
	if client.Username == state.Player1.Username {
		state.UpdatePlayer1Paddle(data.Paddle_y)
		data.Paddle_y = state.Player1.Paddle.Position.Y
	} else {
		state.UpdatePlayer2Paddle(data.Paddle_y)
		data.Paddle_y = state.Player2.Paddle.Position.Y
	}
	bytes, err := utils.EncodeJSON(data)
	if err != nil {
		client.SendErrorEventWithMessage(event, err.Error())
		return
	}
	for _, c := range room.Clients {
		if c.Username != client.Username {
			moveEvent := ws.NewSimpleEvent(event.Type, c.Username)
			moveEvent.Data = bytes
			go func(targetClient *ws.Client, moveEvent ws.Event) {
				targetClient.SendEvent(&moveEvent)
			}(c, moveEvent)
		}
	}
}

func (s *PongService) HandleEventBallShot(event *ws.Event, client *ws.Client) {
	state, ok := s.GameStates[client.RoomCode]
	if !ok {
		client.SendErrorEventWithMessage(event, "Invalid Room")
		return
	}
	if state.Ball.Direction != ball_direction_none {
		return
	}

	if state.Player1.Username == client.Username {
		state.Ball.Direction = ball_direction_right
		state.Ball.Dx = state.Ball.Speed
	} else {
		state.Ball.Direction = ball_direction_left
		state.Ball.Dx = -state.Ball.Speed
	}

	go s.UpdateBall(state, client.RoomCode)
}

const deltaTime = 0.016

func (s *PongService) UpdateBall(state *GameState, code string) {
	if state == nil {
		return
	}
	for state.Ball.Direction != ball_direction_none {
		time.Sleep(10 * time.Millisecond)

		if state.Ball.is_collision(state.Player1.Paddle) {
			state.Ball.handle_collision(state.Player1.Paddle)
		} else {
			if state.Ball.is_collision(state.Player2.Paddle) {
				state.Ball.handle_collision(state.Player2.Paddle)
			}
		}

		state.Ball.check_wall_collision(25, state.Canvas.Height+25)
		// moe := float32(4) // margin of error
		if state.Ball.Position.X < state.Player1.Paddle.Position.X+float32(state.Player1.Paddle.Width) {
			if state.Ball.Position.X-state.Ball.Radius <= 0 {
				state.Player2.Score += 1
				state.Ball.recenter(state.Canvas.Width, state.Canvas.Height+50)
				go s.UpdatePoints(state, code)
			}
		} else {
			if state.Ball.Position.X > state.Player2.Paddle.Position.X+float32(state.Player1.Paddle.Width) {
				if state.Ball.Position.X+state.Ball.Radius >= state.Canvas.Width {
					state.Player1.Score += 1
					state.Ball.recenter(state.Canvas.Width, state.Canvas.Height+50)
					go s.UpdatePoints(state, code)
					s.UpdateBall(state, code)
				}
			}
		}

		state.Ball.Position.X += state.Ball.Dx
		state.Ball.Position.Y += state.Ball.Dy

		data, err := utils.EncodeJSON(state.Ball.Position)
		if err == nil {
			go func() {
				for _, client := range s.Hub.Rooms[code].Clients {
					event := ws.NewSimpleEvent(EventTypeBallUpdate, client.Username)
					event.Data = data
					client.SendEvent(&event)
				}

			}()
		}
	}
}

func two_points_distance(point1 *models.Vector2D, point2 *models.Vector2D) float32 {
	return float32(math.Sqrt(float64((math.Pow((float64(point1.X)-float64(point2.X)), 2) + math.Pow((float64(point1.Y)-float64(point2.Y)), 2)))))
}

func (ball *Ball) is_collision(paddle *Paddle) bool {

	// Calcular os limites da raquete
	paddle_x1 := paddle.Position.X
	paddle_x2 := paddle.Position.X + float32(paddle.Width)
	paddle_y1 := paddle.Position.Y
	paddle_y2 := paddle.Position.Y + float32(paddle.Length)

	// Encontrar o ponto mais próximo na raquete em relação à bola
	closest_x := math.Max(float64(paddle_x1), math.Min(float64(ball.Position.X), float64(paddle_x2)))
	closest_y := math.Max(float64(paddle_y1), math.Min(float64(ball.Position.Y), float64(paddle_y2)))

	// Calcular a distância do centro da bola até o ponto mais próximo
	distance := math.Sqrt(math.Pow(float64(ball.Position.X)-closest_x, 2) + math.Pow(float64(ball.Position.Y)-closest_y, 2))

	// Verificar se a distância é menor ou igual ao raio da bola
	return distance <= float64(ball.Radius)
	// paddle_x1 := paddle.Position.X                         // top left
	// paddle_x2 := paddle.Position.X + float32(paddle.Width) // bottom right

	// paddle_y1 := paddle.Position.Y
	// paddle_y2 := paddle.Position.Y + float32(paddle.Length)

	// if ball.Position.X < paddle_x1 || ball.Position.X > paddle_x2 || ball.Position.Y < paddle_y1 || ball.Position.Y > paddle_y2 {
	// 	return false
	// }

	// closest_x := math.Max(float64(paddle_x1), math.Min(float64(ball.Position.X), float64(paddle_x2)))
	// closest_y := math.Max(float64(paddle_y1), math.Min(float64(ball.Position.Y), float64(paddle_y2)))

	// distance := math.Sqrt(math.Pow(float64(ball.Position.X)-closest_x, 2) + math.Pow(float64(ball.Position.Y)-closest_y, 2))
	// return distance <= float64(ball.Radius)

	// if ball.Position.X+ball.Radius < paddle.Position.X ||
	// 	ball.Position.X-ball.Radius > paddle.Position.X+float32(paddle.Width) {
	// 	return false
	// }

	// if ball.Position.Y+ball.Radius < paddle.Position.Y ||
	// 	ball.Position.Y-ball.Radius > paddle.Position.Y+float32(paddle.Length) {
	// 	return false
	// }
	// slog.Info("COLLISION")

	// return true
}

func (ball *Ball) recenter(width float32, height float32) {
	ball.Direction = ball_direction_none
	ball.Position.X = (width / 2)
	ball.Position.Y = (height / 2)
	ball.Dx = 0
	ball.Dy = 0
}

func (ball *Ball) check_wall_collision(top_y float32, bottom_y float32) {
	angle_modifer := float32(1.1)
	if ball.Position.Y-ball.Radius <= top_y {
		ball.Dy = -ball.Dy * angle_modifer
		ball.Position.Y = top_y + ball.Radius
	}

	if ball.Position.Y+ball.Radius >= bottom_y {
		ball.Dy = -ball.Dy * angle_modifer
		ball.Position.Y = bottom_y - ball.Radius
	}
}
func (ball *Ball) handle_collision(paddle *Paddle) {
	ball.invert_direction()
	ball.Dx = -ball.Dx

	offset := (ball.Position.Y - paddle.Position.Y) / float32(paddle.Length/2)
	ball.Dy += offset * 0.9

	// speed := float32(math.Sqrt(float64(ball.Dx*ball.Dx + ball.Dy*ball.Dy)))
	// ball.Dx = (ball.Dx / speed) * ball.Speed
	// ball.Dy = (ball.Dy / speed) * ball.Speed
}

func (ball *Ball) invert_direction() {
	if ball.Direction == ball_direction_left {
		ball.Direction = ball_direction_right
	} else {
		if ball.Direction == ball_direction_right {
			ball.Direction = ball_direction_left
		}
	}
}

func (s *PongService) UpdatePoints(state *GameState, code string) {
	type Points struct {
		Player1Score int `json:"player1_score"`
		Player2Score int `json:"player2_score"`
	}
	points := Points{
		Player1Score: state.Player1.Score,
		Player2Score: state.Player2.Score,
	}
	data, err := utils.EncodeJSON(points)
	if err == nil {
		for _, client := range s.Hub.Rooms[code].Clients {
			event := ws.NewSimpleEvent(EventTypeGoal, client.Username)
			event.Data = data
			client.SendEvent(&event)
		}
	}
}
