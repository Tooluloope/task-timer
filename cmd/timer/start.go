package timer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a timer for the specified task",
	Long: `This command starts a timer for the task specified by <task_name>. If the task does not exist, it creates a new task entry with the given name. Example: task-timer start "Project Design"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
	},
}

func init() {
}
