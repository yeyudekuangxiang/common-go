/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package topic

import (
	"github.com/spf13/cobra"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/service/community"
)

// importCmd represents the import command
var hotTopicCmd = &cobra.Command{
	Use:   "hot-topic",
	Short: "hot",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Initialize("./config.ini")

		//err := cmd.ParseFlags(args)
		//if err != nil {
		//	log.Fatalf("args not found: %s", err.Error())
		//}
		//
		//configPath := cmd.Flag("config").Value.String()
		//initialize.Initialize(configPath)
		//
		//if configPath == "" {
		//	log.Fatal("config path is required")
		//}
		community.NewTopicService(context.NewMioContext()).SetWeekTopic()
	},
}

func init() {
	TopicCmd.AddCommand(hotTopicCmd)

	//hotTopicCmd.Flags().StringP("config", "-c", "", "config file relative path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
