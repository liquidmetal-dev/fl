package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display version information",
		Run: func(c *cobra.Command, _ []string) {
			fmt.Println("to do, add version information")
		},
	}

}
