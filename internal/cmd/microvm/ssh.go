package microvm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/liquidmetal-dev/fl/pkg/app"
)

const (
	sshExamples = `
# SSH into a microvm
fl microvm ssh --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ

# SSH into a microvm as a specific user
fl microvm ssh --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ -- -l ubuntu
`
)

func newSSHCommand() *cobra.Command {
	sshInput := &app.SSHInput{}

	cmd := &cobra.Command{
		Use:     "ssh <vmid> [-- ssh-args...]",
		Short:   "ssh into a microvm",
		Example: sshExamples,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sshInput.UID = args[0]
			sshInput.SSHArgs = args[1:]

			a := app.New(zap.S().With("action", "ssh"))

			if err := a.SSH(cmd.Context(), sshInput); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					os.Exit(exitErr.ExitCode())
				}

				return fmt.Errorf("running ssh: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&sshInput.Host, "host", "", "the flintlock host to ssh through")

	cmd.MarkFlagRequired("host")

	return cmd
}
