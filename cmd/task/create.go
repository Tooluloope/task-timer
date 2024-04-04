package task

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/pkg/storage"
)
var taskName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new task entry",
	Long: `This command creates a new task entry with the given name. Example: task-timer task create "Project Design"`,
	Run: func(cmd *cobra.Command, args []string) {
		createTask(taskName)
	},
}

func createTask(name string) {

	task:= storage.Task{Name: name}

	if err := storage.Data.SaveTask(task); err != nil {
		fmt.Println(err)
	}
	
}

func init() {
	TaskCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task")

	
	err:= createCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}

}
