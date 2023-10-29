package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usesCmd)
	setAnalyzeLimit(usesCmd)
	setExported(usesCmd)
	setPackageName(usesCmd)
	setIdentName(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses [--analyzeLimit LIMIT] [--exported] [--package] [--ident] [packages]",
	Short: "List uses",
	Long:  `Displays a list of identifiers and their corresponding objects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &runner.ListUses{
			Patterns:     getPatterns(args),
			AnalyzeLimit: int(getAnalyzeLimit(cmd)),
			Exported:     getExported(cmd),
			PackageName:  getPackageName(cmd),
			IdentName:    getIdentName(cmd),
		}
		return r.Run(cmd.Context())
	},
}
