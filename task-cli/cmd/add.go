package cmd

import (
	"fmt"
	"gophercises/task-cli/db"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, errAddingTask := db.AddTask(task)
		if errAddingTask != nil {
			fmt.Printf("something went wrong: %v\n", errAddingTask)
			return
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
