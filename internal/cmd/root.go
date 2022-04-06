package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/weaveworks-experiments/fl/internal/cmd/microvm"
	"github.com/weaveworks-experiments/fl/pkg/logging"
)

const (
	logLevelFlag = "log-level"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fl",
		Short: "The experimental cli for flintlock",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			logLevelFlag, _ := cmd.Flags().GetString(logLevelFlag)
			err := logging.Configure(logLevelFlag)
			if err != nil {
				return err
			}

			return nil
		},
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				zap.S().Debugw("ingoring cobra error",
					"error",
					err.Error())
			}
		},
	}

	cmd.PersistentFlags().String(logLevelFlag, "debug", "set the level of the logger")

	versionCmd := newVersionCommand()
	cmd.AddCommand(versionCmd)

	microvmCmd := microvm.NewCommand()
	cmd.AddCommand(microvmCmd)

	return cmd
}
