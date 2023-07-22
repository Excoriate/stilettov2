package clients

import (
	"context"
	"github.com/excoriate/stiletto/internal/core/adapters"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/observability"
	"github.com/excoriate/stiletto/internal/tui"
	"github.com/excoriate/stiletto/internal/utils"
	"go.uber.org/zap"
)

// Builder InstanceClient is the client for the pipeline instance.
// It's used to leverage the builder pattern to create the pipeline instance.
type Builder struct {
	id        string
	ctx       *context.Context
	clientCfg *entities.ClientConfig
	cliCfg    *entities.CLIConfig
	apiCfg    *entities.APIConfig
	hostCfg   *entities.HostConfig
	daggerCfg *entities.DaggerConfig
	dirCfg    *entities.DirCfg

	// Cross-functional configuration as part of the builder pattern.
	logger *zap.Logger
	error  error
}

// WithCLI is used to configure the CLI client.
func (b *Builder) WithCLI(opt entities.CLIConfigArgs) *Builder {
	if opt.DirCfg == nil {
		opt.DirCfg = b.dirCfg
	}

	cliCfg := &entities.CLIConfig{
		UXLog:    tui.NewTUIMessage(),
		UXTitles: tui.NewTitle(),
		DirCfg:   opt.DirCfg,
	}

	b.cliCfg = cliCfg

	return b
}

// WithHost is used to configure the host client.
func (b *Builder) WithHost() *Builder {
	if b.hostCfg != nil {
		return b
	}

	var envVarsFromHost map[string]string
	if utils.MapIsNulOrEmpty(b.clientCfg.HostEnvVars) {
		envVarsFromHost, err := utils.FetchAllEnvVarsFromHost()
		if err != nil {
			b.error = errors.NewConfigurationError(
				"Failed to configure the 'host' for this client, could not fetch environment"+
					" variables from host", err)
			return b
		}

		b.clientCfg.HostEnvVars = envVarsFromHost
	}

	if b.dirCfg == nil {
		b.dirCfg = entities.GetDirCfg()
	}

	hostCfg := &entities.HostConfig{
		EnvVars: envVarsFromHost,
		DirCfg:  b.dirCfg,
	}

	b.hostCfg = hostCfg
	return b
}

// WithDagger is used to configure the dagger client.
func (b *Builder) WithDagger() *Builder {
	if b.daggerCfg != nil {
		return b
	}

	client, err := adapters.NewDaggerClient(b.ctx)

	if err != nil {
		b.error = errors.NewConfigurationError("Could not create this client with Dagger", err)
		return b
	}

	daggerCfg := &entities.DaggerConfig{
		Client: client,
		DirCfg: b.dirCfg,
	}

	b.daggerCfg = daggerCfg
	return b
}

// Build is used to build the pipeline instance.
func (b *Builder) Build() (*entities.Client, error) {
	if b.error != nil {
		return nil, b.error
	}

	return &entities.Client{
		Id:        utils.GetUUID(),
		CfgCore:   b.clientCfg,
		CfgCLI:    b.cliCfg,
		CfgAPI:    b.apiCfg,
		CfgHost:   b.hostCfg,
		Logger:    b.logger,
		CfgDagger: b.daggerCfg,
		CfgDir:    b.dirCfg,
		Ctx:       b.ctx,
	}, nil
}

// NewClient is used to create a new pipeline instance.
func NewClient(clientType string) *Builder {
	logger := observability.NewLogger()

	ctx := context.Background()

	envVarsFromHost, _ := utils.FetchAllEnvVarsFromHost()

	dirCfg := entities.GetDirCfg()

	coreCfg := &entities.ClientConfig{
		Classification: utils.NormaliseStringUpper(clientType),
		HostEnvVars:    envVarsFromHost,
		DirCfg:         dirCfg,
	}

	return &Builder{
		id:        utils.GetUUID(),
		clientCfg: coreCfg,
		logger:    logger,
		ctx:       &ctx,
		// Configurations that can be passed selectively.
		cliCfg: &entities.CLIConfig{
			DirCfg: dirCfg,
		},
		hostCfg: &entities.HostConfig{
			DirCfg: dirCfg,
		},
		apiCfg: &entities.APIConfig{},
		daggerCfg: &entities.DaggerConfig{
			DirCfg: dirCfg,
		},
		dirCfg: dirCfg,
		error:  nil,
	}
}
