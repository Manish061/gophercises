package cmd

import (
	"fmt"
	"gophercises/task-cli/db"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "list all completed task for today",
	Run: func(cmd *cobra.Command, args []string) {
		completedTasks, errFetchingCompletedTasks := db.GetCompletedTasksForToday()
		if errFetchingCompletedTasks != nil {
			fmt.Printf("Something went wrong while fetching tasks: %v\n", errFetchingCompletedTasks)
			return
		}
		allTasks, errFetchingTasks := db.AllTasks()
		if errFetchingTasks != nil {
			fmt.Printf("Something went wrong while fetching tasks: %v\n", errFetchingTasks)
			return
		}
		if len(allTasks) < 1 && len(completedTasks) < 1 {
			fmt.Printf("You don't have any tasks %s\n", string(0x1f601))
			return
		} else if len(completedTasks) < 1 {
			fmt.Printf("You have not completed any of your task%s\n", string(0x1f625))
			return
		}
		fmt.Println("You have completed the following tasks:")
		for i, task := range completedTasks {
			fmt.Printf("%d %v\n", i+1, task.Value)
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(completedCmd)
}
