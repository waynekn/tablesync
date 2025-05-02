package logging

import (
	"log/slog"
	"os"
)

// InitLogger initializes a default slog logger with JSON output to stdout.
// It sets this logger as the default for the slog package.
func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
