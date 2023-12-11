package logger

import (
	"github.com/pterm/pterm"
	"log/slog"
)

func New(debug bool) *slog.Logger {
	level := pterm.LogLevelInfo

	if debug {
		level = pterm.LogLevelDebug
	}

	return slog.New(pterm.NewSlogHandler(pterm.DefaultLogger.WithLevel(level)))
}
