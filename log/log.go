package log

import (
	"log/slog"
	"os"
)

// Log is the default logger for any log message
var Log *slog.Logger

// Audit is the degault logger for any audit message
var Audit *slog.Logger

func init() {
	logHadler := slog.NewJSONHandler(os.Stdout, nil)
	Log = slog.New(logHadler)

	auditHandler := slog.NewJSONHandler(os.Stdout, nil)
	Audit = slog.New(auditHandler)
}
