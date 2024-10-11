package logger

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultServiceLoggerFilePath         = "./logs/services/"
	defaultHandlerLoggerFilePath         = "./logs/handler/"
	defaultSQLiteRepostoryLoggerFilePath = "./logs/database/sqlite/"
	defaultServerLoggerFilePath          = "./logs/httpserver/"
)

type PrettyLogs struct {
	Time   time.Time   `json:"time"`
	Level  string      `json:"level"`
	Source slog.Source `json:"source"`
	Msg    string      `json:"msg"`
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

func NewServerLogger(name string, path string, showOnConsole bool) (*slog.Logger, error) {
	if path == "" {
		path = defaultServerLoggerFilePath
		createFilePath(path)
	}

	return buildLogger(name, path, showOnConsole)
}

func NewRepositoryLogger(dbType, name, path string, showOnConsole bool) (*slog.Logger, error) {
	if path == "" {
		switch dbType {
		case "sqlite":
			path = defaultSQLiteRepostoryLoggerFilePath
			createFilePath(path)

		}
	}

	return buildLogger(name, path, showOnConsole)
}

func GetServiceLogs(serviceName string) ([]PrettyLogs, error) {
	logContent, err := os.ReadFile("./logs/services/" + serviceName + ".log")
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

func buildLogger(name string, path string, showOnConsole bool) (*slog.Logger, error) {
	file, err := os.OpenFile(path+name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger
	var multiHandler *MultiHandler

	jsonHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
	})

	if showOnConsole {
		textHandler := slog.Default().Handler()
		multiHandler = NewMultiHandler(textHandler, jsonHandler)
	} else {
		multiHandler = NewMultiHandler(jsonHandler)
	}

	logger = slog.New(multiHandler)
	return logger, nil
}

func createFilePath(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}

	return os.Create(path)
}
