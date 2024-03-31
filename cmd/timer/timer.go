/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package timer

import (
	"fmt"

	"github.com/spf13/cobra"
)

// timerCmd represents the timer command
var TimerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Logic for starting, stopping, pausing, and resuming timers",
	Long: `This command is used to start, stop, pause, and resume timers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("timer called")
		cmd.Help()
	},
}

func init() {
	// rootCmd.AddCommand(timerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
