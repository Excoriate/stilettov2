package runner

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/excoriate/stiletto/internal/core/adapters"
	"github.com/excoriate/stiletto/internal/core/daggerio"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/core/scheduler"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
	"go.uber.org/zap"
	"path/filepath"
)

type DaggerRunner struct {
	Id           string
	Client       *entities.Client
	Jobs         []entities.Job
	DaggerClient *dagger.Client
	Logger       *zap.Logger
	Ctx          *context.Context
	BaseDir      string
	BaseDirAbs   string
	Options      DaggerRunnerOptions
}

type DaggerRunnerBuilder struct {
	id           string
	client       *entities.Client
	jobs         []entities.Job
	daggerClient *dagger.Client
	error        error
	logger       *zap.Logger
	ctx          *context.Context
	baseDir      string
	baseDirAbs   string
	Options      DaggerRunnerOptions
}

type DaggerRunnerOptions struct {
	ShowEnvVars bool
}

func (r *DaggerRunner) RunInDagger(jobs []entities.Job) error {
	if len(jobs) == 0 {
		return errors.NewRunnerConfigurationError("No jobs to run", nil)
	}

	daggerClient := r.DaggerClient

	if daggerClient == nil && r.Client.CfgDagger.Client == nil {
		return errors.NewRunnerConfigurationError("No Dagger engine ("+
			"or client) found in either the client instance or the scheduler.", nil)
	}

	// Builder for host/common operations.
	daggerFs, err := daggerio.NewDaggerFsBuilder(r.Client.Ctx).
		WithLogger(r.Logger).
		WithDaggerClient(daggerClient).
		Build()

	if err != nil {
		return errors.NewRunnerConfigurationError("Failed to run jobs in Dagger", err)
	}

	defer daggerClient.Close()

	for _, job := range jobs {
		if len(job.Tasks) == 0 {
			errMsg := fmt.Sprintf("Job %s with id %s has no tasks. Continuing... ", job.Name,
				job.Id)
			r.Logger.Warn(errMsg)

			continue
		}

		// Get the current host directory.
		baseDirAbs := job.BaseDirAbs
		r.Logger.Info(fmt.Sprintf("Job %s will be executed from base directory %s", job.Name, baseDirAbs))

		for _, task := range job.Tasks {
			// Directory to copy to the container, aka 'mount directory'.
			mountDirPathAbs := filepath.Join(baseDirAbs, task.MountDir)
			r.Logger.Info(fmt.Sprintf("Task %s with id %s will be executed from mount directory %s", task.Name, task.Id, mountDirPathAbs))

			if err := daggerFs.ValidateEntries(mountDirPathAbs); err != nil {
				return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run task %s with id %s", task.Name, task.Id), err)
			}

			mountDir, _ := daggerFs.GetDaggerDir(mountDirPathAbs)

			// Mounting/copying the directory to the container.
			container := daggerClient.Container().From(task.ContainerImage)
			container = container.WithDirectory(daggerFs.GetMntDir(), mountDir)

			_ = daggerFs.PrintEntries(mountDir)

			// WorkDir validation within dagger.
			workDirPathAbs := filepath.Join(mountDirPathAbs, task.Workdir)
			r.Logger.Info(fmt.Sprintf("Task %s with id %s will be executed from work directory %s", task.Name, task.Id, workDirPathAbs))

			if err := daggerFs.ValidateEntries(workDirPathAbs); err != nil {
				return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run task %s with id %s", task.Name, task.Id), err)
			}

			workDir, _ := daggerFs.GetDaggerDir(workDirPathAbs)
			_ = daggerFs.PrintEntries(workDir)

			if !utils.MapIsNulOrEmpty(task.EnvVars) {
				container, _ = daggerio.SetEnvVarsInContainer(container, task.EnvVars)
			}

			if r.Options.ShowEnvVars {
				envVars, err := daggerio.GetEnvVarsSetInContainer(container, r.Ctx)
				if err != nil {
					return errors.NewTaskExecutionError(fmt.Sprintf("Failed to run task %s with id %s", task.Name, task.Id), err)
				}

				for _, envVar := range envVars {
					name, _ := envVar.Name(*r.Ctx)
					value, _ := envVar.Value(*r.Ctx)
					r.Logger.Info(fmt.Sprintf("EnvVar: %s=%s", name, value))
				}
			}

			workDirPath := filepath.Join(daggerFs.GetMntDir(), task.Workdir)
			container = container.WithWorkdir(workDirPath)

			// Run specific set of commands per task.
			for _, cmd := range task.CommandsCfg {
				_, err := container.WithExec(cmd.Commands).Sync(*job.Client.Ctx)
				if err != nil {
					r.Logger.Error(fmt.Sprintf("Task %s with id %s failed to run", task.Name, task.Id))
					return errors.NewTaskExecutionError(fmt.Sprintf("Task %s with id %s failed to run", task.Name, task.Id), err)
				}
			}
		}
	}

	r.Logger.Info("All jobs were executed successfully")
	return nil

}

func (b *DaggerRunnerBuilder) WithOptions(opt DaggerRunnerOptions) *DaggerRunnerBuilder {
	if opt.ShowEnvVars {
		b.logger.Info("The environment variables will be shown")
	}

	b.Options = opt

	return b
}

func (b *DaggerRunnerBuilder) WithDaggerClient(c *dagger.Client) *DaggerRunnerBuilder {
	if c != nil {
		b.daggerClient = c
		b.logger.Info("An explicit 'Dagger' client was passed to the Dagger runner builder")
		return b
	}

	daggerClientInScheduledJobs := b.daggerClient
	if daggerClientInScheduledJobs != nil {
		b.daggerClient = daggerClientInScheduledJobs
		return b
	}

	if b.client.CfgDagger.Client == nil {
		newDaggerClient, err := adapters.NewDaggerClient(b.ctx)
		if err != nil {
			daggerErr := errors.NewConfigurationError("Failed to create a new dagger client", err)
			b.error = daggerErr
			b.logger.Error(daggerErr.Error())

			return b
		}

		b.daggerClient = newDaggerClient
		return b
	}

	return b
}

func (b *DaggerRunnerBuilder) Build() (*DaggerRunner, error) {
	if b.error != nil {
		return nil, b.error
	}

	return &DaggerRunner{
		Id:           b.id,
		Client:       b.client,
		Jobs:         b.jobs,
		DaggerClient: b.daggerClient,
		Logger:       b.logger,
		Ctx:          b.ctx,
		BaseDir:      b.baseDir,
		BaseDirAbs:   b.baseDirAbs,
		Options:      b.Options,
	}, nil
}

func NewRunnerDagger(s *scheduler.ScheduledJobs) *DaggerRunnerBuilder {
	c := s.Client

	return &DaggerRunnerBuilder{
		id:           utils.GetUUID(),
		jobs:         []entities.Job{},
		client:       c,
		error:        nil,
		logger:       c.Logger,
		daggerClient: s.DaggerClient,
		ctx:          c.Ctx,
		baseDir:      s.Client.CfgDir.BaseDir,
		baseDirAbs:   s.Client.CfgDir.BaseDirAbs,
	}
}
