package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/cmd/task"
	"github.com/tooluloope/task-timer/cmd/timer"
)


var rootCmd = &cobra.Command{
	Use:   "task-timer",
	Short: "Task Timer is a CLI tool for tracking time spent on tasks.",
	Long: `Task Timer is a CLI tool for tracking time spent on tasks. It allows you to create, update, and delete tasks, as well as start and stop timers for tasks.`,
	
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubCommands() {
	rootCmd.AddCommand(timer.TimerCmd)
	rootCmd.AddCommand(task.TaskCmd)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task-timer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubCommands()

}


