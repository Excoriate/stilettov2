package entities

import (
	"github.com/excoriate/stiletto/internal/core/commands"
)

type Job struct {
	// Id The id of the job.
	Id string

	// Name The name of the job.
	Name string

	// Client The client to be used to execute the commands.
	Client *Client

	// Tasks The tasks to be executed.
	Tasks []Task

	// BaseDir is the current directory where the commands will be executed.
	BaseDir string

	// BaseDirAbs is the current directory where the commands will be executed.
	BaseDirAbs string

	// EnvVars is the environment variables to be passed to the container.
	EnvVars map[string]string
}

type Task struct {
	// Id The id of the task.
	Id string

	// Name The name of the task.
	Name string

	// Binary The binary to be executed.
	ContainerImage string

	// Workdir Where the commands will be executed.
	Workdir string

	// MountDir The directory that'll be mounted in the container.
	MountDir string

	// BaseDir is the current directory where the commands will be executed.
	BaseDir string

	// BaseDir is the current directory where the commands will be executed.
	BaseDirAbs string

	// If passed, it'll pass these environment variables to the container.
	EnvVars map[string]string

	// CommandsCfg is the configuration of the jobcmd to be executed.
	// It includes the main binary, and the commands passed to it.
	CommandsCfg []*commands.CMD
}
