package services

import (
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/logger"
)

type HandGameService struct {
	Name   string
	Status *Status
	Log    *slog.Logger
}

const (
	logFileName = "HandgameService"
)

func NewHandGameService() *HandGameService {
	lo, err := logger.NewServiceLogger(logFileName, "", true)
	if err != nil {
		lo = slog.Default()
	}
	return &HandGameService{
		Name:   "HandGameService",
		Status: NewStatus(),
		Log:    lo,
	}
}

func (s *HandGameService) Start() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Started")
	return nil
}

func (s *HandGameService) Stop() error {
	s.Status.SetInactive()
	s.Log.Warn(s.Name + " Stopped")
	return nil
}

func (s *HandGameService) Resume() error {
	s.Status.SetActive()
	s.Log.Info(s.Name + " Resumed")
	return nil
}

func (s *HandGameService) GetStatus() StatusChecker {
	return s.Status
}

func (s *HandGameService) GetRoute() string {
	return "/handgame"
}

func (s *HandGameService) GetName() string {
	return s.Name
}

func (s *HandGameService) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetServiceLogs(logFileName)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
