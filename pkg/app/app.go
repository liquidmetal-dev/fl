package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	flintlockv1 "github.com/weaveworks/flintlock/api/services/microvm/v1alpha1"
)

type App interface {
	Create(ctx context.Context, input *CreateInput) error
	Get(ctx context.Context, input *GetInput) error
	Delete(ctx context.Context, input *DeleteInput) error
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
