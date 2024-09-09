package handgame

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/middleware"
)

type HandGameApp struct {
	Name        string
	RoutePrefix string
	Status      string
	Server      *app.Server
}

var (
	ErrServerRouterNotFound = errors.New("Could not find server router")
)

func NewHandGameApp(name, routePrefix string, server *app.Server) *HandGameApp {
	return &HandGameApp{
		Name:        name,
		RoutePrefix: routePrefix,
		Status:      "inactive",
		Server:      server,
	}
}

func (hga *HandGameApp) Start() error {
	slog.Info("Starting HandGame App...")

	return hga.setupRoutes()
}

func (hga *HandGameApp) setupRoutes() error {
	if hga.Server.Router == nil {
		return ErrServerRouterNotFound
	}

	appMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
	)

	hga.Server.Router.Handle(hga.RoutePrefix+"/home", appMiddlewares(hga.Server.Handlers.HomeHandler))

	slog.Info(hga.Name + " App routes have been setup")
	return nil
}

func (hga *HandGameApp) Stop() error {
	slog.Warn("Stopping HandGame App...")
	hga.Server.BlockAppRoutes(hga.RoutePrefix)

	slog.Info("HandGame App routes has stopped")
	hga.Status = "Stopped"

	return nil
}

func (hga *HandGameApp) GetAppName() string {
	return hga.Name
}
