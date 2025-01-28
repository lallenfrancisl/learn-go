package log

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func (s *Logger) Fatal(msg string, args ...any) {
	s.Error(msg, args...)
}

// Return a logger instance
//
// handler can be passed nil to instantiate a new handler automatically
func New(handler slog.Handler) *Logger {
	if handler != nil {
		return &Logger{slog.New(handler)}
	}

	return &Logger{slog.New(slog.NewTextHandler(os.Stdout, nil))}
}
