package microvm

import (
	"fmt"

	"github.com/spf13/cobra"
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

func newGetCommand() *cobra.Command {
	getInput := &app.GetInput{}

	cmd := &cobra.Command{
		Use:   "get [vmid]",
		Short: "get details of a microvm(s) from a host",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				getInput.UID = args[0]
			}

			a := app.New(zap.S().With("action", "get"))
			err := a.Get(cmd.Context(), getInput)
			if err != nil {
				return fmt.Errorf("getting microvnm: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&getInput.Host, "host", "", "the flintlock host to create the microvm on")
	cmd.Flags().StringVar(&getInput.Namespace, "namespace", defaultNamespace, "the namespace for the microvm")

	cmd.MarkFlagRequired("host")

	return cmd
}
