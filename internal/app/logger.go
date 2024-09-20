package app

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	defaultAppLoggerFilePath = "./logs/apps/"
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
