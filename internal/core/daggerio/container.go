package daggerio

import (
	"context"
	"dagger.io/dagger"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
)

// SetEnvVarsInContainer sets the environment variables in the container.
func SetEnvVarsInContainer(c *dagger.Container, envVars map[string]string) (*dagger.Container,
	error) {

	if utils.MapIsNulOrEmpty(envVars) {
		return nil, errors.NewArgumentError("No environment variables are passed, skipping the environment variable configuration step", nil)
	}

	for k, v := range envVars {
		c = c.WithEnvVariable(k, v)
	}

	return c, nil
}

// GetEnvVarsSetInContainer returns the environment variables set in the container.
func GetEnvVarsSetInContainer(c *dagger.Container, ctx *context.Context) ([]dagger.EnvVariable,
	error) {

	if c == nil {
		return nil, errors.NewArgumentError("No container was passed", nil)
	}

	envVars, err := c.EnvVariables(*ctx)
	if err != nil {
		return nil, errors.NewRunnerExecutionError("Could not get the environment variables from the container", err)
	}

	return envVars, nil
}
