package cmd

import (
	"github.com/berquerant/godepdump/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [packages]",
	Short: "List imports",
	Long: `The list of imported packages is output in the following format.

{
  "src": {
    "name": "PACKAGE_NAME of import source package",
    "path": "PACKAGE_PATH of import source package"
  },
  "dst": {
    "name": "PACKAGE_NAME of import destination package",
    "path": "PACKAGE_PATH of import destination package"
  }
}

e.g.
package example
import "os"

result:
{
  "src": {"name": "example", "path": "path/to/example"},
  "dst": {"name": "os", "path": "os"}
}`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &runner.ListImport{
			Patterns: getPatterns(args),
		}
		return r.Run(cmd.Context())
	},
}
