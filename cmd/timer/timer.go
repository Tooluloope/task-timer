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

func addSubCommands(){
	TimerCmd.AddCommand(StartCmd)
}

func init() {
	addSubCommands()
	}
