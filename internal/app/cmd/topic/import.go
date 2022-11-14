///*/*
//Copyright © 2022 NAME HERE <EMAIL ADDRESS>
//
//*/
package topic

//
//import (
//	"github.com/spf13/cobra"
//	"log"
//	"mio/internal/pkg/core/context"
//	"mio/internal/pkg/service/kumiaoCommunity"
//	"strconv"
//)
//
//// importCmd represents the import command
//var ImportCmd = &cobra.Command{
//	Use:   "import",
//	Short: "A brief description of your command",
//	Long: `A longer description that spans multiple lines and likely contains examples
//and usage of using your command. For example:
//
//Cobra is a CLI library for Go that empowers applications.
//This application is a tool to generate the needed files
//to quickly create a Cobra application.`,
//	Run: func(cmd *cobra.Command, args []string) {
//		//initialize.Initialize("./config-prod.ini")
//
//		err := cmd.ParseFlags(args)
//		if err != nil {
//			log.Fatal(err)
//		}
//		topicPath := cmd.Flag("topic").Value.String()
//		userPath := cmd.Flag("user").Value.String()
//		if topicPath == "" {
//			log.Fatal("参数错误")
//		}
//
//		importIdStr := cmd.Flag("importId").Value.String()
//		importId, err := strconv.Atoi(importIdStr)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		if userPath != "" {
//			err = kumiaoCommunity.NewTopicService(context.NewMioContext()).ImportUser(userPath)
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//
//		err = kumiaoCommunity.DefaultTopicService.ImportTopic(topicPath, importId)
//		if err != nil {
//			log.Fatal(err)
//		}
//	},
//}
//
//func init() {
//	TopicCmd.AddCommand(importCmd)
//
//	importCmd.Flags().StringP("topic", "t", "", "topic file path")
//	importCmd.Flags().StringP("user", "u", "", "user file path")
//	importCmd.Flags().StringP("importId", "i", "", "base importId")
//
//	// Here you will define your flags and configuration settings.
//
//	// Cobra supports Persistent Flags which will work for this command
//	// and all subcommands, e.g.:
//	// importCmd.PersistentFlags().String("foo", "", "A help for foo")
//
//	// Cobra supports local flags which will only run when this command
//	// is called directly, e.g.:
//	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}
