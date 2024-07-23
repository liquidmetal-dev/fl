package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/liquidmetal-dev/fl/internal/cmd/microvm"
	"github.com/liquidmetal-dev/fl/pkg/logging"
)

const (
	logLevelFlag = "log-level"
)

func NewApp() *cli.App {
	app := &cli.App{
		Name:     "fl",
		Usage:    "The experimental cli for flintlock",
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  logLevelFlag,
				Usage: "set the level of the debugger",
				Value: "debug",
			},
		},
		Before: func(ctx *cli.Context) error {
			logLevelFlag := ctx.String(logLevelFlag)
			err := logging.Configure(logLevelFlag)
			if err != nil {
				return err
			}

			return nil
		},
	}

	versionCmd := newVersionCommand()
	app.Commands = append(app.Commands, versionCmd)

	microvmCmd := microvm.NewCommand()
	app.Commands = append(app.Commands, microvmCmd)

	return app
}
