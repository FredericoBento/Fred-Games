package tictactoe

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/ws"
	"github.com/gorilla/websocket"
)

type TicTacToeService struct {
	Name       string
	Status     *services.Status
	Log        *slog.Logger
	Hub        *ws.Hub
	GameStates map[string]*GameState
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  512,
		WriteBufferSize: 512,
	}

	ErrInvalidCode = errors.New("Invalid Code, Room does not exists")
)

func NewTicTacToeService() *TicTacToeService {
	lo, err := logger.NewServiceLogger("TicTacToeService", "", true)
	if err != nil {
		lo = slog.Default()
	}
	service := &TicTacToeService{
		Name:       "TicTacToeService",
		Status:     services.NewStatus(),
		Log:        lo,
		Hub:        ws.NewHub(),
		GameStates: make(map[string]*GameState),
	}
	go service.Run(service.Hub)
	return service
}

func (s *TicTacToeService) ReadMessageHandler(client *ws.Client, event ws.Event) {
	switch event.Type {
	case EventTypeCreateGame:
		s.HandleEventCreateGame(&event, client)
		break
	case EventTypeJoinGame:
		s.HandleEventJoinGame(&event, client)
		break
	case EventTypeMakePlay:
		s.HandleEventMakePlay(&event, client)
		break
	default:
		slog.Error("Unknown event received")
		return
	}
}

func (s *TicTacToeService) HandleWebSocketConnection() http.HandlerFunc {
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

func (s *TicTacToeService) Start() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Started")
	return nil
}

func (s *TicTacToeService) Stop() error {
	s.Status.SetInactive()
	s.Log.Warn(s.Name + " Stopped")
	return nil
}

func (s *TicTacToeService) Resume() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Resumed")
	return nil
}

func (s *TicTacToeService) GetStatus() services.StatusChecker {
	return s.Status
}

func (s *TicTacToeService) GetRoute() string {
	return "/tictactoe"
}

func (s *TicTacToeService) GetName() string {
	return s.Name
}

func (s *TicTacToeService) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetServiceLogs(s.Name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
