package minervalog

import (
	"fmt"
	"os"
	"log/slog"
)

func Init() {
	level := slog.LevelInfo
	
	switch os.Getenv("MINERVA_LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
	Info("Log level set to %s", level)
}



func Info(format string, args ...any) {
    slog.Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
    slog.Warn(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...any) {
    slog.Debug(fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
    slog.Error(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...any) {
    slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
