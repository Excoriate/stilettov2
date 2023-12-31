package job

import (
	"github.com/excoriate/stiletto/internal/core/env"
	"github.com/excoriate/stiletto/internal/utils"
)

type EnvVarBehaviourOptions struct {
	RequiredEnvVars       []string
	FailIfNotSet          bool
	Enabled               bool
	DotFiles              []string
	IgnoreIfNotSetOrEmpty []string
	RemoveEnvVarsIfFound  []string
}

type EnvVarsOptions struct {
	InheritEnvVarsFromJob bool
	EnvVarsHostCfg        EnvVarBehaviourOptions
	EnvVarsAWSCfg         EnvVarBehaviourOptions
	EnvVarsTerraformCfg   EnvVarBehaviourOptions
	EnvVarsFromDotFileCfg EnvVarBehaviourOptions
	EnvVarsExplicit       map[string]string
}

// DecorateWithEnvVars  decorates the job with the env vars.
func DecorateWithEnvVars(opts EnvVarsOptions) (map[string]string,
	error) {
	var envVars map[string]string

	// AWS env vars
	if opts.EnvVarsAWSCfg.Enabled {
		awsEnvVars, err := env.GetEnvVarsByType(env.VarsByTypeOpt{
			FailIfNotSet:          opts.EnvVarsAWSCfg.FailIfNotSet,
			RequiredEnvVars:       opts.EnvVarsAWSCfg.RequiredEnvVars,
			IgnoreIfNotSetOrEmpty: opts.EnvVarsAWSCfg.IgnoreIfNotSetOrEmpty,
			RemoveEnvVarsIfFound:  opts.EnvVarsAWSCfg.RemoveEnvVarsIfFound,
			Prefix:                "AWS_",
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, awsEnvVars)
	}

	// TF env vars
	if opts.EnvVarsTerraformCfg.Enabled {
		tfEnvVars, err := env.GetEnvVarsByType(env.VarsByTypeOpt{
			FailIfNotSet:          opts.EnvVarsTerraformCfg.FailIfNotSet,
			RequiredEnvVars:       opts.EnvVarsTerraformCfg.RequiredEnvVars,
			IgnoreIfNotSetOrEmpty: opts.EnvVarsTerraformCfg.IgnoreIfNotSetOrEmpty,
			RemoveEnvVarsIfFound:  opts.EnvVarsTerraformCfg.RemoveEnvVarsIfFound,
			Prefix:                "TF_",
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, tfEnvVars)
	}

	// Host env vars
	if opts.EnvVarsHostCfg.Enabled {
		hostEnvVars, err := env.GetAllEnvVarsFromHost(env.VarsHostOpt{
			FailIfNotSet:    opts.EnvVarsHostCfg.FailIfNotSet,
			RequiredEnvVars: opts.EnvVarsHostCfg.RequiredEnvVars,
		})

		if err != nil {
			return nil, err
		}

		envVars = utils.MergeEnvVars(envVars, hostEnvVars)
	}

	// Env vars custom
	if !utils.MapIsNulOrEmpty(opts.EnvVarsExplicit) {
		envVars = utils.MergeEnvVars(envVars, opts.EnvVarsExplicit)
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
