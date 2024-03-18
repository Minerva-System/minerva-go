package minervalog

import (
	"fmt"
	"os"
	"runtime"
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
	pc, file, line, ok := runtime.Caller(1)
	var fn string
	if !ok {
		file = "unknown"
		line = 0
		fn = "unknown"
	} else {
		fn = runtime.FuncForPC(pc).Name()
	}
	
    slog.Info(
		fmt.Sprintf(format, args...),
		"file", file,
		"line", line,
		"function", fn,
	)
}

func Warn(format string, args ...any) {
	pc, file, line, ok := runtime.Caller(1)
	var fn string
	if !ok {
		file = "unknown"
		line = 0
		fn = "unknown"
	} else {
		fn = runtime.FuncForPC(pc).Name()
	}
	
    slog.Warn(
		fmt.Sprintf(format, args...),
		"file", file,
		"line", line,
		"function", fn,
	)
}

func Debug(format string, args ...any) {
	pc, file, line, ok := runtime.Caller(1)
	var fn string
	if !ok {
		file = "unknown"
		line = 0
		fn = "unknown"
	} else {
		fn = runtime.FuncForPC(pc).Name()
	}
	
    slog.Debug(
		fmt.Sprintf(format, args...),
		"file", file,
		"line", line,
		"function", fn,
	)
}

func Error(format string, args ...any) {
	pc, file, line, ok := runtime.Caller(1)
	var fn string
	if !ok {
		file = "unknown"
		line = 0
		fn = "unknown"
	} else {
		fn = runtime.FuncForPC(pc).Name()
	}
	
    slog.Error(
		fmt.Sprintf(format, args...),
		"file", file,
		"line", line,
		"function", fn,
	)
}

func Fatal(format string, args ...any) {
	pc, file, line, ok := runtime.Caller(1)
	var fn string
	if !ok {
		file = "unknown"
		line = 0
		fn = "unknown"
	} else {
		fn = runtime.FuncForPC(pc).Name()
	}
	
    slog.Error(
		fmt.Sprintf(format, args...),
		"file", file,
		"line", line,
		"function", fn,
	)
	os.Exit(1)
}
