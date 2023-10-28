package logx

import (
	"io"
	"log/slog"
	"sync"
)

type internalLogger interface {
	debug(msg string, v ...Attr)
	info(msg string, v ...Attr)
	warn(msg string, v ...Attr)
	error(msg string, v ...Attr)
}

var (
	setupOnce sync.Once
	instance  internalLogger
	levelVar  *slog.LevelVar
	output    io.Writer = io.Discard
)

func setup(w io.Writer) {
	setupOnce.Do(func() {
		output = w
		setupInstance(output)
	})
}

func setupInstance(w io.Writer) {
	levelVar = new(slog.LevelVar)
	instance = &logger{
		Logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: levelVar,
		})),
	}
}

type logger struct {
	*slog.Logger
}

func (l *logger) debug(msg string, v ...Attr) { l.logAttrs(slog.LevelDebug, msg, v...) }
func (l *logger) info(msg string, v ...Attr)  { l.logAttrs(slog.LevelInfo, msg, v...) }
func (l *logger) warn(msg string, v ...Attr)  { l.logAttrs(slog.LevelWarn, msg, v...) }
func (l *logger) error(msg string, v ...Attr) { l.logAttrs(slog.LevelError, msg, v...) }

func (l *logger) logAttrs(level slog.Level, msg string, v ...Attr) {
	attrs := make([]slog.Attr, len(v))
	for i, attr := range v {
		attrs[i] = slog.Attr(attr)
	}
	l.LogAttrs(nil, level, msg, attrs...)
}

func getInstance() internalLogger {
	setup(output)
	return instance
}
