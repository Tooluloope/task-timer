package task

import (
	"fmt"

	"github.com/spf13/cobra"
)





var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Logic for creating, updating, and deleting tasks",
	Long: `This command is used to create, update, and delete tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("task called")
	},
}


func init() {
	
}
