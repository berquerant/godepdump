package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(defDeclCmd)
	setAnalyzeLimit(defDeclCmd)
}

var defDeclCmd = &cobra.Command{
	Use:   "defdecl [--analyzeLimit LIMIT] [packages]",
	Short: "List def decls",
	Long:  `List functions, types, variables at the top level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.ListDefDecls(
			cmd.Context(),
			getPatterns(args),
			int(getAnalyzeLimit(cmd)),
		)
	},
}
