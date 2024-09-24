package handgame

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/middleware"
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
	lo, err := app.NewAppLogger(name, "", false)
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

func (hga *HandGameApp) Start() error {
	hga.log.Info("Starting HandGame App...")
	err := hga.setupRoutes()
	if err != nil {
		hga.log.Error(" - Failed")
		hga.log.Error(err.Error())
		return err
	}
	hga.log.Info(" - Ok")
	hga.status.SetActive()
	return nil
}

func (hga *HandGameApp) setupRoutes() error {
	if hga.server.Router == nil {
		return ErrServerRouterNotFound
	}

	appMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
		middleware.RequiredLogged,
	)

	hga.server.Router.Handle(hga.routePrefix+"/home", appMiddlewares(hga.server.Handlers.HomeHandler))

	return nil
}

func (hga *HandGameApp) Stop() error {
	hga.log.Warn("Stopping HandGame App...")
	hga.server.BlockAppRoutes(hga.routePrefix)

	hga.log.Info(" - Ok")
	hga.status.SetInactive()

	return nil
}

func (hga *HandGameApp) Resume() error {
	if hga.status.HasStartedOnce() == false {
		hga.log.Error(ErrHasNotStartedYet.Error())
		return ErrHasNotStartedYet
	}
	hga.log.Warn("Resuming HandGame App...")
	hga.server.UnblockAppRoutes(hga.routePrefix)

	hga.log.Info(" - Ok")
	hga.status.SetActive()

	return nil
}

func (hga *HandGameApp) GetAppName() string {
	return hga.name
}

func (aa *HandGameApp) GetStatus() app.AppStatusChecker {
	return aa.status
}

func (da *HandGameApp) GetLogs() ([]app.PrettyLogs, error) {
	logs, err := app.GetAppLogs(da.name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
