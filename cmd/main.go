package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/FredericoBento/HandGame/internal/app"
)

var (
	handgame          app.App
	exitCode          = 1
	exitCodeInterrupt = 2
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	defer func() {
		signal.Stop(channel)
		cancel()
	}()

	go func() {
		select {
		case <-channel:
			cancel()
		case <-ctx.Done():
		}
		<-channel

		err := handgame.Stop()
		if err != nil {
			slog.Error(err.Error())
		}
		os.Exit(exitCodeInterrupt)

	}()

	if err := run(ctx); err != nil {
		slog.Error(err.Error())

		err = handgame.Stop()
		if err != nil {
			slog.Error(err.Error())
		}

		os.Exit(exitCode)
	}
}

func run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			handgame = app.NewHandGameApp()
			return handgame.Start()
		}
	}

}
