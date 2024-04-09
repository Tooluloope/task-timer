package task

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/pkg/storage"
)

var taskName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new task entry",
	Long:  `This command creates a new task entry with the given name. Example: task-timer task create "Project Design"`,
	Run: func(cmd *cobra.Command, args []string) {
		createTask()
	},
}

func createTask() {

	currentTime := time.Now()

	task := storage.Task{Name: taskName, UpdatedAt: currentTime, CreatedAt: currentTime}

	id, err := storage.Data.SaveTask(task)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Task created with ID: %s\n Use the ID to start your new task\n", id)

}

func init() {
	TaskCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&taskName, "name", "n", "", "Name of the task")

	err := createCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}

}
