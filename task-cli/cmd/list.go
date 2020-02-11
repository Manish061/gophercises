package cmd

import (
	"fmt"
	"gophercises/task-cli/db"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, errFetchingTasks := db.AllTasks()
		if errFetchingTasks != nil {
			fmt.Printf("Something went wrong while fetching tasks: %v\n", errFetchingTasks)
			return
		}
		if len(tasks) < 1 {
			// grinFace := html.UnescapeString("&#" + strconv.Itoa(128513) + ";")
			fmt.Printf("You don't have any tasks %s\n", string(0x1f601))
			return
		}
		fmt.Println("You have the following tasks:")
		fmt.Printf("#\t TaskName\tCompletedOn\n")
		for i, task := range tasks {
			if task.Completed != "" {
				timeC, errParsing := time.Parse("2006-01-02 15:04:05 -0700 MST", task.Completed)
				if errParsing != nil {
					fmt.Printf("Something went wrong while fetching tasks: %v\n", errFetchingTasks)
					return
				}
				fmt.Printf("%d.\t%v\t%d-%02d-%02d\n", i+1, task.Value,
					timeC.Year(), timeC.Month(), timeC.Day())
			} else {
				fmt.Printf("%d.\t%v\t%v\n", i+1, task.Value, task.Completed)
			}
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
