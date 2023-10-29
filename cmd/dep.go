package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(depCmd)
	setAnalyzeLimit(depCmd)
	setExported(depCmd)
	setPackageName(depCmd)
	setIdentName(depCmd)
}

var depCmd = &cobra.Command{
	Use:   "dep [--anlyzeLimit LIMIT] [--exported] [--package] [--ident] [packages]",
	Short: "List dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &runner.ListDeps{
			Patterns:     getPatterns(args),
			AnalyzeLimit: int(getAnalyzeLimit(cmd)),
			Exported:     getExported(cmd),
			PackageName:  getPackageName(cmd),
			IdentName:    getIdentName(cmd),
		}
		return r.Run(cmd.Context())
	},
}
