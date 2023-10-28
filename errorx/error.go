package errorx

import (
	"errors"
	"fmt"
)

func Errorf(err error, format string, v ...any) error {
	return fmt.Errorf("%w %s", err, fmt.Sprintf(format, v...))
}

func New(msg string) error { return errors.New(msg) }
