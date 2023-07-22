package validation

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
	"os"
	"path/filepath"
)

// DirIsValid validates the workDir argument.
func DirIsValid(dirToEvaluate string) error {
	if dirToEvaluate == "" {
		return errors.NewArgumentError("The 'dir' argument is required. "+
			"It can be defined or set as '.' to use the current directory.", nil)
	}

	if err := utils.IsValidDir(dirToEvaluate); err != nil {
		return errors.NewArgumentError(fmt.Sprintf(
			"Invalid directory passed: %s", err), nil)
	}

	return nil
}

type WorkDirIsValidArgs struct {
	BaseDir  string
	WorkDir  string
	MountDir string
}

func WorkDirIsValid(opts WorkDirIsValidArgs) error {
	if opts.BaseDir == "" {
		return errors.NewArgumentError("The 'baseDir' argument is required. "+
			"It can be defined or set as '.' to use the current directory.", nil)
	}

	if opts.WorkDir == "" {
		return errors.NewArgumentError("The 'workDir' argument is required. "+
			"It can be defined or set as '.' to use the current directory.", nil)
	}

	if opts.MountDir == "" {
		return errors.NewArgumentError("The 'mountDir' argument is required. "+
			"It can be defined or set as '.' to use the current directory.", nil)
	}

	if opts.BaseDir == "." {
		opts.BaseDir, _ = os.Getwd()
	}

	if filepath.IsAbs(opts.WorkDir) {
		return errors.NewArgumentError(fmt.Sprintf(
			"The 'workDir' argument cannot be an absolute path: %s", opts.WorkDir), nil)
	}

	if filepath.IsAbs(opts.MountDir) {
		return errors.NewArgumentError(fmt.Sprintf(
			"The 'mountDir' argument cannot be an absolute path: %s", opts.MountDir), nil)
	}

	if !filepath.IsAbs(opts.BaseDir) {
		return errors.NewArgumentError(fmt.Sprintf(
			"The 'baseDir' argument must be an absolute path: %s", opts.BaseDir), nil)
	}

	mountDirFullPath := filepath.Join(opts.BaseDir, opts.MountDir)
	if err := utils.IsValidDir(mountDirFullPath); err != nil {
		return errors.NewArgumentError(fmt.Sprintf(
			"Invalid mountdir passed: %s", err), nil)
	}

	workDirFullPath := filepath.Join(opts.BaseDir, opts.MountDir, opts.WorkDir)

	if err := utils.IsValidDir(workDirFullPath); err != nil {
		return errors.NewArgumentError(fmt.Sprintf(
			"Invalid workdir passed: %s", err), nil)
	}

	return nil
}
