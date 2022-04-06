package microvm

import (
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/weaveworks-experiments/fl/pkg/app"
	"github.com/weaveworks-experiments/fl/pkg/flags"
)

const (
	defaultNamespace   = "default"
	defaultVCPU        = 2
	defaultMemoryMb    = 2048
	defaultKernelImage = "ghcr.io/weaveworks/flintlock-kernel:5.10.77"
	defaultKernelFile  = "boot/vmlinux"
	defaultRootImage   = "ghcr.io/weaveworks/capmvm-kubernetes:1.21.8"
)

func newCreateCommand() *cobra.Command {
	createInput := &app.CreateInput{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new microvm",
		PreRun: func(cmd *cobra.Command, args []string) {
			flags.BindFlags(cmd)
		},
		Run: func(c *cobra.Command, _ []string) {
			a := app.New(zap.S().With("action", "create"))
			if err := a.Create(c.Context(), createInput); err != nil {
				zap.S().Errorw("failed creating microvm", "error", err)
			}
		},
	}

	cmd.Flags().StringVar(&createInput.Host, "host", "", "the flintlock host to create the microvm on")
	cmd.MarkFlagRequired("host")

	defaultName := namesgenerator.GetRandomName(10)
	cmd.Flags().StringVar(&createInput.Name, "name", defaultName, "the name of the microvm, auto-generated if not supplied")
	cmd.Flags().StringVar(&createInput.Namespace, "namespace", defaultNamespace, "the namespace for the microvm")
	cmd.Flags().IntVar(&createInput.VCPU, "vcpu", defaultVCPU, "the number of vcpus")
	cmd.Flags().IntVar(&createInput.MemoryInMb, "memory", defaultMemoryMb, "the memory in mb")
	cmd.Flags().StringVar(&createInput.KernelImage, "kernel-image", defaultKernelImage, "the image to use for the kernel")
	cmd.Flags().BoolVar(&createInput.KernelAddNetConf, "add-netconf", true, "automatically add network configuration to the kernel cmd line")
	cmd.Flags().StringVar(&createInput.KernelFileName, "kernel-filename", defaultKernelFile, "name of the kernel file in the image")
	cmd.Flags().StringVar(&createInput.RootImage, "root-image", defaultRootImage, "the image to use for the root volume")
	cmd.Flags().StringVar(&createInput.InitrdImage, "initrd-image", "", "the image to use for the initial ramdisk")
	cmd.Flags().StringVar(&createInput.InitrdFilename, "initrd-filename", "", "name of the file in the image to use for the initial ramdisk")
	cmd.Flags().StringSliceVar(&createInput.NetworkInterfaces, "network-interface", nil, "specify the network interfaces to attach. In the following format: name:type:[macaddress]:[ipaddress]")
	cmd.Flags().StringSliceVar(&createInput.MetadataFromFile, "metadata-from-file", nil, "specify metadata to be available to your microvm. In the following format key=pathtofile")
	cmd.Flags().StringVar(&createInput.Hostname, "hostname", "", "the hostname of the the microvm")
	cmd.Flags().StringVar(&createInput.SSHKeyFile, "ssh-key-file", "", "an ssh key to use")
	//TODO: additional command line args for kernel
	//TODO: add additional volumes

	return cmd
}
