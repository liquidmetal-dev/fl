package app

import (
	"context"
	"fmt"
	"os"

	flintlockv1 "github.com/weaveworks/flintlock/api/services/microvm/v1alpha1"
	"gopkg.in/yaml.v2"
)

func (a *app) Get(ctx context.Context, input *GetInput) error {
	if input.UID == "" {
		a.logger.Debugw("getting all microvms", "host", input.Host)

		return a.listMicrovms(ctx, input)
	} else {
		a.logger.Debugw("getting microvm", "uid", input.UID, "host", input.Host)

		return a.getMicrovm(ctx, input)
	}
}

func (a *app) listMicrovms(ctx context.Context, input *GetInput) error {
	client, err := a.createFlintlockClient(input.Host)
	if err != nil {
		return fmt.Errorf("creating flintlock client for %s: %w", input.Host, err)
	}

	listInput := &flintlockv1.ListMicroVMsRequest{
		Namespace: input.Namespace,
	}

	listOutput, err := client.ListMicroVMs(ctx, listInput)
	if err != nil {
		return fmt.Errorf("listing microvms: %w", err)
	}

	if len(listOutput.Microvm) == 0 {
		fmt.Fprintln(os.Stderr, "No microvms found.")
	}

	for i := range listOutput.Microvm {
		microvm := listOutput.Microvm[i]

		fmt.Fprintf(os.Stderr, "%s (%s/%s): %s\n", *microvm.Spec.Uid, microvm.Spec.Namespace, microvm.Spec.Id, microvm.Status.State.String())
	}

	return nil
}

func (a *app) getMicrovm(ctx context.Context, input *GetInput) error {
	client, err := a.createFlintlockClient(input.Host)
	if err != nil {
		return fmt.Errorf("creating flintlock client for %s: %w", input.Host, err)
	}

	getInput := &flintlockv1.GetMicroVMRequest{
		Uid: input.UID,
	}

	getOutput, err := client.GetMicroVM(ctx, getInput)
	if err != nil {
		return fmt.Errorf("getting microvm with uid %s: %w", input.UID, err)
	}

	if getOutput.Microvm == nil {
		a.logger.Infow("microvm not found", "uid", input.UID, "host", input.Host)

		return nil
	}

	data, err := yaml.Marshal(getOutput.Microvm)
	if err != nil {
		return fmt.Errorf("marshalling microvm to yaml: %w", err)
	}

	fmt.Fprintln(os.Stderr, string(data))

	return nil
}
