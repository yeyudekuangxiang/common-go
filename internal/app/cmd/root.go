/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"mio/internal/app/cmd/auth"
	"mio/internal/app/cmd/certificate"
	"mio/internal/app/cmd/coupon"
	"mio/internal/app/cmd/db"
	"mio/internal/app/cmd/gitlab"
	"mio/internal/app/cmd/quiz"
	"mio/internal/app/cmd/topic"
	"mio/internal/app/cmd/user"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mp2c",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.command.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringP("config", "c", "./config.ini", "配置文件路径")

	cmds()
}
func cmds() {
	RootCmd.AddCommand(user.UserCmd)
	RootCmd.AddCommand(topic.TopicCmd)
	RootCmd.AddCommand(gitlab.GitlabCmd)
	RootCmd.AddCommand(db.DBCmd)
	RootCmd.AddCommand(certificate.CertificateCmd)
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(quiz.QuizCmd)
	RootCmd.AddCommand(coupon.CouponCmd)
}
