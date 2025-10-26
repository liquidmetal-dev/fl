package cmd

import (
	"github.com/spf13/cobra"

	"github.com/liquidmetal-dev/fl/internal/cmd/microvm"
	"github.com/liquidmetal-dev/fl/internal/cmd/version"
	"github.com/liquidmetal-dev/fl/pkg/logging"
)

const (
	logLevelFlag = "log-level"
)

func NewRootCmd() *cobra.Command {
	var logLevel string

	cmd := &cobra.Command{
		Use:   "fl",
		Short: "The experimental cli for flintlock",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := logging.Configure(logLevel)
			if err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.Flags().StringVar(&logLevel, logLevelFlag, "debug", "set the level of the debugger")

	versionCmd := version.NewVersionCommand()
	cmd.AddCommand(versionCmd)

	microvmCmd := microvm.NewCommand()
	cmd.AddCommand(microvmCmd)

	return cmd
}
