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
	Name        string `json:"name"`
	RoutePrefix string `json:"routePrefix"`
	Active      int    `json:"active"`
}

type Config struct {
	Server       ServerConfig                 `json:"server"`
	Applications map[string]ApplicationConfig `json:"applications"`
}

func main() {

	authHandler := handler.NewAuthHandler()
	homeHandler := handler.NewHomeHandler()

	serverHandlers := app.NewServerHandlers(authHandler, homeHandler)

	config, err := loadConfig("/home/fredarch/Documents/Github/HandGame/config.json")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}

	slog.Info("Port: " + strconv.Itoa(config.Server.Port))

	server := app.NewServer(
		app.WithHost(config.Server.Host),
		app.WithPort(config.Server.Port),
		app.WithHandlers(serverHandlers),
	)

	appManager := app.NewAppsManager(
		app.WithServer(server),
	)

	for _, appConfig := range config.Applications {
		if appConfig.Active == 1 {
			app := createApp(appConfig, server)
			appManager.AddApp(app)
		}
	}

	go catchInterrupt(appManager)

	_, err = appManager.StartAll()
	if err != nil {
		slog.Error(err.Error())
	} else {
		err = appManager.StartServer()
	}

	slog.Warn("Going to block Handgame")
	err = appManager.StopApp("handgame")
	if err != nil {
		slog.Error("Could not block Handgame")
	} else {
		slog.Info("Handgame routes do not work anymore")
	}

	catchInterrupt(appManager)
}

func createApp(appConfig ApplicationConfig, server *app.Server) app.App {
	switch strings.ToLower(appConfig.Name) {
	case "handgame":
		return handgame.NewHandGameApp(appConfig.Name, appConfig.RoutePrefix, server)
	}
	return nil
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
