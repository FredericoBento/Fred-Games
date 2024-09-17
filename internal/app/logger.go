package app

import (
	"io"
	"log/slog"
	"os"
)

const (
	defaultAppLoggerFilePath = "./logs/apps/"
)

func NewAppLogger(name string, filePath string, showOnConsole bool) (*slog.Logger, error) {
	if filePath == "" {
		filePath = defaultAppLoggerFilePath
	}

	file, err := os.OpenFile(filePath+name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger
	if showOnConsole {
		multiWriter := io.MultiWriter(file, os.Stdout)
		logger = slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}))
	}

	// logger = logger.With("App", name)

	return logger, nil
}
