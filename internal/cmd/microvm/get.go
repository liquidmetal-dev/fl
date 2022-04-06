package microvm

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/weaveworks-experiments/fl/pkg/app"
	"github.com/weaveworks-experiments/fl/pkg/flags"
)

const (
	examples = `
# Get all microvms from a host
fl microvm get --host host1:9090

# Get a microvm with a specific id
fl microvm get --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ
`
)

func newGetCommand() *cobra.Command {
	getInput := &app.GetInput{}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get details of a microvm(s) from a host",
		Example: examples,
		Args:    cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			flags.BindFlags(cmd)
		},
		Run: func(c *cobra.Command, args []string) {
			if len(args) > 0 {
				getInput.UID = args[0]
			}

			a := app.New(zap.S().With("action", "get"))
			err := a.Get(c.Context(), getInput)
			if err != nil {
				zap.S().Errorw("failed getting microvm(s)", "error", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&getInput.Host, "host", "", "the flintlock host to get the microvms from")
	cmd.Flags().StringVar(&getInput.Namespace, "namespace", defaultNamespace, "the namespace to get the microvms from")

	cmd.MarkFlagRequired("host")

	return cmd
}
