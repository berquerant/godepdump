package cmd

import (
	"regexp"

	"github.com/berquerant/godepdump/logx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func displayFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		logx.Debug(
			"flags",
			logx.S("name", f.Name),
			logx.B("changed", f.Changed),
			logx.S("value", f.Value.String()),
		)
	})
}

func getPatterns(args []string) []string {
	if len(args) > 0 {
		return args
	}
	return []string{"./..."}
}

func setAnalyzeLimit(cmd *cobra.Command) {
	cmd.Flags().UintP("analyzeLimit", "l", 1,
		"type analyze limit. The higher the number, the more detailed the type will be analyzed.")
}

func getAnalyzeLimit(cmd *cobra.Command) uint {
	v, _ := cmd.Flags().GetUint("analyzeLimit")
	return v
}

func setExported(cmd *cobra.Command) {
	cmd.Flags().Bool("exported", false, "list exported node only if true")
}

func getExported(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("exported")
	return v
}

func setPackageName(cmd *cobra.Command) {
	cmd.Flags().String("package", "", "select packages by regexp")
}

func getPackageName(cmd *cobra.Command) *regexp.Regexp {
	return getRegexp(cmd, "package")
}

func setIdentName(cmd *cobra.Command) {
	cmd.Flags().String("ident", "", "select identifiers by regexp")
}

func getIdentName(cmd *cobra.Command) *regexp.Regexp {
	return getRegexp(cmd, "ident")
}

func getRegexp(cmd *cobra.Command, name string) *regexp.Regexp {
	v, err := cmd.Flags().GetString(name)
	if err != nil {
		return nil
	}
	return regexp.MustCompile(v)
}
