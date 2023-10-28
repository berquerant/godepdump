package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(depCmd)
	setAnalyzeLimit(depCmd)
}

var depCmd = &cobra.Command{
	Use:   "dep [--anlyzeLimit LIMIT] [packages]",
	Short: "List dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.ListDeps(
			cmd.Context(),
			getPatterns(args),
			int(getAnalyzeLimit(cmd)),
		)
	},
}
