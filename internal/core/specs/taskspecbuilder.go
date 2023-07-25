package specs

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/core/job"
	"github.com/excoriate/stiletto/internal/core/validation"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/utils"
	"github.com/excoriate/stiletto/internal/yamlparser"
	"go.uber.org/zap"
	"path/filepath"
)

type Builder struct {
	manifestFile     string
	manifestType     string
	taskManifestSpec *TaskManifestSpec

	// Cross-functional configuration as part of the builder pattern.
	logger     *zap.Logger
	client     *entities.Client
	err        error
	baseDir    string
	baseDirAbs string
}

type NewOpts struct {
	ManifestType string
	Client       *entities.Client
	ManifestFile string
}

type TaskFromManifestConverter interface {
	Convert() (*ConvertedTask, error)
}

type ConvertedTask struct {
	Task       *job.TaskNewArgs
	TaskEnvCfg *job.EnvVarsOptions
}

func (s *TaskManifestSpec) Convert() (*ConvertedTask, error) {
	if s == nil {
		return nil, errors.NewManifestError(
			"Cannot convert a manifest to a Task that is empty or nil", nil)
	}

	var taskCommandArgs []job.TaskNewCMDArgs
	for _, command := range s.Spec.CommandsSpec {
		var newCmd = job.TaskNewCMDArgs{}
		if command.Binary != "" {
			newCmd.Binary = command.Binary
		}

		for _, cmd := range command.Commands {
			newCmd.CommandArgs = cmd
		}

		taskCommandArgs = append(taskCommandArgs, newCmd)
	}

	var envVarsOptions job.EnvVarsOptions

	if s.Spec.EnvVarsSpec.EnvVarsScanned.ScanTerraformEnvVars.Enabled {
		envVarsOptions.EnvVarsTerraformCfg = job.EnvVarBehaviourOptions{
			Enabled:               true,
			FailIfNotSet:          s.Spec.EnvVarsSpec.EnvVarsScanned.ScanTerraformEnvVars.FailIfNotSet,
			RequiredEnvVars:       s.Spec.EnvVarsSpec.EnvVarsScanned.ScanTerraformEnvVars.RequiredEnvVars,
			IgnoreIfNotSetOrEmpty: s.Spec.EnvVarsSpec.EnvVarsScanned.ScanTerraformEnvVars.IgnoreIfNotSetOrEmpty,
			RemoveEnvVarsIfFound:  s.Spec.EnvVarsSpec.EnvVarsScanned.ScanTerraformEnvVars.RemoveEnvVarsIfFound,
		}
	}

	if s.Spec.EnvVarsSpec.EnvVarsScanned.ScanAWSEnvVars.Enabled {
		envVarsOptions.EnvVarsAWSCfg = job.EnvVarBehaviourOptions{
			Enabled:               true,
			FailIfNotSet:          s.Spec.EnvVarsSpec.EnvVarsScanned.ScanAWSEnvVars.FailIfNotSet,
			RequiredEnvVars:       s.Spec.EnvVarsSpec.EnvVarsScanned.ScanAWSEnvVars.RequiredEnvVars,
			IgnoreIfNotSetOrEmpty: s.Spec.EnvVarsSpec.EnvVarsScanned.ScanAWSEnvVars.IgnoreIfNotSetOrEmpty,
			RemoveEnvVarsIfFound:  s.Spec.EnvVarsSpec.EnvVarsScanned.ScanAWSEnvVars.RemoveEnvVarsIfFound,
		}
	}

	if !utils.MapIsNulOrEmpty(s.Spec.EnvVarsSpec.EnvVars) {
		envVarsOptions.EnvVarsExplicit = s.Spec.EnvVarsSpec.EnvVars
	}

	if len(s.Spec.EnvVarsSpec.DotFiles) != 0 {
		envVarsOptions.EnvVarsFromDotFileCfg = job.EnvVarBehaviourOptions{
			Enabled:      true,
			FailIfNotSet: true,
			DotFiles:     s.Spec.EnvVarsSpec.DotFiles,
		}
	}

	return &ConvertedTask{
		Task: &job.TaskNewArgs{
			Name:           s.Metadata.Name,
			ContainerImage: s.Spec.ContainerImage,
			WorkDir:        s.Spec.Workdir,
			MountDir:       s.Spec.MountDir,
			BaseDir:        s.Spec.BaseDir,
			Commands:       taskCommandArgs,
		},
		TaskEnvCfg: &envVarsOptions,
	}, nil
}

// WithGeneratedTaskManifest WithTaskManifests WithJobManifests adds job manifests to the builder.
func (b *Builder) WithGeneratedTaskManifest() *Builder {
	if b.manifestType != entities.ManifestTypeTask {
		errMsg := fmt.Sprintf("invalid manifest type: %s", b.manifestType)
		b.logger.Error(errMsg)
		b.err = errors.NewArgumentError(errMsg, nil)

		return b
	}

	if err := yamlparser.YamlStructureIsValid(b.manifestFile, &TaskManifestSpec{}); err != nil {
		errMsg := fmt.Sprintf("Cannot add task manifests, "+
			"invalid yaml structure for file %s", b.manifestFile)

		b.logger.Error(errMsg)
		b.err = errors.NewArgumentError(errMsg, err)
		return b
	}

	taskManifestSpec := &TaskManifestSpec{}

	if err := yamlparser.YamlToStruct(b.manifestFile, taskManifestSpec); err != nil {
		errMsg := fmt.Sprintf("Cannot add task manifests, "+
			"cannot parse yaml file %s", b.manifestFile)

		b.logger.Error(errMsg)
		b.err = errors.NewArgumentError(errMsg, err)
		return b
	}

	if taskManifestSpec.Spec.BaseDir == "" {
		b.logger.Info("he 'baseDir' in the task manifest isn't set, so it'll be resolved to the current directory")
		taskManifestSpec.Spec.BaseDir = b.baseDirAbs
	}

	b.taskManifestSpec = taskManifestSpec
	b.logger.Info("task manifest added to the builder")

	return b
}

