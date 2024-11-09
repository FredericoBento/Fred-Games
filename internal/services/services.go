package services

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/ws"
)

type Service interface {
	// GetStatus() StatusChecker
	// GetLogs() ([]logger.PrettyLogs, error)
}

type GameService interface {
	Start() error
	Stop() error
	Resume() error
	GetName() string
	GetStatus() StatusChecker
	GetRoute() string
	GetLogs() ([]logger.PrettyLogs, error)
	HandleWebSocketConnection() http.HandlerFunc
	ReadMessageHandler(client *ws.Client, event ws.Event)
}
