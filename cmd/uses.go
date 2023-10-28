package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usesCmd)
	setAnalyzeLimit(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses [--analyzeLimit LIMIT] [packages]",
	Short: "List uses",
	Long:  `Displays a list of identifiers and their corresponding objects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.ListUses(
			cmd.Context(),
			getPatterns(args),
			int(getAnalyzeLimit(cmd)),
		)
	},
}
