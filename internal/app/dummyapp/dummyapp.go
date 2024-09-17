package dummyapp

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/app"
)

var (
	ErrHasNotStartedYet     = errors.New("app hasnt started yet")
	ErrCouldNotCreateLogger = errors.New("could not create app logger")
)

type DummyApp struct {
	name        string
	routePrefix string
	status      *app.AppStatus
	server      *app.Server
	log         *slog.Logger
}

func NewDummyApp(name, routePrefix string, server *app.Server) *DummyApp {
	lo, err := app.NewAppLogger(name, "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error(), err)
		lo = slog.Default()
	}

	return &DummyApp{
		name:        name,
		routePrefix: routePrefix,
		status:      app.NewAppStatus(),
		server:      server,
		log:         lo,
	}
}

func (da *DummyApp) Start() error {
	da.log.Info("Starting " + da.name + " App...")
	da.status.SetActive()
	da.log.Info(" - Ok")
	return nil
}

func (da *DummyApp) Stop() error {
	da.log.Info("Stopping " + da.name + " App...")
	da.status.SetInactive()
	da.log.Info(" - Ok")
	return nil
}

func (da *DummyApp) Resume() error {
	if da.status.HasStartedOnce() == false {
		return ErrHasNotStartedYet
	}
	da.log.Info("Resuming %d App...", da.name)
	da.status.SetActive()
	da.log.Info(" - Ok")
	return nil
}

func (da *DummyApp) GetStatus() app.AppStatusChecker {
	return da.status
}

func (da *DummyApp) GetAppName() string {
	return da.name
}
