package microvm

import (
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "microvm",
		Usage:       "perform microvm operations",
		Subcommands: []*cli.Command{},
	}

	subCommands := []*cli.Command{
		newCreateCommand(),
		newGetCommand(),
		newDeleteCommand(),
	}
	cmd.Subcommands = append(cmd.Subcommands, subCommands...)

	return cmd
}
