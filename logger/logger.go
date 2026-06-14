package logger

import (
	"log/slog"
	"os"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

type Options struct {
	Level  slog.Level
	Format Format
}

// New creates a new slog.Logger with the given options.
func New(opts Options) *slog.Logger {
	ho := &slog.HandlerOptions{Level: opts.Level}

	var h slog.Handler
	if opts.Format == FormatJSON {
		h = slog.NewJSONHandler(os.Stdout, ho)
	} else {
		h = slog.NewTextHandler(os.Stdout, ho)
	}

	return slog.New(h)
}

// Development returns a debug-level text logger suitable for local development.
func Development() *slog.Logger {
	return New(Options{Level: slog.LevelDebug, Format: FormatText})
}

// Production returns an info-level JSON logger suitable for production.
func Production() *slog.Logger {
	return New(Options{Level: slog.LevelInfo, Format: FormatJSON})
}
