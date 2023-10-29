package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"

	"github.com/berquerant/godepdump/logx"

	"github.com/spf13/cobra"
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
