package timer

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/pkg/storage"
)

var stopCmd = &cobra.Command{
	Use:   "stop [task_ID]",
	Short: "Stops a timer for the specified task",
	Long:  `This command stops a timer for the task specified by <ID> or --name <task_name>.`,
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		taskName, _ := cmd.Flags().GetString("name")

		if (len(args) == 0 && taskName == "") || (len(args) == 1 && taskName != "") {
			return fmt.Errorf("you must specify either a task ID as an argument or a task name using the --name flag, but not both")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 1 {
			taskId = args[0]
			stopTaskByID()
		} else {
			taskName, _ = cmd.Flags().GetString("name")
			stopTaskByName()
		}
	},
}

func stopTaskByID() {
	task, err := storage.Data.GetTaskByID(taskId)
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Data.StopTask(task); err != nil {
		log.Fatal(err)
	}
	fmt.Println(task)
}

func stopTaskByName() {

	tasks, err := storage.Data.GetTasksByName(taskName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(tasks)
}

func init() {
	TimerCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task")
}