// WithStrictDeepValidation adds strict deep validation to the builder.
func (b *Builder) WithStrictDeepValidation() *Builder {
	specContent := b.taskManifestSpec
	if specContent == nil {
		errMsg := "task manifest is required prior to this API execution. " +
			"Ensure that you've called the WithGeneratedTaskManifest method"

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.Kind != "Task" && specContent.Kind != "Job" && specContent.
		Kind != "Workflow" {
		errMsg := fmt.Sprintf("invalid manifest kind: %s. Should be 'Job', 'Task' or 'Workflow'",
			specContent.Kind)

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.APIVersion != "v1" {
		errMsg := fmt.Sprintf("invalid manifest api version: %s. Should be 'v1'", specContent.APIVersion)

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.Metadata.Name == "" {
		errMsg := "manifest name is required. Give it a proper name. E.g.: 'my-task'"

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.Spec.ContainerImage == "" {
		errMsg := "container image is required. " +
			"It's required to bootstrap a container for the 'Dagger' runtime."

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.Spec.Workdir == "" {
		errMsg := "workDir is required. "

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if specContent.Spec.MountDir == "" {
		errMsg := "mountDir is required. "

		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)
		return b
	}

	if err := validation.WorkDirIsValid(validation.WorkDirIsValidArgs{
		BaseDir:  specContent.Spec.BaseDir,
		WorkDir:  specContent.Spec.Workdir,
		MountDir: specContent.Spec.MountDir,
	}); err != nil {
		errMsg := fmt.Sprintf("The manifest directory configuration is invalid: %s", err)
		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)

		return b
	}

	if len(specContent.Spec.CommandsSpec) == 0 {
		errMsg := "The manifest commands are invalid. They should have at least one command"
		b.logger.Error(errMsg)
		b.err = errors.NewManifestError(errMsg, nil)

		return b
	}

	for _, cmd := range specContent.Spec.CommandsSpec {
		if len(cmd.Commands) == 0 {
			errMsg := "The manifest commands are invalid. It was detected a configuration, " +
				"but without any command to execute"
			b.logger.Error(errMsg)
			b.err = errors.NewManifestError(errMsg, nil)
			return b
		}
	}

	return b
}

// Build builds the manifest.
func (b *Builder) Build() (*TaskManifestSpec, error) {
	if b.err != nil {
		return &TaskManifestSpec{}, errors.NewConfigurationError(fmt.Sprintf(
			"cannot build manifest of type '%s'", b.manifestType), b.err)
	}

	return &TaskManifestSpec{
		APIVersion: b.taskManifestSpec.APIVersion,
		Kind:       b.taskManifestSpec.Kind,
		Metadata: TaskMetadata{
			Name: b.taskManifestSpec.Metadata.Name,
		},
		Spec: TaskSpec{
			ContainerImage: b.taskManifestSpec.Spec.ContainerImage,
			Workdir:        b.taskManifestSpec.Spec.Workdir,
			MountDir:       b.taskManifestSpec.Spec.MountDir,
			BaseDir:        b.taskManifestSpec.Spec.BaseDir,
			CommandsSpec:   b.taskManifestSpec.Spec.CommandsSpec,
			EnvVarsSpec:    b.taskManifestSpec.Spec.EnvVarsSpec,
		},
	}, nil
}

// NewTaskSpecBuilder creates a new instance of ManifestBuilder.
func NewTaskSpecBuilder(opts NewOpts) (*Builder, error) {
	if opts.Client == nil {
		errMsg := "A valid client instance is required."
		return nil, errors.NewArgumentError(errMsg, nil)
	}

	logger := opts.Client.Logger

	if opts.ManifestFile == "" {
		errMsg := "The manifest file is required. " +
			"Ensure it's passed as a relative path of the current directory"
		logger.Error(errMsg)
		return nil, errors.NewArgumentError(errMsg, nil)
	}

	if opts.ManifestType != entities.ManifestTypeTask && opts.ManifestType != entities.
		ManifestTypeJob && opts.ManifestType != entities.ManifestTypeWorkflow {
		errMsg := fmt.Sprintf("Cannot create a manifest builder client. Invalid manifest type: %s",
			opts.ManifestType)
		logger.Error(errMsg)
		return nil, errors.NewArgumentError(errMsg, nil)
	}

	// Joining the manifest filepath with the current directory.
	manifestFileFull := filepath.Join(opts.Client.CfgDir.BaseDir, opts.ManifestFile)
	logger.Info(fmt.Sprintf("The manifest full path is resolved as: %s", manifestFileFull))

	opts.ManifestFile = manifestFileFull

	if err := yamlparser.YamlFileIsValid(opts.ManifestFile); err != nil {
		errMsg := fmt.Sprintf("Cannot create a manifest builder client. Invalid manifest file: %s",
			opts.ManifestFile)

		logger.Error(errMsg)
		return nil, errors.NewArgumentError(errMsg, err)
	}

	return &Builder{
		manifestType:     opts.ManifestType,
		manifestFile:     opts.ManifestFile,
		taskManifestSpec: nil,
		client:           opts.Client,
		logger:           opts.Client.Logger,
		baseDir:          opts.Client.CfgDir.BaseDir,
		baseDirAbs:       opts.Client.CfgDir.BaseDirAbs,
	}, nil
}
