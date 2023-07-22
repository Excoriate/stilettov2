package job

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/core/commands"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/core/validation"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
	"go.uber.org/zap"
	"path/filepath"
)

type Builder struct {
	id         string
	client     *entities.Client
	job        *entities.Job
	tasks      []entities.Task
	baseDir    string
	baseDirAbs string
	envVars    map[string]string

	// For TUI and nice pipeline's output.
	logger *zap.Logger
	error  error
}

type TaskNewArgs struct {
	Name           string
	ContainerImage string
	Commands       []TaskNewCMDArgs
	BaseDir        string // Equivalent to the current Dir
	WorkDir        string // Equivalent to the directory that'll be used to perform tasks.
	MountDir       string // The directory that'll be mounted in the container.
}

type TaskNewCMDArgs struct {
	Binary      string
	CommandArgs string
}

type NewArgs struct {
	Name string
}

func (b *Builder) Build() (*entities.Job, error) {
	if b.error != nil {
		return nil, b.error
	}

	return &entities.Job{
		Id:         b.id,
		Name:       b.job.Name,
		Client:     b.client,
		Tasks:      b.tasks,
		BaseDir:    b.baseDir,
		BaseDirAbs: b.baseDirAbs,
		EnvVars:    b.envVars,
	}, nil
}

// WithTasks adds the tasks to the job.
func (b *Builder) WithTasks(opts []TaskNewArgs,
	envVarsOpt EnvVarsOptions) *Builder {
	if len(opts) == 0 {
		b.client.Logger.Info(fmt.Sprintf("No tasks provided for job '%s' with Id '%s'. "+
			"Skipping.", b.job.Name, b.job.Id))
		return b
	}

	for _, task := range opts {
		if task.ContainerImage == "" {
			taskErr := errors.NewTaskConfigurationError(fmt.Sprintf("The 'containerImage' argument is required for task '%s' with id '%s'.", task.Name, b.id), nil)
			b.client.Logger.Error(taskErr.Error())
			b.error = taskErr
			return b
		}

		// Task directories validation.
		workDir := task.WorkDir
		baseDir := task.BaseDir
		mountDir := task.MountDir

		if filepath.IsAbs(workDir) {
			taskErr := errors.NewTaskConfigurationError(fmt.Sprintf("The 'workDir' argument cannot be an absolute path: %s", workDir), nil)
			b.client.Logger.Error(taskErr.Error())
			b.error = taskErr
			return b
		}

		if filepath.IsAbs(mountDir) {
			taskErr := errors.NewTaskConfigurationError(fmt.Sprintf("The 'mountDir' argument cannot be an absolute path: %s", mountDir), nil)
			b.client.Logger.Error(taskErr.Error())
			b.error = taskErr
			return b
		}

		if baseDir == "" {
			b.logger.Warn(fmt.Sprintf("No taskBaseDir found for task '%s' with id '%s'. "+
				"It'll be replaced by the base directory set at the job level '%s'.", task.Name,
				b.id, b.baseDir))

			baseDir = b.job.BaseDirAbs
		}

		if !filepath.IsAbs(baseDir) {
			taskErr := errors.NewTaskConfigurationError(fmt.Sprintf("The 'baseDir' argument must be an absolute path: %s", baseDir), nil)
			b.client.Logger.Error(taskErr.Error())
			b.error = taskErr
			return b
		}

		if err := validation.WorkDirIsValid(validation.WorkDirIsValidArgs{
			BaseDir:  baseDir,
			WorkDir:  workDir,
			MountDir: mountDir,
		}); err != nil {
			taskErr := errors.NewTaskConfigurationError(fmt.Sprintf(
				"Cannot configure task '%s' with id '%s', ", task.Name, b.id), err)
			b.client.Logger.Error(taskErr.Error())
			b.error = taskErr

			return b
		}

		b.logger.Info(fmt.Sprintf("Configuring task '%s' with id '%s'.", task.Name, b.id))
		taskId := utils.GetUUID()

		// Env vars for the task.
		var taskEnvVars map[string]string
		if envVarsOpt.InheritEnvVarsFromJob {
			jobEnvVars := b.job.EnvVars
			if utils.MapIsNulOrEmpty(jobEnvVars) {
				b.client.Logger.Warn(fmt.Sprintf("No env vars found for job '%s' with id '%s', "+
					"it means this task %s marked to inherit env vars from job will not have any ", b.job.Name, b.job.Id, task.Name))
			} else {
				b.client.Logger.Info(fmt.Sprintf(
					"Inheriting env vars from job '%s' with id '%s' in task '%s' with id '%s'.", b.job.Name, b.job.Id, task.Name, taskId))
				taskEnvVars = jobEnvVars
			}
		} else {
			tempEnvVars, err := DecorateWithEnvVars(envVarsOpt)
			if err != nil {
				taskErr := errors.NewArgumentError(fmt.Sprintf("Error decorating env vars for task '%s' with id '%s'.", task.Name, taskId), err)
				b.client.Logger.Error(taskErr.Error())
				b.error = taskErr

				return b
			}

			b.client.Logger.Info(fmt.Sprintf("Decorating env vars for task '%s' with id '%s'.", task.Name, taskId))
			taskEnvVars = tempEnvVars
		}

		// Building the required commands for the task.
		var taskCommands []*commands.CMD
		if len(task.Commands) != 0 {
			cmdBuilder := commands.NewCMD()
			for _, cmd := range task.Commands {
				newCMD, cmdErr := cmdBuilder.WithBinary(cmd.Binary).WithCommands(cmd.CommandArgs).
					Build()

				if cmdErr != nil {
					taskErr := errors.NewArgumentError(fmt.Sprintf(
						"Error building jobcmd for task '%s' with id '%s'.", task.Name, taskId), cmdErr)
					b.client.Logger.Error(taskErr.Error())
					b.error = taskErr

					return b
				}

				taskCommands = append(taskCommands, newCMD)
			}
		}

		b.tasks = append(b.tasks, entities.Task{
			Id:             taskId,
			Name:           task.Name,
			ContainerImage: task.ContainerImage,
			Workdir:        workDir,
			MountDir:       mountDir,
			BaseDir:        baseDir,
			BaseDirAbs:     baseDir,
			EnvVars:        taskEnvVars,
			CommandsCfg:    taskCommands,
		})

		b.client.Logger.Info(fmt.Sprintf("Task '%s' with id '%s' added to the job '%s' with id"+
			" '%s'.", task.Name, taskId, b.job.Name, b.job.Id))

	}

	return b
}

