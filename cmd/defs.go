package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(defsCmd)
	setAnalyzeLimit(defsCmd)
	setExported(defsCmd)
	setPackageName(defsCmd)
	setIdentName(defsCmd)
}

var defsCmd = &cobra.Command{
	Use:   "defs [--analyzeLimit LIMIT] [--exported] [--package] [--ident] [packages]",
	Short: "List defs",
	Long:  `Display a list of identifiers and their corresponding definitions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &runner.ListDefs{
			Patterns:     getPatterns(args),
			AnalyzeLimit: int(getAnalyzeLimit(cmd)),
			Exported:     getExported(cmd),
			PackageName:  getPackageName(cmd),
			IdentName:    getIdentName(cmd),
		}
		return r.Run(cmd.Context())
	},
}
