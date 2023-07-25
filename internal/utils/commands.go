package utils

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/google/shlex"
)

// GetCommandArgs parses the job command and returns the arguments.
func GetCommandArgs(cmd string) ([]string, error) {
	args, err := shlex.Split(cmd)

	if err != nil {
		return nil, errors.NewArgumentError(fmt.Sprintf("Could not parse the jobcmd: %s", err), nil)
	}

	return args, nil
}
