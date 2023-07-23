package cli

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/tui"
	"github.com/excoriate/stiletto/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Job
	jobName string

	// dotFiles is the list of dotfiles to be copied to the container.
	dotFiles []string

	// workDir is the directory where the commands will be executed.
	workDir string

	// mountDir is the directory that'll be mounted in the container.
	mountDir string

	// showEnvVars is a flag that indicates if the environment variables should be shown.
	showEnvVars bool
)

var JobCMD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "job",
	Long: `The 'job' command perform a job (
single job-task) that can be executed with any of the supported runners (E.g.: Dagger).
A 'Job' is an entity that groups one or many Tasks.`,
	Example: `
	  stiletto job --debug --workdir=/tmp --mountdir=/tmp --dotfiles=.env,.env2 --job-name=job1`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// CLI UX utilities.
		cliLog := tui.NewTUIMessage()

		// Setting the job's name.
		if jobName == "" {
			jobNameSuffix := utils.GenerateRandomString(utils.RandomStringOptions{
				Length:            5,
				UpperLowerRandom:  false,
				AllowSpecialChars: false,
			})

			jobName = fmt.Sprintf("job-%s", jobNameSuffix)

			cliLog.ShowInfo("", fmt.Sprintf("The jobName is not set, using the random name: %s", jobName))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func addPersistentFlagsToJobCMD() {
	JobCMD.PersistentFlags().StringVarP(&jobName,
		"job-name",
		"", "",
		"Job name that will be used to identify the job. If it's not set, it'll default to a random name.")

	JobCMD.PersistentFlags().StringSliceVarP(&dotFiles,
		"dotfiles",
		"", []string{""},
		"List of dotfiles that will be copied to the container. ")

	JobCMD.PersistentFlags().StringVarP(&workDir,
		"workdir",
		"", "",
		"WorkDir directory that represent that directory that'll be used to perform tasks. "+
			"It is normally a subdirectory (relative) of the 'mountdir' directory.")

	JobCMD.PersistentFlags().StringVarP(&mountDir,
		"mountdir",
		"", "",
		"MountDir directory that represent that directory that'll be mounted or copied to the"+
			" Dagger container")

	JobCMD.PersistentFlags().BoolVarP(&showEnvVars,
		"show-env-vars",
		"", false,
		"Show the environment variables that'll be used in the job.")

	_ = viper.BindPFlag("jobName", JobCMD.PersistentFlags().Lookup("job-name"))
	_ = viper.BindPFlag("dotFiles", JobCMD.PersistentFlags().Lookup("dotfiles"))
	_ = viper.BindPFlag("workDir", JobCMD.PersistentFlags().Lookup("workdir"))
	_ = viper.BindPFlag("mountDir", JobCMD.PersistentFlags().Lookup("mountdir"))
	_ = viper.BindPFlag("showEnvVars", JobCMD.PersistentFlags().Lookup("show-env-vars"))
}

func init() {
	addPersistentFlagsToJobCMD()
	JobCMD.AddCommand(DaggerCMD)
}
