package app

import (
	"context"
	"fmt"
	"io"
	"os"

	flintlockv1 "github.com/liquidmetal-dev/flintlock/api/services/microvm/v1alpha1"
	microvmexecv1 "github.com/liquidmetal-dev/flintlock/api/services/microvmexec/v1alpha1"
)

const execStdinChunkSize = 32 * 1024

func (a *app) Exec(ctx context.Context, input *ExecInput) (int, error) {
	a.logger.Debugw("exec on microvm", "uid", input.UID, "host", input.Host)

	address, err := a.resolveGuestAgentAddress(ctx, input.Host, "exec", func(r *flintlockv1.ServerInfoResponse) *flintlockv1.GuestAgentServiceInfo {
		return r.Exec
	})
	if err != nil {
		return 0, err
	}

	client, err := a.createExecClient(address)
	if err != nil {
		return 0, fmt.Errorf("creating exec client for %s: %w", address, err)
	}

	stream, err := client.ExecCommand(ctx)
	if err != nil {
		return 0, fmt.Errorf("opening exec stream: %w", err)
	}

	startReq := &microvmexecv1.ExecCommandRequest{
		Payload: &microvmexecv1.ExecCommandRequest_Start{
			Start: &microvmexecv1.ExecStart{
				Uid:            input.UID,
				Cmd:            input.Cmd,
				Args:           input.Args,
				Cwd:            input.Cwd,
				Env:            input.Env,
				Shell:          input.Shell,
				User:           input.User,
				TimeoutSeconds: input.TimeoutSeconds,
				HasStdin:       input.Stdin,
			},
		},
	}

	if err := stream.Send(startReq); err != nil {
		return 0, fmt.Errorf("sending exec start: %w", err)
	}

	if input.Stdin {
		go streamStdin(stream)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return 0, nil
			}

			return 0, fmt.Errorf("receiving exec response: %w", err)
		}

		switch payload := resp.Payload.(type) {
		case *microvmexecv1.ExecCommandResponse_Stdout:
			os.Stdout.Write(payload.Stdout)
		case *microvmexecv1.ExecCommandResponse_Stderr:
			os.Stderr.Write(payload.Stderr)
		case *microvmexecv1.ExecCommandResponse_Error:
			a.logger.Errorw("exec error", "uid", input.UID, "host", input.Host, "error", payload.Error)
		case *microvmexecv1.ExecCommandResponse_ExitCode:
			return int(payload.ExitCode), nil
		}
	}
}

func streamStdin(stream microvmexecv1.MicroVMExec_ExecCommandClient) {
	buf := make([]byte, execStdinChunkSize)

	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])

			if sendErr := stream.Send(&microvmexecv1.ExecCommandRequest{
				Payload: &microvmexecv1.ExecCommandRequest_Stdin{Stdin: chunk},
			}); sendErr != nil {
				return
			}
		}

		if err != nil {
			stream.Send(&microvmexecv1.ExecCommandRequest{
				Payload: &microvmexecv1.ExecCommandRequest_StdinEof{StdinEof: true},
			})

			return
		}
	}
}
