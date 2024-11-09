package pong

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

type PongService struct {
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

func NewPongService() *PongService {
	lo, err := logger.NewServiceLogger("PongService", "", true)
	if err != nil {
		lo = slog.Default()
	}
	service := &PongService{
		Name:       "PongService",
		Status:     services.NewStatus(),
		Log:        lo,
		Hub:        ws.NewHub(),
		GameStates: make(map[string]*GameState),
	}

	// go service.Hub.Run()
	go service.Run(service.Hub)
	return service
}

func (s *PongService) ReadMessageHandler(client *ws.Client, event ws.Event) {
	switch event.Type {
	case ws.EventTypePing:
		ws.HandleEventPing(&event, client)
		return

	case EventTypeMessage:
		s.HandleEventMessage(&event, client)
		return

	case EventTypeCreateRoom:
		s.HandleEventCreateRoom(&event, client)
		return

	case EventTypeJoinRoom:
		s.HandleEventJoinRoom(&event, client)
		return

	case EventTypePaddleMoved:
		s.HandleEventPaddleMove(&event, client)
		return

	case EventTypeBallShot:
		s.HandleEventBallShot(&event, client)
		return

	default:
		slog.Error("Unknown event received")
		return
	}
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

		// go func() {
		// 	ticker := time.NewTicker(5 * time.Second)
		// 	defer ticker.Stop()
		// 	for {
		// 		<-ticker.C
		// 		if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		// 			s.Log.Error("Failed to send ping:", err)
		// 			return
		// 		}
		// 	}
		// }()

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
