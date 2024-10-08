package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/app/handgame"
	"github.com/FredericoBento/HandGame/internal/app/pong"
	"github.com/FredericoBento/HandGame/internal/database/repository"
	"github.com/FredericoBento/HandGame/internal/database/sqlite"
	"github.com/FredericoBento/HandGame/internal/handler"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/services"

	_ "net/http/pprof"

	_ "github.com/mattn/go-sqlite3"
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

type DatabaseConfig struct {
	Type string `json:"type"`
}

type ApplicationConfig struct {
	Name           string `json:"name"`
	RoutePrefix    string `json:"routePrefix"`
	Active         int    `json:"active"`
	StartAtStartup int    `json:"startAtStartup"`
}

type Config struct {
	Server       ServerConfig                 `json:"server"`
	Database     DatabaseConfig               `son:"database"`
	Applications map[string]ApplicationConfig `json:"applications"`
}

const (
	dbFile = "./simple.db"
)

var (
	handGameHandler *handler.HandGameHandler
	handGameApp     *handgame.HandGameApp

	pongHandler *handler.PongHandler
	pongApp     *pong.PongApp
)

func main() {
	pprofRun()
	config, err := loadConfig("/home/fredarch/Documents/Github/HandGame/config.json")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}

	db, err := getDB(config.Database)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}
	defer db.Close()

	appManager := app.NewAppsManager()

	userRepository := repository.NewSQLiteUserRepository(db)

	userService := services.NewUserService(userRepository, time.Minute*10)
	authService := services.NewAuthService(userService)

	middleware.SetAuthService(authService)

	authHandler := handler.NewAuthHandler(authService, userService)
	adminHandler := handler.NewAdminHandler(appManager, userService)
	homeHandler := handler.NewHomeHandler(appManager, authService)

	handGameHandler = handler.NewHandGameHandler(handGameApp)
	pongHandler = handler.NewPongHandler(pongApp, authService)

	serverHandlers := app.NewServerHandlers(authHandler, adminHandler, homeHandler, handGameHandler, pongHandler)

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
				if appConfig.StartAtStartup == 1 {
					err = appManager.StartApp(app.GetName())
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

func getDB(databaseConfig DatabaseConfig) (db *sql.DB, err error) {
	switch databaseConfig.Type {
	case "sqlite":
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}
		err = sqlite.CreateTables(db)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("database type/driver not found")
	}

	return db, nil
}

func createApp(appConfig ApplicationConfig, server *app.Server) (app.App, error) {
	switch strings.ToLower(appConfig.Name) {
	case "handgame":
		handGameApp = handgame.NewHandGameApp(appConfig.Name, appConfig.RoutePrefix, server)
		return handGameApp, nil

	case "pong":
		pongApp = pong.NewPongApp(appConfig.Name, appConfig.RoutePrefix, server)
		return pongApp, nil

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
			slog.Error("Could not stop " + app.GetName())
			appsStillRunning = append(appsStillRunning, app.GetName())
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

func pprofRun() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
