package task

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tooluloope/task-timer/pkg/storage"
)

var taskName string
var tags []string

var createCmd = &cobra.Command{
	Use:   "create [taskName]",
	Short: "Creates a new task",
	Long:  `This command creates a new task entry with the given name. Example: task-timer task create "Project Design"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName = args[0]
		createTask(taskName)
	},
}

func createTask(taskName string) {
	currentTime := time.Now()

	task := storage.Task{Name: taskName, Tags: tags, Status: storage.Created.String(), UpdatedAt: currentTime, CreatedAt: currentTime}

	id, err := storage.Data.SaveTask(task)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Task created with ID: %s\n Use the ID to start your new task\n", id)
}

func init() {
	TaskCmd.AddCommand(createCmd)
	createCmd.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "Tags assigned to this task")
}
