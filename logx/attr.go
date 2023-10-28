package logx

import (
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/exp/constraints"
)

type Attr slog.Attr

func B(k string, v bool) Attr                              { return Attr(slog.Bool(k, v)) }
func S(k, v string) Attr                                   { return Attr(slog.String(k, v)) }
func SS[T ~[]string](k string, v T) Attr                   { return Attr(slog.Any(k, v)) }
func Any(k string, v any) Attr                             { return Attr(slog.Any(k, v)) }
func I[T constraints.Integer](k string, v T) Attr          { return Attr(slog.Int(k, int(v))) }
func II[U ~[]T, T constraints.Integer](k string, v U) Attr { return Attr(slog.Any(k, v)) }
func Err(err error) Attr                                   { return Attr(slog.Any("error", err)) }
func D(k string, v time.Duration) Attr                     { return Attr(slog.Duration(k, v)) }
func V(k string, v any) Attr                               { return Attr(slog.String(k, fmt.Sprintf("%#v", v))) }
