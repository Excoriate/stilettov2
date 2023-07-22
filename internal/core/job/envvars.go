package job

import (
	"github.com/excoriate/stiletto/internal/core/env"
	"github.com/excoriate/stiletto/internal/utils"
)

type EnvVarBehaviourOptions struct {
	RequiredEnvVars []string
	FailIfNotSet    bool
	Enabled         bool
	DotFiles        []string
}

type EnvVarsOptions struct {
	InheritEnvVarsFromJob bool
	EnvVarsHostCfg        EnvVarBehaviourOptions
	EnvVarsAWSCfg         EnvVarBehaviourOptions
	EnvVarsTerraformCfg   EnvVarBehaviourOptions
	EnvVarsCustomCfg      EnvVarBehaviourOptions
	EnvVarsFromDotFileCfg EnvVarBehaviourOptions
}

func DecorateWithEnvVars(opts EnvVarsOptions) (map[string]string,
	error) {
	var envVars map[string]string

	// AWS env vars
	if opts.EnvVarsAWSCfg.Enabled {
		awsEnvVars, err := env.GetEnvVarsByType(env.EnvVarsByTypeOpt{
			FailIfNotSet:    opts.EnvVarsAWSCfg.FailIfNotSet,
			RequiredEnvVars: opts.EnvVarsAWSCfg.RequiredEnvVars,
			Prefix:          "AWS_",
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, awsEnvVars)
	}

	// TF env vars
	if opts.EnvVarsTerraformCfg.Enabled {
		tfEnvVars, err := env.GetEnvVarsByType(env.EnvVarsByTypeOpt{
			FailIfNotSet:    opts.EnvVarsTerraformCfg.FailIfNotSet,
			RequiredEnvVars: opts.EnvVarsTerraformCfg.RequiredEnvVars,
			Prefix:          "TF_",
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, tfEnvVars)
	}

	// Host env vars
	if opts.EnvVarsHostCfg.Enabled {
		hostEnvVars, err := env.GetAllEnvVarsFromHost(env.EnvVarsHostOpt{
			FailIfNotSet:    opts.EnvVarsHostCfg.FailIfNotSet,
			RequiredEnvVars: opts.EnvVarsHostCfg.RequiredEnvVars,
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, hostEnvVars)
	}

	// Env vars custom
	if opts.EnvVarsCustomCfg.Enabled {
		customEnvVars, err := env.GetEnvVarsBySpecificKeys(env.EnvVarsSpecificKeysOpt{
			FailIfNotSet: opts.EnvVarsCustomCfg.FailIfNotSet,
			EnvVarKeys:   opts.EnvVarsCustomCfg.RequiredEnvVars,
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, customEnvVars)
	}

	// DotFiles
	if opts.EnvVarsFromDotFileCfg.Enabled {
		var dotFileEnvVars map[string]string
		for _, dotFile := range opts.EnvVarsFromDotFileCfg.DotFiles {
			envVars, err := env.GetEnvVarsFromDotFile(dotFile)
			if err != nil {
				return nil, err
			}

			dotFileEnvVars = utils.MergeEnvVars(dotFileEnvVars, envVars)
		}

		envVars = utils.MergeEnvVars(envVars, dotFileEnvVars)
	}

	return envVars, nil
}
