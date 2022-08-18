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
	"time"
)

// autoCmd represents the auto command
var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "release 和 hotfix分支合并到master分支时自动合并到develop分支",
	Long:  `release 和 hotfix分支合并到master分支时自动合并到develop分支`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "需要三个参数")
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
		startTime := time.Now()
		source := args[2]
		retry := 0
		for {
			if time.Now().Sub(startTime) > time.Minute*30 {
				fmt.Fprintf(os.Stderr, "自动合并失败,请手动将%s分支合并到develop分支", source)
				os.Exit(1)
			}
			state, err := system.DefaultGitlabService.MergeState(projectId, mergeRequestId)
			if err != nil && retry < 5 {
				retry++
				fmt.Fprintln(os.Stderr, "获取合并状态失败 5秒后自动重试", retry, err)
				goto next
			}

			switch state {
			case "merged":
				err := system.DefaultGitlabService.MergeBranch(projectId, source, "develop")
				if err != nil {
					fmt.Fprintln(os.Stderr, "合并失败", err)
					os.Exit(1)
				} else {
					fmt.Fprintf(os.Stdout, "%s 已自动合并到 develop\n", source)
					os.Exit(0)
				}
			case "closed":
				fmt.Fprintf(os.Stderr, "%s到master的合并请求已关闭 取消自动合并\n", source)
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "合并状态 %s \n 5秒后自动重试", source)
			}
		next:
			time.Sleep(time.Second * 5)
		}
	},
}

func init() {
	mergeCmd.AddCommand(autoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// autoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
