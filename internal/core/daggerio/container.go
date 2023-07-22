package daggerio

import (
	"dagger.io/dagger"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
)

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
