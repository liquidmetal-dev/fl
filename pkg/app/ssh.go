package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"

	flintlockv1 "github.com/liquidmetal-dev/flintlock/api/services/microvm/v1alpha1"
	microvmsshproxyv1 "github.com/liquidmetal-dev/flintlock/api/services/microvmsshproxy/v1alpha1"
)

func (a *app) SSH(ctx context.Context, input *SSHInput) error {
	a.logger.Debugw("ssh to microvm", "uid", input.UID, "host", input.Host)

	address, err := a.resolveGuestAgentAddress(ctx, input.Host, "ssh", func(r *flintlockv1.ServerInfoResponse) *flintlockv1.GuestAgentServiceInfo {
		return r.SshProxy
	})
	if err != nil {
		return err
	}

	client, err := a.createSSHProxyClient(address)
	if err != nil {
		return fmt.Errorf("creating ssh-proxy client for %s: %w", address, err)
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := client.SSHProxy(sessionCtx)
	if err != nil {
		return fmt.Errorf("opening ssh-proxy stream: %w", err)
	}

	startReq := &microvmsshproxyv1.SSHProxyRequest{
		Payload: &microvmsshproxyv1.SSHProxyRequest_Uid{Uid: input.UID},
	}

	if err := stream.Send(startReq); err != nil {
		return fmt.Errorf("sending ssh-proxy start: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("starting local ssh proxy listener: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	acceptErrCh := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptErrCh <- fmt.Errorf("accepting local ssh connection: %w", err)

			return
		}
		defer conn.Close()

		bridgeSSHProxy(conn, stream, cancel)
		acceptErrCh <- nil
	}()

	sshArgs := []string{
		"-p", strconv.Itoa(port),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	sshArgs = append(sshArgs, input.SSHArgs...)
	sshArgs = append(sshArgs, "127.0.0.1")

	sshCmd := exec.CommandContext(ctx, "ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	runErr := sshCmd.Run()

	listener.Close()
	cancel()

	if bridgeErr := <-acceptErrCh; bridgeErr != nil {
		a.logger.Debugw("ssh proxy bridge ended", "uid", input.UID, "host", input.Host, "error", bridgeErr)
	}

	if runErr != nil {
		return fmt.Errorf("running ssh: %w", runErr)
	}

	return nil
}

func bridgeSSHProxy(conn net.Conn, stream microvmsshproxyv1.MicroVMSSHProxy_SSHProxyClient, cancel context.CancelFunc) {
	// stop unblocks whichever side is still copying once the other side has finished:
	// closing conn unblocks a pending conn.Read/Write, cancelling the session unblocks a
	// pending stream.Recv/Send. Both are safe to call more than once.
	stop := func() {
		conn.Close()
		cancel()
	}

	done := make(chan struct{}, 2)

	go func() {
		defer func() { done <- struct{}{} }()
		defer stop()

		buf := make([]byte, 32*1024)

		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])

				if sendErr := stream.Send(&microvmsshproxyv1.SSHProxyRequest{
					Payload: &microvmsshproxyv1.SSHProxyRequest_Data{Data: chunk},
				}); sendErr != nil {
					return
				}
			}

			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer func() { done <- struct{}{} }()
		defer stop()

		for {
			resp, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					return
				}

				return
			}

			if _, err := conn.Write(resp.Data); err != nil {
				return
			}
		}
	}()

	<-done
	<-done
}
