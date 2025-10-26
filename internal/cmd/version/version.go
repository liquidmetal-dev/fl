package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display version information",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("to do, add version information")

			return nil
		},
	}
}
