package admin_service

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/server"
	"github.com/FredericoBento/HandGame/internal/services"
)

type AdminService struct {
	Name         string
	Status       *services.Status
	GameServices []services.GameService
	Log          *slog.Logger
	Server       *server.Server
}

var (
	ErrNoServerProvided   = errors.New("server was not provided")
	ErrGameNotFound       = errors.New("game was not found")
	ErrCouldNotStopGame   = errors.New("game could not be stopped")
	ErrCouldNotStartGame  = errors.New("game could not be started")
	ErrCouldNotResumeGame = errors.New("game could not be resumed")
	ErrGameServiceUnknown = errors.New("unknown game service name")
)

func NewAdminService(server *server.Server, gameServices []services.GameService) *AdminService {
	lo, err := logger.NewServiceLogger("AdminService", "", true)
	if err != nil {
		lo = slog.Default()
	}
	return &AdminService{
		Name:         "AdminService",
		Status:       services.NewStatus(),
		GameServices: gameServices,
		Log:          lo,
		Server:       server,
	}
}

func (s *AdminService) GetName() string {
	return s.Name
}

func (s *AdminService) GetStatus() services.Status {
	return *s.Status
}

func (s *AdminService) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetServiceLogs(s.Name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *AdminService) StartGame(name string) error {
	game, ok := s.GetGame(name)
	if ok {
		err := game.Start()
		if err != nil {
			s.Log.Error(err.Error())
			return ErrCouldNotStartGame
		}
		s.Server.UnblockRoutes(game.GetRoute())
		return nil
	}

	return ErrGameNotFound
}

func (s *AdminService) StopGame(name string) error {
	game, ok := s.GetGame(name)
	if ok {
		err := game.Stop()
		if err != nil {
			s.Log.Error(err.Error())
			return ErrCouldNotStopGame
		}
		s.Server.BlockRoutes(game.GetRoute())
		return nil
	}

	return ErrGameNotFound
}

func (s *AdminService) ResumeGame(name string) error {
	game, ok := s.GetGame(name)
	if ok {
		err := game.Resume()
		if err != nil {
			s.Log.Error(err.Error())
			return ErrCouldNotResumeGame
		}
		s.Server.UnblockRoutes(game.GetRoute())
		return nil
	}

	return ErrGameNotFound
}

func (s *AdminService) GetGame(gameName string) (services.GameService, bool) {
	for _, service := range s.GameServices {
		if service.GetName() == gameName {
			return service, true
		}
	}
	return nil, false
}
