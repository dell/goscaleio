package log

import (
	"log/slog"
	"os"
)

var (
	logLevel = new(slog.LevelVar) // Info by default
	Log      = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	Debug    = false // False by default, will turn true if level is set to debug
)

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
	if level == slog.LevelDebug {
		Debug = true
	}
}

func DoLog(
	l func(msg string, args ...any),
	msg string,
) {
	if Debug {
		l(msg)
	}
}
