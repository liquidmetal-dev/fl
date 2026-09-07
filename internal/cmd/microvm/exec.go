package microvm

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/liquidmetal-dev/fl/pkg/app"
)

const (
	execExamples = `
# Run a command in a microvm
fl microvm exec --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ -- echo hello

# Run a command with stdin forwarded
fl microvm exec --host host1:9090 01FZZJV1XD2FKH2KY0NDB4MBRQ --stdin -- cat
`
)

func newExecCommand() *cobra.Command {
	execInput := &app.ExecInput{}
	envVars := []string{}

	cmd := &cobra.Command{
		Use:     "exec [vmid] -- <cmd> [args...]",
		Short:   "execute a command in a microvm",
		Example: execExamples,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.ArgsLenAtDash() != 1 {
				return errors.New("you must separate the vmid and command with --, e.g. exec <vmid> -- <cmd> [args...]")
			}

			execInput.UID = args[0]

			cmdArgs := args[1:]
			if len(cmdArgs) == 0 {
				return errors.New("you must supply the command to run after --")
			}
			execInput.Cmd = cmdArgs[0]
			execInput.Args = cmdArgs[1:]

			env, err := parseEnvVars(envVars)
			if err != nil {
				return fmt.Errorf("parsing env vars: %w", err)
			}
			execInput.Env = env

			a := app.New(zap.S().With("action", "exec"))

			exitCode, err := a.Exec(cmd.Context(), execInput)
			if err != nil {
				return fmt.Errorf("executing command: %w", err)
			}

			if exitCode != 0 {
				os.Exit(exitCode)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&execInput.Host, "host", "", "the flintlock host to run the command on")
	cmd.Flags().StringVar(&execInput.Cwd, "cwd", "", "the working directory to run the command in")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "environment variables to set, in the format KEY=VALUE")
	cmd.Flags().BoolVar(&execInput.Shell, "shell", false, "run the command via a shell")
	cmd.Flags().StringVar(&execInput.User, "user", "", "the guest user to run the command as")
	cmd.Flags().Int32Var(&execInput.TimeoutSeconds, "timeout", 0, "timeout in seconds for the command, 0 means unbounded")
	cmd.Flags().BoolVarP(&execInput.Stdin, "stdin", "i", false, "forward local stdin to the command")

	cmd.MarkFlagRequired("host")

	return cmd
}

func parseEnvVars(vars []string) (map[string]string, error) {
	env := map[string]string{}

	for _, v := range vars {
		key, value, ok := strings.Cut(v, "=")
		if !ok {
			return nil, fmt.Errorf("invalid env var %q, expected KEY=VALUE", v)
		}

		env[key] = value
	}

	return env, nil
}
