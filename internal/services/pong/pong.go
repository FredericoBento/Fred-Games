package pong

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/utils"
	"github.com/FredericoBento/HandGame/internal/ws"
	"github.com/gorilla/websocket"
)

type PongService struct {
	Name   string
	Status *services.Status
	Log    *slog.Logger
	Hub    *ws.Hub
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrInvalidCode = errors.New("Invalid Code, Room does not exists")
)

func NewPongService() *PongService {
	lo, err := logger.NewServiceLogger("PongService", "", true)
	if err != nil {
		lo = slog.Default()
	}
	service := &PongService{
		Name:   "PongService",
		Status: services.NewStatus(),
		Log:    lo,
		Hub:    ws.NewHub(),
	}

	go service.Hub.Run()
	return service
}

func (s *PongService) ReadMessageHandler(client *ws.Client, message []byte) {
	event := ws.Event{}
	err := json.Unmarshal(message, &event)
	if err != nil {
		slog.Error("Invalid message no event: " + string(message))
		return
	}
	event.From = client.Username

	switch event.Type {
	case EventTypeMessage:
		message := EventMessage{}
		err := json.Unmarshal(event.Data, &message)
		if err != nil {
			slog.Error("Invalid data for EventMessage: " + err.Error())
			return
		}
		slog.Info("Got Message from " + client.Username + ": " + message.Message)
	case EventTypeCreateRoom:
		code := utils.RandomString(4)
		_, exist := s.Hub.Rooms[code]
		for exist {
			code = utils.RandomString(4)
			_, exist = s.Hub.Rooms[code]
		}
		room := ws.NewRoom(code, 2)
		err := room.AddClient(client)
		if err != nil {
			client.SendErrorEvent(&event)
			return
		}
		client.RoomCode = code
		// event.RoomCode = code
		// client.SendSimpleMessage("/joined")
		return
	case EventTypeJoinRoom:
		data := EventJoinRoomData{}
		err := json.Unmarshal([]byte(event.Data), data)
		if err != nil {
			slog.Error("Invalid data for JoinRoomData")
			return
		}
		if _, ok := s.Hub.Rooms[data.Code]; !ok {
			client.SendErrorEventWithMessage(&event, ErrInvalidCode.Error())
		}

	default:
		slog.Error("Unknown event received")
		return
	}
}

func (s *PongService) JoinRoom() {

}

func (s *PongService) HandleWebSocketConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, isLogged := middleware.GetUserFromContext(r)
		if !isLogged {
			s.Log.Error("Error User not logged:")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Log.Error("Error upgrading to WebSocket:", err)
			return
		}

		client := ws.NewClient(conn, user.Username)
		s.Hub.Register <- client

		go client.ReadPump(s.Hub, s.ReadMessageHandler)
		go client.WritePump()
	}
}

func (s *PongService) Start() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Started")
	return nil
}

func (s *PongService) Stop() error {
	s.Status.SetInactive()
	s.Log.Warn(s.Name + " Stopped")
	return nil
}

func (s *PongService) Resume() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Resumed")
	return nil
}

func (s *PongService) GetStatus() services.StatusChecker {
	return s.Status
}

func (s *PongService) GetRoute() string {
	return "/pong"
}

func (s *PongService) GetName() string {
	return s.Name
}

func (s *PongService) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetServiceLogs(s.Name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
