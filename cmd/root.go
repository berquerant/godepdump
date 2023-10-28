package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"

	"github.com/berquerant/godepdump/logx"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rootCmd = &cobra.Command{
		Use:   "godepdump",
		Short: "Dump dependencies",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var (
				debug, _ = cmd.Flags().GetBool("debug")
				quiet, _ = cmd.Flags().GetBool("quiet")
			)
			if quiet {
				cmd.SetOut(io.Discard)
				logx.Setup(io.Discard, debug)
			} else {
				logx.Setup(os.Stderr, debug)
			}
			displayFlags(cmd.Flags())
		},
	}
)

func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logx.Error("execute", logx.Err(err))
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logs")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet logs")
}

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
