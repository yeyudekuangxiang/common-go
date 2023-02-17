/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package user

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxapp"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"time"
)

// validatorCmd represents the validator command
var validatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validator called")
		initialize.Initialize("./config-prod.ini")
		list, err := service.DefaultUserService.GetUserListBy(repository.GetUserListBy{
			StartTime: time.Date(2022, 5, 7, 0, 0, 0, 0, time.Local),
			EndTime:   time.Date(2022, 5, 7, 23, 59, 59, 0, time.Local),
		})
		if err != nil {
			log.Fatal(err)
		}

		if len(list) > 50000 {
			log.Fatal("数量太多")
		}
		fmt.Println("数量", len(list))

		for i, user := range list {
			fmt.Println(i)
			rest, err := app.Weapp.GetUserRiskRank(wxapp.UserRiskRankParam{
				AppId:    config.Config.Weapp.AppId,
				OpenId:   user.OpenId,
				Scene:    0,
				ClientIp: "192.168.0.1",
			})
			if rest.ErrCode != 0 && rest.ErrCode != 61010 {
				app.Logger.Errorf("Validator Cmd error: %+v", rest)
			}
			log.Println(rest, err)
		}
	},
}

func init() {
	UserCmd.AddCommand(validatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
