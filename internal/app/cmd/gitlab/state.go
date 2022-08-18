/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package gitlab

import (
	"fmt"
	"github.com/spf13/cobra"
	"mio/internal/pkg/service/system"
	"os"
	"strconv"
)

// stateCmd represents the status command
var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 1 {
			fmt.Fprintln(os.Stderr, "需要两个参数")
			return
		}
		projectId, err := strconv.Atoi(args[0])
		if len(args) <= 1 {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		mergeRequestId, err := strconv.Atoi(args[1])
		if len(args) <= 1 {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		status, err := system.DefaultGitlabService.MergeState(projectId, mergeRequestId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stdout, status)
	},
}

func init() {
	stateCmd.PersistentFlags()
	mergeCmd.AddCommand(stateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
