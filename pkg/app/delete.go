package app

import (
	"context"
	"fmt"

	flintlockv1 "github.com/liquidmetal-dev/flintlock/api/services/microvm/v1alpha1"
)

func (a *app) Delete(ctx context.Context, input *DeleteInput) error {
	a.logger.Debugw("deleting microvm", "uid", input.UID, "host", input.Host)

	client, err := a.createFlintlockClient(input.Host)
	if err != nil {
		return fmt.Errorf("creating flintlock client for %s: %w", input.Host, err)
	}

	deleteInput := &flintlockv1.DeleteMicroVMRequest{
		Uid: input.UID,
	}

	_, err = client.DeleteMicroVM(ctx, deleteInput)
	if err != nil {
		return fmt.Errorf("deleting microvm with uid %s: %w", input.UID, err)
	}

	a.logger.Infow("deleted microvm", "uid", input.UID, "host", input.Host)

	return nil
}
