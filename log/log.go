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

// set level for Log
func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}
