package main

import (
	"database/sql"
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
	"github.com/FredericoBento/HandGame/internal/database/sqlite"
	"github.com/FredericoBento/HandGame/internal/handler"
	"github.com/FredericoBento/HandGame/internal/services"

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

func main() {
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

	userRepository := sqlite.NewSQLiteUserRepository(db)

	userService := services.NewUserService(userRepository)

	authHandler := handler.NewAuthHandler()
	homeHandler := handler.NewHomeHandler()
	adminHandler := handler.NewAdminHandler(appManager, userService)

	serverHandlers := app.NewServerHandlers(authHandler, homeHandler, adminHandler)

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
	}

	return db, nil
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
