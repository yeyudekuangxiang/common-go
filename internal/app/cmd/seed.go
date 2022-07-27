/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/event"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.DB.Transaction(func(tx *gorm.DB) error {
			eventCategory(tx)
			return nil
		})
	},
}

func init() {
	dbCmd.AddCommand(seedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func eventCategory(db *gorm.DB) {
	categoryList := []event.EventCategory{
		{
			EventCategoryId: "79550af260ecf9df2635751c3273b269",
			Title:           "生态环保",
			Active:          true,
			ImageUrl:        "https://resources.miotech.com/static/mp2c/images/event/shouye/sy_sthb.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/eep.png",
			Sort:            2,
		},
		{
			EventCategoryId: "cbddf0af60ecf9f11676bcbd6482736f",
			Title:           "公益善心",
			Active:          true,
			ImageUrl:        "https://resources.miotech.com/static/mp2c/images/event/shouye/sy_rw.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/hc.png",
			Sort:            1,
		},
		{
			EventCategoryId: "79550af260ecfcd4263627ff7c516d0b",
			Title:           "碳减排证书",
			Active:          true,
			ImageUrl:        "https//resources.miotech.com//static/mp2c/images/event/shouye/sy_dtjp.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/lcaer.png",
			Sort:            3,
		},
	}

	for _, category := range categoryList {
		old := event.EventCategory{}
		err := db.Where("event_category_id = ?", category.EventCategoryId).First(&old).Error
		panicOnErr(err)

		if old.ID > 0 {
			category.ID = old.ID
			panicOnErr(db.Save(&category).Error)
		} else {
			panicOnErr(db.Create(&category).Error)
		}
	}
}
func panicOnErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
