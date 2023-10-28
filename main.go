package main

import (
	"github.com/berquerant/godepdump/cmd"
	"github.com/berquerant/godepdump/osx"
)

func main() {
	if err := cmd.Execute(); err != nil {
		osx.Exit(osx.Efailure)
	}
}
