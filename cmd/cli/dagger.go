package cli

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/core/job"
	"github.com/excoriate/stiletto/internal/core/runner"
	"github.com/excoriate/stiletto/internal/core/scheduler"
	"github.com/excoriate/stiletto/internal/core/specs"
	"github.com/excoriate/stiletto/internal/tui"
	"github.com/excoriate/stiletto/pkg/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	taskFiles []string
)

var DaggerCMD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "dagger",
	Long: `The 'dagger' command is an special type of 'Job' that runs tasks on top of
Dagger (Container).`,
	Example: `
stiletto job dagger --task-files=../../stiletto/tasks/terragrunt-plan.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		// CLI UX utilities.
		cliLog := tui.NewTUIMessage()
		cliUX := tui.NewTitle()

		// Specific flags
		workDir := viper.GetString("workDir")
		mountDir := viper.GetString("mountDir")
		showEnvVars := viper.GetBool("showEnvVars")

		// Task files/manifests to mount.
		taskFilesCfg := viper.GetStringSlice("taskFiles")

		if len(taskFilesCfg) == 0 {
			cliLog.ShowError("", "No task files (specs, or manifests) were provided",
				nil)
			os.Exit(1)
		}

		// New client builder.
		c := clients.NewClient(entities.ClientTypeCli)

		// Client instance.
		i, err := c.WithCLI(entities.CLIConfigArgs{}).WithHost().Build()

		if err != nil {
			cliLog.ShowError("CLIENT-ERROR", err.Error(), nil)
			os.Exit(1)
		}

		cliUX.ShowTitleAndDescription("STILETTO",
			"Automated pipelines, "+
				"workflows and whatever can be containerized in your own laptop 👨🏻‍💻("+
				"powered by Dagger.IO)")

		var tasksConvertedFromManifest []specs.ConvertedTask
		for _, taskFile := range taskFilesCfg {
			manifestBuilder, err := specs.NewTaskSpecBuilder(specs.NewOpts{
				ManifestType: entities.ManifestTypeTask,
				ManifestFile: taskFile,
				Client:       i,
			})

			if err != nil {
				cliLog.ShowError("", err.Error(), nil)
				os.Exit(1)
			}

			// Task manifest, ready to be transformed into a valid Dagger job (task).
			taskManifest, err := manifestBuilder.WithGeneratedTaskManifest().
				WithStrictDeepValidation().
				Build()

			if err != nil {
				cliLog.ShowError("", err.Error(), nil)
				os.Exit(1)
			}

			convertedTask, err := taskManifest.Convert()

			if err != nil {
				cliLog.ShowError("", err.Error(), nil)
				os.Exit(1)
			}

			if workDir != "" && convertedTask.Task.WorkDir != "" {
				cliLog.ShowWarning("", fmt.Sprintf("The workDir '%s' was set in the CLI, "+
					"but also set in the manifest file '%s'. The CLI value will be used.", workDir, taskFile))

				convertedTask.Task.WorkDir = workDir
			}

			if mountDir != "" && convertedTask.Task.MountDir != "" {
				cliLog.ShowWarning("", fmt.Sprintf("The mountDir '%s' was set in the CLI, "+
					"but also set in the manifest file '%s'. The CLI value will be used.", mountDir, taskFile))

				convertedTask.Task.MountDir = mountDir
			}

			tasksConvertedFromManifest = append(tasksConvertedFromManifest, *convertedTask)
		}

		daggerClient := job.NewDaggerClient(i)

		var jobs []entities.Job
		for _, task := range tasksConvertedFromManifest {
			j, err := daggerClient.WithJob(job.NewArgs{
				Name: fmt.Sprintf("job-task-%s", task.Task.Name),
			}, job.EnvVarsOptions{}).WithTasks([]job.TaskNewArgs{*task.Task},
				*task.TaskEnvCfg).Build()

			if err != nil {
				cliLog.ShowError("JOB-ERROR", err.Error(), nil)
				os.Exit(1)
			}

			jobs = append(jobs, *j)
		}

		// Running the Dagger job.
		scheduler := scheduler.NewScheduler()

		// Create a runner with the jobs to run.
		scheduleJobs, err := scheduler.WithClient(i).WithJobsToRun(jobs).WithDaggerEngine().Build()
		if err != nil {
			cliLog.ShowError("SCHEDULER-ERROR", err.Error(), nil)
			os.Exit(1)
		}

		// Run the jobs.
		runnerBuilder := runner.NewRunnerDagger(scheduleJobs)
		runner, err := runnerBuilder.WithDaggerClient(nil).WithOptions(runner.DaggerRunnerOptions{
			ShowEnvVars: showEnvVars,
		}).Build()
		if err != nil {
			cliLog.ShowError("RUNNER-ERROR", err.Error(), nil)
			os.Exit(1)
		}

		err = runner.RunInDagger(jobs)
		if err != nil {
			cliLog.ShowError("RUNNER-ERROR", err.Error(), nil)
		}
	},
}

func addFlagsToDaggerCMD() {
	DaggerCMD.Flags().StringSliceVarP(&taskFiles, "task-files",
		"", []string{}, "The tasks  in .yml format that'll be executed")

	_ = viper.BindPFlag("taskFiles", DaggerCMD.Flags().Lookup("task-files"))

	if err := DaggerCMD.MarkFlagRequired("task-files"); err != nil {
		panic(err)
	}
}

func init() {
	addFlagsToDaggerCMD()
}
