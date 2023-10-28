package cmd

import (
	"errors"

	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	setAnalyzeLimit(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search ident|decl NAME [--analyzeLimit LIMIT] [patterns]",
	Short: "Search identifiers",
	Long:  `Search identifiers by name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			return errors.New("subcommand required")
		case 1:
			return errors.New("NAME required")
		}

		var (
			command  = args[0]
			name     = args[1]
			patterns = getPatterns(args[2:])
			limit    = int(getAnalyzeLimit(cmd))
		)
		switch command {
		case "ident":
			return runner.SearchIdent(
				cmd.Context(),
				patterns,
				name,
				limit,
			)
		case "decl":
			return runner.SearchDecl(
				cmd.Context(),
				patterns,
				name,
				limit,
			)
		default:
			return errors.New("unknown subcommand")
		}
	},
}
