package logger

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultAppLoggerFilePath             = "./logs/apps/"
	defaultServiceLoggerFilePath         = "./logs/services/"
	defaultHandlerLoggerFilePath         = "./logs/handler/"
	defaultSQLiteRepostoryLoggerFilePath = "./logs/database/sqlite/"
)

type PrettyLogs struct {
	Time   time.Time   `json:"time"`
	Level  string      `json:"level"`
	Source slog.Source `json:"source"`
	Msg    string      `json:"msg"`
}

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
		logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}))
	}

	// logger = logger.With("App", name)

	return logger, nil
}

func GetAppLogs(appName string) ([]PrettyLogs, error) {
	logContent, err := os.ReadFile("./logs/apps/" + appName + ".log")
	if err != nil {
		return nil, err
	}

	var logs []PrettyLogs
	var logAux PrettyLogs
	for _, line := range strings.Split(string(logContent), "\n") {
		if line == "" {
			continue
		}
		err = json.Unmarshal([]byte(line), &logAux)
		if err != nil {
			return nil, err
		}
		logs = append(logs, logAux)
	}

	return logs, nil
}

func NewServiceLogger(name string, path string, showOnConsole bool) (*slog.Logger, error) {
	if path == "" {
		path = defaultServiceLoggerFilePath
		createFilePath(path)
	}

	return buildLogger(name, path, showOnConsole)

}

func NewHandlerLogger(name string, path string, showOnConsole bool) (*slog.Logger, error) {
	if path == "" {
		path = defaultHandlerLoggerFilePath
		createFilePath(path)
	}

	return buildLogger(name, path, showOnConsole)

}

func buildLogger(name string, path string, showOnConsole bool) (*slog.Logger, error) {
	file, err := os.OpenFile(path+name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger
	if showOnConsole {
		multiWriter := io.MultiWriter(file, os.Stdout)
		logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}))
	}

	return logger, nil
}

func NewRepositoryLogger(dbType, name, path string, showOnConsole bool) (*slog.Logger, error) {
	if path == "" {
		switch dbType {
		case "sqlite":
			path = defaultSQLiteRepostoryLoggerFilePath
			createFilePath(path)

		}
	}

	file, err := os.OpenFile(path+name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger
	if showOnConsole {
		multiWriter := io.MultiWriter(file, os.Stdout)
		logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}))
	}

	return logger, nil
}

func createFilePath(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}

	return os.Create(path)
}
