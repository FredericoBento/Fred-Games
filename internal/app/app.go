package app

// import (
// 	"errors"
// 	"log/slog"
// 	"sort"
// 	"strings"

// 	"github.com/FredericoBento/HandGame/internal/logger"
// )

// type App interface {
// 	Start() error
// 	Stop() error
// 	Resume() error
// 	GetName() string
// 	GetRoute() string
// 	GetStatus() AppStatusChecker
// 	GetLogs() ([]logger.PrettyLogs, error)
// }

// type AppStatusChecker interface {
// 	IsActive() bool
// 	IsInactive() bool
// 	SetActive()
// 	SetInactive()
// 	SetError(error)
// 	GetErrors() []error
// 	HasStartedOnce() bool
// 	HasErrors() bool
// }

// type AppStatus struct {
// 	value          string
// 	statusErrors   []error
// 	hasStartedOnce bool
// }

// type AppsManager struct {
// 	Apps   map[string]App
// 	Server *Server
// }

// type AppsManagerOption func(*AppsManager)

// const (
// 	statusInactive           = "inactive"
// 	statusActive             = "active"
// 	statusInactiveWithErrors = "inactive with errors"
// 	statusActiveWithErrors   = "active with errors"
// )

// var (
// 	ErrAppNameNotFound = errors.New("Could not find requested app by such name")
// 	ErrAppNameInUse    = errors.New("App name is alread in use")
// )

// func NewAppStatus() *AppStatus {
// 	return &AppStatus{
// 		value:          "inactive",
// 		statusErrors:   make([]error, 0),
// 		hasStartedOnce: false,
// 	}
// }

// func (as *AppStatus) IsActive() bool {
// 	return as.value == statusActive
// }

// func (as *AppStatus) IsInactive() bool {
// 	return as.value == statusInactive
// }

// func (as *AppStatus) SetInactive() {
// 	as.value = statusInactive
// }

// func (as *AppStatus) SetActive() {
// 	as.value = statusActive
// 	as.hasStartedOnce = true
// }

// func (as *AppStatus) SetError(e error) {
// 	as.statusErrors = append(as.statusErrors, e)
// }

// func (as *AppStatus) GetErrors() []error {
// 	return as.statusErrors
// }

// func (as *AppStatus) HasErrors() bool {
// 	return len(as.statusErrors) > 0
// }

// func (as *AppStatus) HasStartedOnce() bool {
// 	return as.hasStartedOnce
// }

// func NewAppsManager(opts ...AppsManagerOption) *AppsManager {
// 	appsManager := &AppsManager{
// 		Apps:   make(map[string]App),
// 		Server: nil,
// 	}

// 	for _, option := range opts {
// 		option(appsManager)
// 	}

// 	return appsManager
// }

// func WithApp(app App) AppsManagerOption {
// 	return func(am *AppsManager) {
// 		am.Apps[app.GetName()] = app
// 	}
// }

// func WithServer(server *Server) AppsManagerOption {
// 	return func(am *AppsManager) {
// 		am.Server = server
// 	}
// }

// func (am *AppsManager) AddApp(app App) error {
// 	if am.Apps[app.GetName()] != nil {
// 		return ErrAppNameInUse
// 	}

// 	am.Apps[app.GetName()] = app
// 	return nil
// }

// func (am *AppsManager) SetServer(server *Server) error {
// 	am.Server = server
// 	return nil
// }

// func (am *AppsManager) StartApp(appName string) error {
// 	if am.Apps[appName] != nil {
// 		return am.Apps[appName].Start()
// 	}
// 	return ErrAppNameNotFound
// }

// func (am *AppsManager) StopApp(appName string) error {
// 	for _, app := range am.Apps {
// 		name := app.GetName()
// 		name = strings.ToLower(name)
// 		if name == strings.ToLower(appName) {
// 			return app.Stop()
// 		}
// 	}
// 	return ErrAppNameNotFound
// }

// func (am *AppsManager) ResumeApp(appName string) error {
// 	for _, app := range am.Apps {
// 		name := app.GetName()
// 		name = strings.ToLower(name)
// 		if name == strings.ToLower(appName) {
// 			return app.Resume()
// 		}
// 	}
// 	return ErrAppNameNotFound
// }

// func (am *AppsManager) StartAll() ([]string, error) {
// 	unableToStart := make([]string, 0)
// 	for _, app := range am.Apps {
// 		err := app.Start()
// 		if err != nil {
// 			slog.Error("Could not start " + app.GetName() + " App: " + err.Error())
// 			unableToStart = append(unableToStart, app.GetName())
// 		}
// 	}
// 	if len(unableToStart) > 0 {
// 		return unableToStart, errors.New("Some apps could not be started")
// 	}
// 	return nil, nil
// }

// func (am *AppsManager) StartServer() error {
// 	err := am.Server.Init()
// 	if err != nil {
// 		return err
// 	}

// 	return am.Server.Run()
// }

// func (am *AppsManager) GetAppsSortedAlphabetic() []App {
// 	apps := make([]App, 0, len(am.Apps))
// 	keys := make([]string, 0, len(am.Apps))

// 	for k := range am.Apps {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	for _, k := range keys {
// 		apps = append(apps, am.Apps[k])
// 	}

// 	return apps
// }
