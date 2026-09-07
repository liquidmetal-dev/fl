package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	flintlockv1 "github.com/liquidmetal-dev/flintlock/api/services/microvm/v1alpha1"
	microvmexecv1 "github.com/liquidmetal-dev/flintlock/api/services/microvmexec/v1alpha1"
	microvmsshproxyv1 "github.com/liquidmetal-dev/flintlock/api/services/microvmsshproxy/v1alpha1"
)

type App interface {
	Create(ctx context.Context, input *CreateInput) error
	Get(ctx context.Context, input *GetInput) error
	Delete(ctx context.Context, input *DeleteInput) error
	Exec(ctx context.Context, input *ExecInput) (int, error)
	SSH(ctx context.Context, input *SSHInput) error
}

func New(logger *zap.SugaredLogger) App {
	return &app{
		logger: logger,
	}
}

type app struct {
	logger *zap.SugaredLogger
}

func (a *app) createFlintlockClient(address string) (flintlockv1.MicroVMClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating grpc connection: %w", err)
	}

	flClient := flintlockv1.NewMicroVMClient(conn)

	return flClient, nil
}

func (a *app) createExecClient(address string) (microvmexecv1.MicroVMExecClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating grpc connection: %w", err)
	}

	return microvmexecv1.NewMicroVMExecClient(conn), nil
}

func (a *app) createSSHProxyClient(address string) (microvmsshproxyv1.MicroVMSSHProxyClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating grpc connection: %w", err)
	}

	return microvmsshproxyv1.NewMicroVMSSHProxyClient(conn), nil
}

// resolveGuestAgentAddress calls ServerInfo on host and returns the address of the
// requested optional guest-agent service (exec or ssh-proxy), or an error if that
// service isn't enabled on the server.
func (a *app) resolveGuestAgentAddress(
	ctx context.Context,
	host string,
	name string,
	service func(*flintlockv1.ServerInfoResponse) *flintlockv1.GuestAgentServiceInfo,
) (string, error) {
	client, err := a.createFlintlockClient(host)
	if err != nil {
		return "", fmt.Errorf("creating flintlock client for %s: %w", host, err)
	}

	info, err := client.ServerInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("getting server info from %s: %w", host, err)
	}

	svc := service(info)
	if svc == nil || !svc.Enabled {
		return "", fmt.Errorf("%s is not enabled on flintlock host %s", name, host)
	}

	return svc.Address, nil
}
