package microvm

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "microvm",
		Short: "perform microvm operations",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				zap.S().Debugw("ingoring cobra error",
					"error",
					err.Error())
			}
		},
	}

	createCmd := newCreateCommand()
	cmd.AddCommand(createCmd)

	getCmd := newGetCommand()
	cmd.AddCommand(getCmd)

	deleteCmd := newDeleteCommand()
	cmd.AddCommand(deleteCmd)

	return cmd
}
