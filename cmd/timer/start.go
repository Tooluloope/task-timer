package timer

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/pkg/storage"
)

var taskName string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a timer for the specified task",
	Long:  `This command starts a timer for the task specified by <task_name>. If the task does not exist, it creates a new task entry with the given name. Example: task-timer start "Project Design"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
	},
}

func startTask() {

	task, err := storage.Data.GetTask(taskName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(task)
}

func init() {
	TimerCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task")

	err := startCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}
}
