/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package quiz

import (
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"sort"
)

// formatCmd represents the quiz command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Flags().Parse(args)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		conPath := cmd.Flag("config").Value.String()
		initialize.Initialize(conPath)

		list := make([]entity.QuizQuestion, 0)

		err = app.DB.Find(&list).Error
		if err != nil {
			log.Panicln(err)
		}

		for _, ques := range list {

			//ques.Choices, ques.AnswerIndex = randomOptions(ques.Choices, ques.AnswerStatement)
			err := app.DB.Save(&ques).Error
			if err != nil {
				log.Println("更新异常", ques.ID, err)
			}
		}

	},
}

func randomOptions(options []string, right string) ([]string, int8) {
	sort.Slice(options, func(i, j int) bool {
		return rand.Intn(50) >= 25
	})
	var rightIndex int8 = -1
	for i, option := range options {
		if option == right {
			rightIndex = int8(i)
		}
	}
	return options, rightIndex
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitlabCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitlabCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	formatCmd.Flags().StringP("config", "c", "./config.ini", "config file")

	QuizCmd.AddCommand(formatCmd)
}
