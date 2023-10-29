package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(defDeclCmd)
	setAnalyzeLimit(defDeclCmd)
	setExported(defDeclCmd)
	setPackageName(defDeclCmd)
	setIdentName(defDeclCmd)
}

var defDeclCmd = &cobra.Command{
	Use:   "defdecl [--analyzeLimit LIMIT] [--exported] [--package] [--ident] [packages]",
	Short: "List def decls",
	Long:  `List functions, types, variables at the top level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &runner.ListDefDecls{
			Patterns:     getPatterns(args),
			AnalyzeLimit: int(getAnalyzeLimit(cmd)),
			Exported:     getExported(cmd),
			PackageName:  getPackageName(cmd),
			IdentName:    getIdentName(cmd),
		}
		return r.Run(cmd.Context())
	},
}
