package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string
	debug   bool
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "stiletto",
	Long: `stiletto is a cmd-line tool that helps you to run your pipelines in a
containerized environment.`,
	Example: `
  stiletto <command> <subcommand> [flags]
  stiletto job --task-files=task1.yaml,task2.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func addPersistentFlagsToRootCMD() {

	rootCmd.PersistentFlags().BoolVarP(&debug,
		"debug",
		"d", false,
		"Enabled debug mode")

	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config",
		"c", "",
		"config file (default is $HOME/.stiletto.yaml)")

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".stiletto")

		_ = viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	addPersistentFlagsToRootCMD()

	cobra.OnInitialize(initConfig)

	// Add Job JobCMD.
	rootCmd.AddCommand(JobCMD)
}
