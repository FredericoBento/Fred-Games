package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/app/admin"
	"github.com/FredericoBento/HandGame/internal/app/dummyapp"
	"github.com/FredericoBento/HandGame/internal/app/handgame"
	"github.com/FredericoBento/HandGame/internal/handler"
)

var (
	exitCode                  = 1
	exitCodeInterrupt         = 2
	ErrCouldNotReadConfigFile = errors.New("Could not read config.json, give full path")
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ApplicationConfig struct {
	Name           string `json:"name"`
	RoutePrefix    string `json:"routePrefix"`
	Active         int    `json:"active"`
	StartAtStartup int    `json:"startAtStartup"`
}

type Config struct {
	Server       ServerConfig                 `json:"server"`
	Applications map[string]ApplicationConfig `json:"applications"`
}

func main() {

	appManager := app.NewAppsManager()

	authHandler := handler.NewAuthHandler()
	homeHandler := handler.NewHomeHandler()
	adminHandler := handler.NewAdminHandler(appManager)

	serverHandlers := app.NewServerHandlers(authHandler, homeHandler, adminHandler)

	config, err := loadConfig("/home/fredarch/Documents/Github/HandGame/config.json")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}

	server := app.NewServer(
		app.WithHost(config.Server.Host),
		app.WithPort(config.Server.Port),
		app.WithHandlers(serverHandlers),
	)

	err = appManager.SetServer(server)

	if err != nil {
		slog.Error("Couldnt setup app manager server")
		os.Exit(exitCode)
	}

	for _, appConfig := range config.Applications {
		if appConfig.Active == 1 {
			app, err := createApp(appConfig, server)
			if err != nil {
				slog.Error(err.Error())
			} else {
				appManager.AddApp(app)
				if appConfig.StartAtStartup != 1 {
					err = appManager.StartApp(app.GetAppName())
					if err != nil {
						slog.Error(err.Error())
					}
				}
			}
		}
	}

	err = appManager.StartServer()
	if err != nil {
		slog.Error(err.Error())
	}

	catchInterrupt(appManager)
}

func createApp(appConfig ApplicationConfig, server *app.Server) (app.App, error) {
	switch strings.ToLower(appConfig.Name) {
	case "handgame":
		return handgame.NewHandGameApp(appConfig.Name, appConfig.RoutePrefix, server), nil

	case "admin":
		return admin.NewAdminApp(appConfig.Name, appConfig.RoutePrefix, server), nil

	case "dummyapp":
		return dummyapp.NewDummyApp(appConfig.Name, appConfig.RoutePrefix, server), nil

	default:
		return nil, errors.New("could not create app with the name " + appConfig.Name + ", app name not found")
	}
}

func loadConfig(filename string) (*Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, ErrCouldNotReadConfigFile
	}
	config := &Config{}
	if err = json.Unmarshal(raw, config); err != nil {
		return nil, err
	}

	return config, nil
}

func catchInterrupt(am *app.AppsManager) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGINT)

	<-channel
	var appsStillRunning []string
	appsStillRunning = make([]string, 0)

	for _, app := range am.Apps {
		if app.Stop() != nil {
			slog.Error("Could not stop " + app.GetAppName())
			appsStillRunning = append(appsStillRunning, app.GetAppName())
		}
	}

	numApps := len(appsStillRunning)
	if numApps > 0 {
		var answer string
		var err error
		for answer == "" || answer != "y" && answer != "n" && err != nil {
			slog.Warn("There are still running " + strconv.Itoa(numApps) + ", do you want to forcefully close them? (y/n)")
			_, err = fmt.Scan(&answer)
		}
		if answer == "y" {
			os.Exit(exitCodeInterrupt)
		} else {
			slog.Error("I dont know what to do...")
		}
	}

	os.Exit(exitCodeInterrupt)
}
