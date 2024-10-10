package services

import "github.com/FredericoBento/HandGame/internal/logger"

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
}
