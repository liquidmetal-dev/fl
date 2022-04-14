package microvm

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/weaveworks-experiments/fl/pkg/app"
)

const (
	deleteExamples = `
# Delete a microvm
fl microvm delete --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ
`
)

func newDeleteCommand() *cli.Command {
	deleteInput := &app.DeleteInput{}

	cmd := &cli.Command{
		Name:  "delete",
		Usage: "delete a microvm from a host",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return errors.New("you must supply the uid as an argument")
			}
			deleteInput.UID = ctx.Args().Get(0)

			a := app.New(zap.S().With("action", "delete"))
			err := a.Delete(ctx.Context, deleteInput)
			if err != nil {
				return fmt.Errorf("deleting microvm: %w", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "the flintlock host to delete the microvm on",
				Destination: &deleteInput.Host,
				Required:    true,
			},
		},
	}

	return cmd
}
