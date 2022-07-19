/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/event"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
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
		app.DB.AutoMigrate(
			&event.Event{},
			&event.EventCategory{},
			&entity.UploadLog{},
			&entity.UploadScene{},
			&entity.Badge{},
			&entity.Banner{},
		)
	},
}

func init() {
	dbCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrateCmd.Flags().StringP("config", "c", "./config.ini", "config file")
}
