package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(defsCmd)
	setAnalyzeLimit(defsCmd)
}

var defsCmd = &cobra.Command{
	Use:   "defs [--analyzeLimit LIMIT] [packages]",
	Short: "List defs",
	Long:  `Display a list of identifiers and their corresponding definitions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.ListDefs(
			cmd.Context(),
			getPatterns(args),
			int(getAnalyzeLimit(cmd)),
		)
	},
}
