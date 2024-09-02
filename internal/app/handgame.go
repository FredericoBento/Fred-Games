package app

import (
	"github.com/FredericoBento/HandGame/internal/handler"
	"log/slog"
)

type App interface {
	Start() error
	Stop() error
}

type HandGameApp struct {
	Server *Server
}

func NewHandGameApp() App {
	return &HandGameApp{}
}

func (hga *HandGameApp) Start() error {
	slog.Info("Starting HandGame App...")
	authHandler := handler.NewAuthHandler()
	homeHandler := handler.NewHomeHandler()

	serverHandlers := NewServerHandlers(authHandler, homeHandler)
	hga.Server = NewServer(
		WithPort(8080),
		WithHandlers(serverHandlers),
	)

	err := hga.Server.Init()
	if err != nil {
		hga.Stop()
		return err
	}

	err = hga.Server.Run()
	if err == nil {
		slog.Info("HandGame App has started ")
	} else {
		slog.Error("HandGame App could not start")
	}

	return err
}

func (hga *HandGameApp) Stop() error {
	slog.Warn("Stopping HandGame App...")
	err := hga.Server.Shutdown()

	if err == nil {
		slog.Info("HandGame App has stopped")
	} else {
		slog.Error("HandGame App could not be stopped")
	}

	return err
}
