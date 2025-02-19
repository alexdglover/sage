package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	FilePath         string
	UserLocalTime    bool
	FileMaxSizeInMB  int
	FileMaxBackups   int
	FileMaxAgeInDays int
	FileCompress     bool
	LogLevel         slog.Level
}

var once sync.Once

var logger *slog.Logger

var loggerConfig LoggerConfig

func init() {

	loggerConfig = LoggerConfig{
		FilePath:         "logs/logs.log",
		UserLocalTime:    false,
		FileMaxSizeInMB:  10,
		FileMaxBackups:   3,
		FileMaxAgeInDays: 30,
		FileCompress:     false,
		LogLevel:         slog.LevelDebug,
	}
}

func Get() *slog.Logger {

	once.Do(func() {
		fileWriter := &lumberjack.Logger{
			Filename:   loggerConfig.FilePath,
			LocalTime:  loggerConfig.UserLocalTime,
			MaxSize:    loggerConfig.FileMaxSizeInMB,
			MaxBackups: loggerConfig.FileMaxBackups,
			MaxAge:     loggerConfig.FileMaxAgeInDays,
			Compress:   loggerConfig.FileCompress,
		}

		consoleWriter := &slog.HandlerOptions{Level: loggerConfig.LogLevel}

		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(fileWriter, os.Stdout), consoleWriter))

	})

	return logger
}
