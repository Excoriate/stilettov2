package env

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
)

type EnvVarsByTypeOpt struct {
	FailIfNotSet    bool
	RequiredEnvVars []string
	Prefix          string
}

type EnvVarsHostOpt struct {
	RequiredEnvVars []string
	FailIfNotSet    bool
}

type EnvVarsSpecificKeysOpt struct {
	FailIfNotSet bool
	EnvVarKeys   []string
}

// GetEnvVarsByType returns the environment variables found in the system.
func GetEnvVarsByType(opt EnvVarsByTypeOpt) (map[string]string, error) {
	allEnvVarsFound, err := utils.FetchEnvVarsWithPrefix(opt.Prefix)

	if err != nil {
		return nil, errors.NewArgumentError(fmt.Sprintf("Could not fetch env vars with prefix '%s'", opt.Prefix), err)
	}

	if utils.MapIsNulOrEmpty(allEnvVarsFound) && opt.FailIfNotSet {
		return nil, errors.NewConfigurationError("No Env Vars found", nil)
	}

	if len(opt.RequiredEnvVars) == 0 {
		return allEnvVarsFound, nil
	}

	for _, envVar := range opt.RequiredEnvVars {
		if _, ok := allEnvVarsFound[envVar]; !ok {
			return nil, errors.NewConfigurationError(fmt.Sprintf(
				"The environment variable %s is not set, but was declared mandatory", envVar), nil)
		}
	}

	return allEnvVarsFound, nil
}

func GetAllEnvVarsFromHost(opt EnvVarsHostOpt) (map[string]string, error) {
	allEnvVarsFromHost, err := utils.FetchAllEnvVarsFromHost()
	if err != nil {
		return nil, errors.NewArgumentError("Could not fetch env vars from host", err)
	}

	if utils.MapIsNulOrEmpty(allEnvVarsFromHost) && opt.FailIfNotSet {
		return nil, errors.NewConfigurationError("No Env Vars found", nil)
	}

	if len(opt.RequiredEnvVars) == 0 {
		return allEnvVarsFromHost, nil
	}

	for _, envVar := range opt.RequiredEnvVars {
		if _, ok := allEnvVarsFromHost[envVar]; !ok {
			return nil, errors.NewConfigurationError(fmt.Sprintf(
				"The environment variable %s is not set, but was declared mandatory", envVar), nil)
		}
	}

	return allEnvVarsFromHost, nil
}

func GetEnvVarsBySpecificKeys(opt EnvVarsSpecificKeysOpt) (map[string]string, error) {
	if len(opt.EnvVarKeys) == 0 {
		return nil, errors.NewArgumentError("No Env Var Keys provided", nil)
	}

	allEnvVarsFound, err := utils.FetchAllEnvVarsFromHost()
	var envVarsFound map[string]string

	if err != nil {
		return nil, errors.NewArgumentError("Could not fetch custom env vars", err)
	}

	if utils.MapIsNulOrEmpty(allEnvVarsFound) && opt.FailIfNotSet {
		return nil, errors.NewConfigurationError("No Env Vars found", nil)
	}

	if len(opt.EnvVarKeys) == 0 {
		return allEnvVarsFound, nil
	}

	for _, envVar := range opt.EnvVarKeys {
		if _, ok := allEnvVarsFound[envVar]; !ok {
			if opt.FailIfNotSet {
				return nil, errors.NewConfigurationError(fmt.Sprintf(
					"The environment variable %s is not set, but was declared mandatory", envVar), nil)
			}
		} else {
			envVarsFound[envVar] = allEnvVarsFound[envVar]
		}
	}

	return envVarsFound, nil
}

func GetEnvVarsFromDotFile(dotFilePath string) (map[string]string, error) {
	if dotFilePath == "" {
		return nil, errors.NewArgumentError("No dot file path provided", nil)
	}

	if err := utils.FileExistAndItIsAFile(dotFilePath); err != nil {
		return nil, errors.NewArgumentError(fmt.Sprintf("The dot file path %s is not valid", dotFilePath), nil)
	}

	envVarsFound, err := utils.GetEnvVarsFromDotFile(dotFilePath)
	if err != nil {
		return nil, errors.NewArgumentError("Could not fetch env vars from dot file", err)
	}

	if utils.MapIsNulOrEmpty(envVarsFound) {
		return nil, errors.NewConfigurationError(fmt.Sprintf("No Env Vars found in dot file %s", dotFilePath), nil)
	}

	return envVarsFound, nil
}
