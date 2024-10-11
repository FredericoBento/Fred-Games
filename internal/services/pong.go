package services

import (
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/logger"
)

type PongService struct {
	Name   string
	Status *Status
	Log    *slog.Logger
}

func NewPongService() *PongService {
	lo, err := logger.NewServiceLogger("PongService", "", true)
	if err != nil {
		lo = slog.Default()
	}
	return &PongService{
		Name:   "PongService",
		Status: NewStatus(),
		Log:    lo,
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

func (s *PongService) GetStatus() StatusChecker {
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
