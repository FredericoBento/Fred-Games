package app

import (
	"errors"
	"log/slog"
	"strings"
)

type App interface {
	Start() error
	Stop() error
	GetAppName() string
}

type AppsManager struct {
	Apps   map[string]App
	Server *Server
}

type AppsManagerOption func(*AppsManager)

const (
	statusInactive           = "inactive"
	statusActive             = "active"
	statusInactiveWithErrors = "inactive and errors"
	statusActiveWithErrors   = "active and errors"
)

var (
	ErrAppNameNotFound = errors.New("Could not find requested app by such name")
	ErrAppNameInUse    = errors.New("App name is alread in use")
)

func NewAppsManager(opts ...AppsManagerOption) *AppsManager {
	appsManager := &AppsManager{
		Apps:   make(map[string]App),
		Server: nil,
	}

	for _, option := range opts {
		option(appsManager)
	}

	return appsManager
}

func WithApp(app App) AppsManagerOption {
	return func(am *AppsManager) {
		am.Apps[app.GetAppName()] = app
	}
}

func WithServer(server *Server) AppsManagerOption {
	return func(am *AppsManager) {
		am.Server = server
	}
}

func (am *AppsManager) AddApp(app App) error {
	if am.Apps[app.GetAppName()] != nil {
		return ErrAppNameInUse
	}

	am.Apps[app.GetAppName()] = app
	return nil
}

func (am *AppsManager) StartApp(appName string) error {
	if am.Apps[appName] != nil {
		return am.Apps[appName].Start()
	}
	return ErrAppNameNotFound
}

func (am *AppsManager) StopApp(appName string) error {
	for _, app := range am.Apps {
		name := app.GetAppName()
		name = strings.ToLower(name)
		if name == strings.ToLower(appName) {
			return app.Stop()
		}
	}
	return ErrAppNameNotFound
}

func (am *AppsManager) StartAll() ([]string, error) {
	unableToStart := make([]string, 0)
	for _, app := range am.Apps {
		err := app.Start()
		if err != nil {
			slog.Error("Could not start " + app.GetAppName() + " App: " + err.Error())
			unableToStart = append(unableToStart, app.GetAppName())
		}
	}
	if len(unableToStart) > 0 {
		return unableToStart, errors.New("Some apps could not be started")
	}
	return nil, nil
}

func (am *AppsManager) StartServer() error {
	err := am.Server.Init()
	if err != nil {
		return err
	}

	return am.Server.Run()
}
