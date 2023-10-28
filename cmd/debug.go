package cmd

import (
	"errors"

	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:    "debug expr|file",
	Short:  "Debug commands",
	Long:   `Read source from stdin and parse it as an expr or a file.`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("subcommand required")
		}
		switch args[0] {
		case "expr":
			return runner.ParseExpr()
		case "file":
			return runner.ParseFile()
		default:
			return errors.New("unknown subcommand")
		}
	},
}
