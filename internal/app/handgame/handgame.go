package handgame

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/models"
)

var (
	ErrServerRouterNotFound = errors.New("Could not find server router")
	ErrHasNotStartedYet     = errors.New("app hasnt started yet")
	ErrCouldNotCreateLogger = errors.New("could not create app logger")
)

type HandGameApp struct {
	name        string
	routePrefix string
	status      *app.AppStatus
	server      *app.Server
	log         *slog.Logger
}

func NewHandGameApp(name, routePrefix string, server *app.Server) *HandGameApp {
	lo, err := logger.NewAppLogger(name, "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error() + " " + err.Error())
		lo = slog.Default()
	}

	return &HandGameApp{
		name:        name,
		routePrefix: routePrefix,
		status:      app.NewAppStatus(),
		server:      server,
		log:         lo,
	}
}

func (a *HandGameApp) Start() error {
	a.log.Info("Starting HandGame App...")
	err := a.setupRoutes()
	if err != nil {
		a.log.Error(" - Failed")
		a.log.Error(err.Error())
		return err
	}
	a.log.Info(" - Ok")
	a.status.SetActive()
	return nil
}

func (a *HandGameApp) setupRoutes() error {
	if a.server.Router == nil {
		return ErrServerRouterNotFound
	}

	appMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
		middleware.RequiredLogged,
	)

	a.server.Router.Handle(a.routePrefix+"/home", appMiddlewares(a.server.Handlers.HandGameHandler))

	return nil
}

func (a *HandGameApp) Stop() error {
	a.log.Warn("Stopping HandGame App...")
	a.server.BlockAppRoutes(a.routePrefix)

	a.log.Info(" - Ok")
	a.status.SetInactive()

	return nil
}

func (a *HandGameApp) Resume() error {
	if a.status.HasStartedOnce() == false {
		a.log.Error(ErrHasNotStartedYet.Error())
		return ErrHasNotStartedYet
	}
	a.log.Warn("Resuming HandGame App...")
	a.server.UnblockAppRoutes(a.routePrefix)

	a.log.Info(" - Ok")
	a.status.SetActive()

	return nil
}

func (a *HandGameApp) GetName() string {
	return a.name
}

func (a *HandGameApp) GetRoute() string {
	return a.routePrefix
}

func (aa *HandGameApp) GetStatus() app.AppStatusChecker {
	return aa.status
}

func (da *HandGameApp) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetAppLogs(da.name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a *HandGameApp) GetRooms() []models.Room {
	rooms := []models.Room{
		{
			ID:   0,
			Name: "Room1",
		},
		{
			ID:   1,
			Name: "Room2",
		},
	}

	return rooms
}
