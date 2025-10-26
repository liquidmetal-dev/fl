package microvm

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/liquidmetal-dev/fl/pkg/app"
)

const (
	deleteExamples = `
# Delete a microvm
fl microvm delete --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ
`
)

func newDeleteCommand() *cobra.Command {
	deleteInput := &app.DeleteInput{}

	cmd := &cobra.Command{
		Use:   "delete [vmid]",
		Short: "delete a microvm from a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must supply the uid as an argument")
			}
			deleteInput.UID = args[0]

			a := app.New(zap.S().With("action", "delete"))
			err := a.Delete(cmd.Context(), deleteInput)
			if err != nil {
				return fmt.Errorf("deleting microvm: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&deleteInput.Host, "host", "", "the flintlock host to delete the microvm on")

	cmd.MarkFlagRequired("host")

	return cmd
}
