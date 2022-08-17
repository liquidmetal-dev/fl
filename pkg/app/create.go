package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/yitsushi/macpot"
	"gopkg.in/yaml.v2"

	flintlockv1 "github.com/weaveworks-liquidmetal/flintlock/api/services/microvm/v1alpha1"
	flintlocktypes "github.com/weaveworks-liquidmetal/flintlock/api/types"
	"github.com/weaveworks-liquidmetal/flintlock/client/cloudinit/userdata"
)

func (a *app) Create(ctx context.Context, input *CreateInput) error {
	a.logger.Debug("creating a microvm")

	spec, err := a.convertCreateInputToReq(input)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if input.NameAutogenerate {
		spec.Id = namesgenerator.GetRandomName(10)
	}

	if !input.Metadata.IsEmpty() {
		if metaErr := a.addUserdata(spec, input); metaErr != nil {
			return fmt.Errorf("adding user-data: %w", metaErr)
		}
	}

	client, err := a.createFlintlockClient(input.Host)
	if err != nil {
		return fmt.Errorf("creating flintlock client for %s: %w", input.Host, err)
	}

	createInput := &flintlockv1.CreateMicroVMRequest{
		Microvm: spec,
	}

	resp, err := client.CreateMicroVM(ctx, createInput)
	if err != nil {
		return fmt.Errorf("creating microvm: %s", err)
	}

	a.logger.Infow("created microvm", "uid", *resp.Microvm.Spec.Uid, "name", input.Name, "namespace", input.Namespace, "host", input.Host)

	return nil
}

func (a *app) addUserdata(spec *flintlocktypes.MicroVMSpec, input *CreateInput) error {
	userData, err := a.createUserData(input.Metadata)
	if err != nil {
		return fmt.Errorf("creating user-data for microvm: %w", err)
	}
	spec.Metadata.Items["user-data"] = userData

	return nil
}

//TODO: this whole thing needs rewriting
func (a *app) convertCreateInputToReq(input *CreateInput) (*flintlocktypes.MicroVMSpec, error) {
	req := &flintlocktypes.MicroVMSpec{
		Id:        input.Name,
		Namespace: input.Namespace,
		Labels: map[string]string{
			"created-with": "fl",
		},
		Vcpu:       int32(input.VCPU),
		MemoryInMb: int32(input.MemoryInMb),
		Kernel: &flintlocktypes.Kernel{
			Image:            input.KernelImage,
			AddNetworkConfig: input.KernelAddNetConf,
			Filename:         &input.KernelFileName,
			//TODO: additional args
		},
		RootVolume: &flintlocktypes.Volume{
			Id:         "root",
			IsReadOnly: false,
			Source: &flintlocktypes.VolumeSource{
				ContainerSource: &input.RootImage,
			},
		},
		Metadata: &flintlocktypes.Metadata{
			Items:     make(map[string]string),
			AddVolume: &input.MetadataAddVolume,
		},
		AdditionalVolumes: []*flintlocktypes.Volume{},
		Interfaces:        []*flintlocktypes.NetworkInterface{},
	}

	if input.InitrdImage != "" {
		req.Initrd = &flintlocktypes.Initrd{
			Image:    input.InitrdImage,
			Filename: &input.InitrdFilename,
		}
	}

	for i := range input.MetadataFromFile {
		metaFromFile := input.MetadataFromFile[i]
		metaparts := strings.Split(metaFromFile, "=")
		if len(metaparts) != 2 {
			//TODO: proper error
			return nil, fmt.Errorf("metadata not in name=pathtofile format")
		}
		content, err := os.ReadFile(metaparts[1])
		if err != nil {
			return nil, fmt.Errorf("reading metadata from file %s: %w", metaparts[1], err)
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		req.Metadata.Items[metaparts[0]] = encoded
	}

	for i := range input.NetworkInterfaces {
		netInt := input.NetworkInterfaces[i]
		netParts := strings.Split(netInt, ":")
		if len(netParts) < 1 || len(netParts) > 4 {
			//TODO: proper error
			return nil, fmt.Errorf("network interfaces not in correct format, expect name:type:[macaddress]:[ipaddress]")
		}

		macAddress := ""
		ipAddress := ""
		name := netParts[0]
		intType := netParts[1] //TODO: validate the types

		if name == "eth0" {
			return nil, fmt.Errorf("you cannot use eth0 as the name of the interface as this is reserved")
		}

		if len(netParts) == 3 {
			macAddress = netParts[2]
		}
		if len(netParts) == 4 {
			ipAddress = netParts[3]
		}
		if macAddress == "" {
			mac, err := macpot.New(macpot.AsLocal(), macpot.AsUnicast())
			if err != nil {
				return nil, fmt.Errorf("creating mac address client: %w", err)
			}
			macAddress = mac.ToString()
		}

		apiIface := &flintlocktypes.NetworkInterface{
			DeviceId: name,
			GuestMac: &macAddress,
		}

		if ipAddress != "" {
			apiIface.Address = &flintlocktypes.StaticAddress{
				Address: ipAddress,
			}
		}

		switch intType {
		case "macvtap":
			apiIface.Type = flintlocktypes.NetworkInterface_MACVTAP
		case "tap":
			apiIface.Type = flintlocktypes.NetworkInterface_TAP
		}

		req.Interfaces = append(req.Interfaces, apiIface)
	}

	return req, nil
}

func (a *app) createUserData(metadata Metadata) (string, error) {
	userMetadata := &userdata.UserData{
		FinalMessage: "The fl booted system is good to go after $UPTIME seconds",
	}

	if metadata.Message != "" {
		userMetadata.FinalMessage = metadata.Message
	}

	if metadata.ResolvdFix {
		userMetadata.BootCommands = []string{
			"ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf",
		}
	}
	//userMetadata.BootCommands = append(userMetadata.BootCommands, "ip route add 169.254.169.254 dev eth0")

	if metadata.Hostname != "" {
		userMetadata.HostName = metadata.Hostname
	}

	if metadata.SSHKeyFile != "" {
		data, err := os.ReadFile(metadata.SSHKeyFile)
		if err != nil {
			return "", fmt.Errorf("reading ssh key file %s: %w", metadata.SSHKeyFile, err)
		}
		sshKey := string(data)

		defaultUser := userdata.User{
			Name: "ubuntu",
		}
		rootUser := userdata.User{
			Name: "root",
		}

		defaultUser.SSHAuthorizedKeys = []string{
			sshKey,
		}
		rootUser.SSHAuthorizedKeys = []string{
			sshKey,
		}

		userMetadata.Users = []userdata.User{defaultUser, rootUser}
	}

	data, err := yaml.Marshal(userMetadata)
	if err != nil {
		return "", fmt.Errorf("marshalling bootstrap data: %w", err)
	}
	dataWithHeader := append([]byte("## template: jinja\n#cloud-config\n\n"), data...)

	return base64.StdEncoding.EncodeToString(dataWithHeader), nil
}
