/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"mio/internal/pkg/service"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//initialize.Initialize("./config-prod.ini")

		err := cmd.ParseFlags(args)
		if err != nil {
			log.Fatal(err)
		}
		topicPath := cmd.Flag("topic").Value.String()
		userPath := cmd.Flag("user").Value.String()
		if topicPath == "" || userPath == "" {
			log.Fatal("参数错误")
		}

		err = service.DefaultTopicService.ImportUser(userPath)
		if err != nil {
			log.Fatal(err)
		}

		err = service.DefaultTopicService.ImportTopic(topicPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	topicCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("topic", "t", "", "topic file path")
	importCmd.Flags().StringP("user", "u", "", "user file path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}