package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the main command for task-cli
var RootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "Task is a CLI task manager",
}
