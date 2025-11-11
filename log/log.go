package log

import (
	"log/slog"
	"os"
)

// Log is the default logger for any log message
var Log *slog.Logger

// Audit is the degault logger for any audit message
var Audit *slog.Logger

var logLevel *slog.LevelVar
var auditLevel *slog.LevelVar

func init() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	logHadler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	Log = slog.New(logHadler)

	auditLevel = new(slog.LevelVar)
	auditLevel.Set(slog.LevelInfo)

	auditHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: auditLevel,
	})
	Audit = slog.New(auditHandler)
}

// Trasform a string into a log level
func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

// set level for Log
func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}
