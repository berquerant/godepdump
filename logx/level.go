package logx

import "log/slog"

type Level int

const (
	Linfo Level = iota
	Ldebug
	Lwarn
)

func (l Level) intoSlog() slog.Level {
	switch l {
	case Linfo:
		return slog.LevelInfo
	case Ldebug:
		return slog.LevelDebug
	case Lwarn:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
