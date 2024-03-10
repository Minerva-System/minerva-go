package minervalog

import (
	"fmt"
	"os"
	"log/slog"
)

func Init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
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
