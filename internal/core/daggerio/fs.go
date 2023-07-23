package daggerio

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
	"go.uber.org/zap"
	"path/filepath"
)

type FsBuilder struct {
	logger       *zap.Logger
	ctx          *context.Context
	daggerClient *dagger.Client
}

type Fs struct {
	DaggerClient *dagger.Client
	Logger       *zap.Logger
	Ctx          *context.Context
}

type FsValidator interface {
	ValidateEntries(dir string) error
	PrintEntries(dir *dagger.Directory) error
	GetDaggerDir(dir string) (*dagger.Directory, error)
	GetMntDir() string
}

func (b *FsBuilder) WithLogger(logger *zap.Logger) *FsBuilder {
	b.logger = logger
	return b
}

func (b *FsBuilder) WithDaggerClient(c *dagger.Client) *FsBuilder {
	b.daggerClient = c
	return b
}

func (b *FsBuilder) WithDir(dir string) *FsBuilder {
	if dir == "" {
		dir = "."
		b.logger.Info("The DaggerIO Fs cheker did not receive a directory, " +
			"using the current directory")

		return b
	}

	if err := utils.IsValidDir(dir); err != nil {
		b.logger.Error("The directory passed to the DaggerIO Fs checker is not valid", zap.Error(err))
		return b
	}

	return b
}

func (b *FsBuilder) Build() (*Fs, error) {
	if b.daggerClient == nil {
		return nil, errors.NewConfigurationError("No dagger client was passed to the DaggerIO Fs builder", nil)
	}

	if b.logger == nil {
		b.logger = zap.NewNop()
	}

	return &Fs{
		DaggerClient: b.daggerClient,
		Logger:       b.logger,
		Ctx:          b.ctx,
	}, nil
}

func (h *Fs) ValidateEntries(dir string) error {
	if dir == "" {
		errMsg := "Cannot validate entries. No directory was passed to the DaggerIOHost instance"
		h.Logger.Error(errMsg)

		return errors.NewArgumentError(errMsg, nil)
	}

	if err := utils.IsValidDir(dir); err != nil {
		errMsg := fmt.Sprintf("Failed to validate dagger entry for dir %s, error: %s", dir, err)
		h.Logger.Error(errMsg)
		return errors.NewConfigurationError(errMsg, err)
	}

	if !filepath.IsAbs(dir) {
		dir, err := utils.PathToAbsolute(dir)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to validate dagger entry for dir %s, error: %s", dir, err)
			h.Logger.Error(errMsg)
			return errors.NewConfigurationError(errMsg, err)
		}

		dir = filepath.Clean(dir)
	}

	dirInDagger := h.DaggerClient.Host().Directory(dir)
	entries, err := dirInDagger.Entries(*h.Ctx)

	if err != nil {
		errMsg := fmt.Sprintf("Could not get entries of directory %s", dir)
		h.Logger.Error(errMsg, zap.Error(err))

		return errors.NewConfigurationError(errMsg, err)
	}

	if len(entries) == 0 {
		errMsg := fmt.Sprintf("No entries found in directory %s", dir)
		h.Logger.Error(errMsg)

		return errors.NewConfigurationError(errMsg, nil)
	}

	h.Logger.Info("Entries found in directory", zap.String("dir", dir), zap.Any("entries", entries))

	return nil
}

func (h *Fs) GetDaggerDir(dir string) (*dagger.Directory, error) {
	if dir == "" {
		errMsg := "Cannot get Dagger directory. No directory was passed to the DaggerIOHost instance"
		h.Logger.Error(errMsg)

		return nil, errors.NewArgumentError(errMsg, nil)
	}

	dirInDagger := h.DaggerClient.Host().Directory(dir)

	return dirInDagger, nil
}

func (h *Fs) GetMntDir() string {
	return "/mnt"
}

func (h *Fs) PrintEntries(dir *dagger.Directory) error {
	if dir == nil {
		errMsg := "Cannot print entries. No directory was passed to the DaggerIOHost instance"
		h.Logger.Error(errMsg)
		return errors.NewArgumentError(errMsg, nil)
	}

	entries, err := dir.Entries(*h.Ctx)

	if err != nil {
		errMsg := fmt.Sprintf("Could not get entries of directory. The 'dir."+
			"Entries' function in Dagger returned an error: %s", err)

		h.Logger.Error(errMsg, zap.Error(err))

		return errors.NewConfigurationError(errMsg, err)
	}

	if len(entries) == 0 {
		errMsg := "No entries found in directory. The entries counted 0 (empty)."
		h.Logger.Error(errMsg)

		return errors.NewConfigurationError(errMsg, nil)
	}

	for _, entry := range entries {
		h.Logger.Info(fmt.Sprintf("Entry found in directory %s", entry))
	}

	return nil
}

func NewDaggerFsBuilder(ctx *context.Context) *FsBuilder {
	return &FsBuilder{
		ctx: ctx,
	}
}
