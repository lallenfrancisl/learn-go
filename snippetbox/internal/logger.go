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

func New() *Logger {
	return &Logger{slog.New(slog.NewTextHandler(os.Stdout, nil))}
}
