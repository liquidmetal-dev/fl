package app

type CreateInput struct {
	Host              string
	Name              string
	NameAutogenerate  bool
	Namespace         string
	VCPU              int
	MemoryInMb        int
	KernelImage       string
	KernelAddNetConf  bool
	KernelFileName    string
	RootImage         string
	InitrdImage       string
	InitrdFilename    string
	NetworkInterfaces []string
	MetadataFromFile  []string
	Metadata          Metadata
	MetadataAddVolume bool
}

type Metadata struct {
	Hostname   string
	SSHKeyFile string
	ResolvdFix bool
	Message    string
}

func (m Metadata) IsEmpty() bool {
	if m.Hostname != "" {
		return false
	}
	if m.SSHKeyFile != "" {
		return false
	}
	if m.Message != "" {
		return false
	}
	if m.ResolvdFix {
		return false
	}

	return true
}

type GetInput struct {
	Host      string
	Namespace string
	UID       string
}

type DeleteInput struct {
	Host string
	UID  string
}
