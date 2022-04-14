package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func newVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "display version information",
		Action: func(ctx *cli.Context) error {
			fmt.Println("to do, add version information")

			return nil
		},
	}
}
