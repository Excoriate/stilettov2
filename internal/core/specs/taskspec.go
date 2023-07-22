package specs

type TaskManifestSpec struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   TaskMetadata
	Spec       TaskSpec `yaml:"spec"`
}

type TaskMetadata struct {
	Name string `yaml:"name"`
}

type TaskSpec struct {
	ContainerImage string          `yaml:"containerImage"`
	Workdir        string          `yaml:"workdir"`
	MountDir       string          `yaml:"mountDir"`
	BaseDir        string          `yaml:"baseDir"` // Optional, normally it's resolved or computed.
	CommandsSpec   []*CommandsSpec `yaml:"commandsSpec"`
	EnvVarsSpec    EnvVarsSpec     `yaml:"envVarsSpec"`
}

type EnvVarsSpec struct {
	EnvVars map[string]string `yaml:"envVars"`
	Options struct {
		ScanAWSEnvVars       []string `yaml:"scanAWSEnvVars"`
		ScanTerraformEnvVars bool     `yaml:"scanTerraformEnvVars"`
		ScanCustomEnvVars    []string `yaml:"scanCustomEnvVars"`
		DotFiles             []string `yaml:"dotFiles"`
	}
}

type CommandsSpec struct {
	Binary   string   `yaml:"binary"`
	Commands []string `yaml:"commands"`
}
