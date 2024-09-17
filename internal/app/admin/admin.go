package admin

import (
	"errors"
	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"log/slog"
)

type AdminApp struct {
	name        string
	routePrefix string
	server      *app.Server
	status      *app.AppStatus
	log         *slog.Logger
}

var (
	ErrAppIsAlreadyActive   = errors.New("Admin App is already active")
	ErrCouldNotCreateLogger = errors.New("could not create app logger")
)

func NewAdminApp(name, routePrefix string, server *app.Server) *AdminApp {
	lo, err := app.NewAppLogger(name, "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error(), err)
		lo = slog.Default()
	}

	return &AdminApp{
		name:        name,
		routePrefix: routePrefix,
		server:      server,
		status:      app.NewAppStatus(),
		log:         lo,
	}
}

func (aa *AdminApp) Start() error {
	if aa.status.IsActive() {
		return ErrAppIsAlreadyActive
	}

	aa.log.Info("Starting " + aa.name + " App...")
	if err := aa.setupRoutes(); err != nil {
		aa.status.SetError(err)
		return err
	}
	aa.status.SetActive()

	aa.log.Info(" - Ok")
	return nil
}

func (aa *AdminApp) Stop() error {
	aa.log.Info("Stopping" + aa.name + " App...")
	aa.server.BlockAppRoutes(aa.routePrefix)

	aa.log.Info(" - Ok")
	aa.status.SetInactive()

	return nil
}

func (aa *AdminApp) setupRoutes() error {
	if aa.server.Router == nil {
		return app.ErrServerRouterNotFound
	}

	appMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
	)

	aa.server.Router.Handle(aa.routePrefix+"/dashboard", appMiddlewares(aa.server.Handlers.AdminHandler))
	aa.server.Router.Handle(aa.routePrefix+"/", appMiddlewares(aa.server.Handlers.AdminHandler))

	return nil
}

func (aa *AdminApp) Resume() error {
	aa.log.Warn("Resuming " + aa.name + "App...")
	aa.server.UnblockAppRoutes(aa.routePrefix)

	aa.log.Info(" - Ok")
	aa.status.SetActive()

	return nil
}

func (aa *AdminApp) GetAppName() string {
	return aa.name
}

func (aa *AdminApp) GetStatus() app.AppStatusChecker {
	return aa.status
}
