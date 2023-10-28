package osx

import "os"

type ExitCode int

const (
	Esuccess ExitCode = iota
	Efailure
)

func Exit(code ExitCode) { os.Exit(int(code)) }
