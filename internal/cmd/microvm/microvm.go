package microvm

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "microvm",
		Short: "perform microvm operations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newDeleteCommand())

	return cmd
}
