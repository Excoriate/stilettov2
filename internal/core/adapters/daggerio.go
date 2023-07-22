package adapters

import (
	"context"
	"dagger.io/dagger"
	"github.com/excoriate/stiletto/internal/errors"
	"os"
)

func NewDaggerClient(ctx *context.Context) (*dagger.
Client, error) {
	client, err := dagger.Connect(*ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return nil, errors.NewConfigurationError(
			"Unable to connect to dagger client (no rootDir passed to be mounted", nil)
	}

	return client, nil
}
