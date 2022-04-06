package microvm

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/weaveworks-experiments/fl/pkg/app"
	"github.com/weaveworks-experiments/fl/pkg/flags"
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
		Use:     "delete",
		Short:   "delete a microvm from a host",
		Example: deleteExamples,
		Args:    cobra.ExactValidArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			flags.BindFlags(cmd)
		},
		Run: func(c *cobra.Command, args []string) {
			deleteInput.UID = args[0]

			a := app.New(zap.S().With("action", "delete"))
			err := a.Delete(c.Context(), deleteInput)
			if err != nil {
				zap.S().Errorw("failed deleting microvm", "error", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&deleteInput.Host, "host", "", "the flintlock host to get the microvms from")
	cmd.MarkFlagRequired("host")

	return cmd
}
