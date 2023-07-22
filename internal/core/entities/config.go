package entities

import (
	"dagger.io/dagger"
	"github.com/excoriate/stiletto/internal/tui"
)

type HostConfig struct {
	EnvVars map[string]string
	DirCfg  *DirCfg
}

type DaggerConfig struct {
	DirCfg *DirCfg
	Client *dagger.Client
}

type CLIConfigArgs struct {
	DirCfg *DirCfg
}

type CLIConfig struct {
	// UXLog is the logger for the UX. Used in the context of the CLI client.
	UXLog tui.UXMessenger

	// UXTitles is the title generator for the UX. Used in the context of the CLI client.
	UXTitles tui.UXTitleGenerator

	// DirCfg is the configuration for the directories.
	DirCfg *DirCfg
}

type ClientConfig struct {
	HostEnvVars    map[string]string
	Classification string
	DirCfg         *DirCfg
}

// APIConfig TODO: Pending to implement.
type APIConfig struct {
	// ClientToken is the token used to identify the client.
	ClientToken string
}

type DirCfg struct {
	BaseDir    string
	BaseDirAbs string
	HomeDir    string
	HomeDirAbs string
	IsGitRepo  bool
	GitDirAbs  string
}
