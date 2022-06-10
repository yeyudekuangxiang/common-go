/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	ebusiness "mio/internal/pkg/model/entity/business"

	"github.com/spf13/cobra"
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
		err := cmd.Flags().Parse(args)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		conPath := cmd.Flag("config").Value.String()
		initialize.Initialize(conPath)
		business(app.DB)
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
	seedCmd.Flags().StringP("config", "c", "./config.ini", "config file")
}
func business(db *gorm.DB) {
	db.Where("1=1").Delete(&ebusiness.User{})
	db.Where("1=1").Delete(&ebusiness.CarbonScene{})
	db.Where("1=1").Delete(&ebusiness.Company{})
	db.Where("1=1").Delete(&ebusiness.CompanyCarbonScene{})
	db.Where("1=1").Delete(&ebusiness.Department{})
	db.Create([]ebusiness.User{
		{
			ID:            1,
			Uid:           "mock-uid-1",
			BCompanyId:    1,
			BDepartmentId: 1,
			Nickname:      "mock-nickname-1",
			Mobile:        "mock-mobile-1",
			TelephoneCode: "86",
			Realname:      "mock-realname-1",
			Avatar:        "mock-avatar-1",
			Status:        1,
		},
		{
			ID:            2,
			Uid:           "mock-uid-2",
			BCompanyId:    1,
			BDepartmentId: 2,
			Nickname:      "mock-nickname-2",
			Mobile:        "mock-mobile-2",
			TelephoneCode: "86",
			Realname:      "mock-realname-2",
			Avatar:        "mock-avatar-2",
			Status:        1,
		},
		{
			ID:            3,
			Uid:           "mock-uid-3",
			BCompanyId:    1,
			BDepartmentId: 2,
			Nickname:      "greencat",
			Mobile:        "13000000000",
			TelephoneCode: "86",
			Realname:      "绿喵",
			Avatar:        "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/images/topic/mio-kol/mio-avatar.jpg",
			Status:        1,
		},
	})
	db.Create([]ebusiness.CarbonScene{
		{
			ID:   1,
			Type: ebusiness.CarbonTypeOnlineMeeting,
			PointRateSetting: ebusiness.PointRateOnlineMeeting{
				OneCity: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        1,
				},
				ManyCity: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        2,
				},
			}.PointRateSetting(),
			Title: "线上会议",
			Desc:  "线上会议描述",
			Icon:  "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/business/carbon/icon/meeting.png",
		},
		{
			ID:   2,
			Type: ebusiness.CarbonTypePublicTransport,
			PointRateSetting: ebusiness.PointRatePublicTransport{
				Bus: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        1,
				},
				Metro: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        2,
				},
			}.PointRateSetting(),
			Title: "公共交通",
			Desc:  "公共交通描述",
			Icon:  "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/business/carbon/icon/transport.png",
		},
		{
			ID:   3,
			Type: ebusiness.CarbonTypeEvCar,
			PointRateSetting: ebusiness.PointRate{
				CarbonCredit: decimal.NewFromInt(1),
				Point:        2,
			}.PointRateSetting(),
			Title: "电车充电",
			Desc:  "电车充电描述",
			Icon:  "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/business/carbon/icon/evcar.png",
		},
		{
			ID:   4,
			Type: ebusiness.CarbonTypeSaveWaterElectricity,
			PointRateSetting: ebusiness.PointRateSaveWaterElectricity{
				Water: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        1,
				},
				Electricity: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        2,
				},
			}.PointRateSetting(),
			Title: "节水节电",
			Desc:  "节水节电描述",
			Icon:  "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/business/carbon/icon/water.png",
		},
	})
	db.Create(&ebusiness.Company{
		ID:       1,
		Cid:      "mock-cid-1",
		Name:     "妙盈科技",
		Email:    "admin@miotech.com",
		Password: "7c4a8d09ca3762af61e59520943dc26494f8941b",
	})
	db.Create([]ebusiness.CompanyCarbonScene{
		{
			ID:            1,
			CarbonSceneId: 1,
			BCompanyId:    1,
			Sort:          1,
			Status:        1,
			PointRateSetting: ebusiness.PointRateOnlineMeeting{
				OneCity: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        1,
				},
				ManyCity: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        2,
				},
			}.PointRateSetting(),
		},
		{
			ID:            2,
			CarbonSceneId: 2,
			BCompanyId:    1,
			Sort:          2,
			Status:        1,
			PointRateSetting: ebusiness.PointRatePublicTransport{
				Bus: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        1,
				},
				Metro: ebusiness.PointRate{
					CarbonCredit: decimal.NewFromInt(1),
					Point:        2,
				},
			}.PointRateSetting(),
		},
		{
			ID:            3,
			CarbonSceneId: 3,
			BCompanyId:    1,
			Sort:          3,
			Status:        1,
			PointRateSetting: ebusiness.PointRate{
				CarbonCredit: decimal.NewFromInt(1),
				Point:        2,
			}.PointRateSetting(),
		},
	})
	db.Create([]ebusiness.Department{
		{
			Title:      "部门1",
			BCompanyId: 1,
			Pid:        0,
			Icon:       "",
		},
		{
			Title:      "部门2",
			BCompanyId: 1,
			Pid:        0,
			Icon:       "",
		},
	})
}
