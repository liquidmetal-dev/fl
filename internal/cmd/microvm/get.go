package microvm

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/liquidmetal-dev/fl/pkg/app"
)

const (
	examples = `
# Get all microvms from a host
fl microvm get --host host1:9090

# Get a microvm with a specific id
fl microvm get --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ
`
)

func newGetCommand() *cli.Command {
	getInput := &app.GetInput{}

	cmd := &cli.Command{
		Name:  "get",
		Usage: "get details of a microvm(s) from a host",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() > 0 {
				getInput.UID = ctx.Args().First()
			}

			a := app.New(zap.S().With("action", "get"))
			err := a.Get(ctx.Context, getInput)
			if err != nil {
				return fmt.Errorf("getting microvnm: %w", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "the flintlock host to get the microvm from",
				Destination: &getInput.Host,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Usage:       "the namespace to get the microvms from",
				Destination: &getInput.Namespace,
				Value:       defaultNamespace,
			},
		},
	}

	return cmd
}
