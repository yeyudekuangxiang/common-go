/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package auth

import (
	"github.com/spf13/cobra"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/util"
	"strconv"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Initialize("./config-dev.ini")
		err := cmd.Flags().Parse(args)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		t := cmd.Flag("type").Value.String()
		idStr := cmd.Flag("id").Value.String()
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		token := ""
		switch t {
		case "user":
			token, err = util.CreateToken(auth.User{
				ID:     id,
				Mobile: "",
			})
		case "admin":
			token, err = util.CreateToken(auth.Admin{
				MioAdminID: int(id),
			})
		case "business":
			token, err = util.CreateToken(auth.BusinessUser{
				ID:        id,
				Mobile:    "",
				Uid:       "",
				CreatedAt: model.Time{},
			})
		default:
			cmd.PrintErrln("不支持的类型", t)
			return
		}
		if err != nil {
			cmd.PrintErrln("不支持的类型", t)
			return
		}
		cmd.Println(token)
	},
}

func init() {
	AuthCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringP("type", "t", "user", "生成token的类型 user,admin,business")
	tokenCmd.Flags().Int64P("id", "i", 0, "用户的id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
