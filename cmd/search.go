package cmd

import (
	"errors"
	"regexp"

	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	setAnalyzeLimit(searchCmd)
	setExported(searchCmd)
	setPackageName(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search ident|decl NAME [--analyzeLimit LIMIT] [--exported] [--package] [patterns]",
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
			name     = regexp.MustCompile(args[1])
			patterns = getPatterns(args[2:])
			limit    = int(getAnalyzeLimit(cmd))
			exported = getExported(cmd)
			pkgName  = getPackageName(cmd)
		)
		switch command {
		case "ident":
			r := &runner.SearchIdent{
				Patterns:     patterns,
				Name:         name,
				AnalyzeLimit: limit,
				Exported:     exported,
				PackageName:  pkgName,
			}
			return r.Run(cmd.Context())
		case "decl":
			r := &runner.SearchDecl{
				Patterns:     patterns,
				Name:         name,
				AnalyzeLimit: limit,
				Exported:     exported,
				PackageName:  pkgName,
			}
			return r.Run(cmd.Context())
		default:
			return errors.New("unknown subcommand")
		}
	},
}