// WithJob creates a new Dagger job.
func (b *Builder) WithJob(args NewArgs,
	envVarOps EnvVarsOptions) *Builder {
	jobName := args.Name

	// Job env var
	var jobEnvVars map[string]string

	// Env vars decoration based on options.
	envVars, err := DecorateWithEnvVars(envVarOps)
	if err != nil {
		jobErr := errors.NewConfigurationError(fmt.Sprintf("Error decorating env vars for job '%s' with id '%s'.", jobName, b.id), err)
		b.client.Logger.Error(jobErr.Error())
		b.error = jobErr

		return b
	}

	jobEnvVars = envVars

	job := entities.Job{
		Id:         b.id,
		Name:       args.Name,
		Client:     b.client,
		Tasks:      []entities.Task{},
		BaseDir:    b.client.CfgDir.BaseDir,
		BaseDirAbs: b.client.CfgDir.BaseDirAbs,
		EnvVars:    jobEnvVars,
	}

	b.job = &job
	b.logger.Info(fmt.Sprintf("Job '%s' with id '%s' created.", jobName, b.id))

	return b
}

func NewDaggerClient(client *entities.Client) *Builder {
	return &Builder{
		id:         utils.GetUUID(),
		envVars:    map[string]string{},
		client:     client,
		job:        &entities.Job{},   // Created when the WithJob method is called.
		tasks:      []entities.Task{}, // Created when the WithTask method is called.
		error:      nil,
		baseDir:    client.CfgDir.BaseDir,
		baseDirAbs: client.CfgDir.BaseDirAbs,
		logger:     client.Logger,
	}
}
