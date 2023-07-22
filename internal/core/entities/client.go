package entities

import (
	"context"
	"go.uber.org/zap"
)

type Client struct {
	// Id is the unique identifier for the instance.
	Id string

	// Ctx is the context for the instance.
	Ctx *context.Context

	// Logger, mostly used for the API server.
	Logger *zap.Logger

	// DaggerClient is the dagger client.
	CfgDagger *DaggerConfig

	// Config is the configuration for the instance.
	CfgCore *ClientConfig

	// CLIConfig is the configuration for the CLI client.
	CfgCLI *CLIConfig

	// HostConfig is the configuration for the host.
	CfgHost *HostConfig

	// APICfg is the configuration for the API client.
	CfgAPI *APIConfig

	// DirCfg is the configuration for the current directory.
	CfgDir *DirCfg
}
