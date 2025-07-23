package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	once sync.Once = sync.Once{}
	l    *logger   = &logger{}
)

type logger struct {
	log *slog.Logger
}

func New(lvl string) *logger {

	once.Do(func() {
		handler := slog.NewTextHandler(os.Stderr, nil)

		switch strings.ToUpper(lvl) {
		case "DEBUG":
			handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
		case "WARN":
			handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})
		case "ERROR":
			handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})
		}
		l.log = slog.New(handler)
	})

	return l
}

func (l logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l logger) Warn(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	l.log.Info(msg, args...)
}
