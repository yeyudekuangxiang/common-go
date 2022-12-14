/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package db

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	sbusiness "mio/internal/pkg/service/business"
	"time"

	"github.com/spf13/cobra"
)

// businessCmd represents the business command
var businessCmd = &cobra.Command{
	Use:   "business",
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
		app.DB.Transaction(func(tx *gorm.DB) error {
			business(tx)
			return nil
		})
	},
}

func init() {
	seedCmd.AddCommand(businessCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// businessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// businessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	businessCmd.Flags().StringP("config", "c", "./config.ini", "config file")
}
func business(db *gorm.DB) {
	db.Where("1=1").Delete(&ebusiness.User{})
	db.Where("1=1").Delete(&ebusiness.CarbonScene{})
	db.Where("1=1").Delete(&ebusiness.Company{})
	db.Where("1=1").Delete(&ebusiness.CompanyCarbonScene{})
	db.Where("1=1").Delete(&ebusiness.Department{})
	db.Where("1=1").Delete(&ebusiness.CarbonCreditsLog{})
	db.Where("1=1").Delete(&ebusiness.CarbonRank{})
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
			BDepartmentId: 1,
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
			Avatar:        config.Config.OSS.CdnDomain + "/static/mp2c/images/topic/mio-kol/mio-avatar.jpg",
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
					Point:        1,
				},
			}.PointRateSetting(),
			MaxCount: 10,
			Title:    "线上会议",
			Desc:     "线上会议描述",
			Icon:     config.Config.OSS.CdnDomain + "/static/mp2c/business/carbon/icon/meeting.png",
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
					Point:        1,
				},
			}.PointRateSetting(),
			MaxCount: 10,
			Title:    "公共交通",
			Desc:     "公共交通描述",
			Icon:     config.Config.OSS.CdnDomain + "/static/mp2c/business/carbon/icon/transport.png",
		},
		{
			ID:   3,
			Type: ebusiness.CarbonTypeEvCar,
			PointRateSetting: ebusiness.PointRate{
				CarbonCredit: decimal.NewFromInt(1),
				Point:        1,
			}.PointRateSetting(),
			MaxCount: 10,
			Title:    "电车充电",
			Desc:     "电车充电描述",
			Icon:     config.Config.OSS.CdnDomain + "/static/mp2c/business/carbon/icon/evcar.png",
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
					Point:        1,
				},
			}.PointRateSetting(),
			MaxCount: 10,
			Title:    "节水节电",
			Desc:     "节水节电描述",
			Icon:     config.Config.OSS.CdnDomain + "/static/mp2c/business/carbon/icon/water.png",
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
					Point:        1,
				},
			}.PointRateSetting(),
			MaxCount: 10,
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
					Point:        1,
				},
			}.PointRateSetting(),
			MaxCount: 10,
		},
		{
			ID:            3,
			CarbonSceneId: 3,
			BCompanyId:    1,
			Sort:          3,
			Status:        1,
			PointRateSetting: ebusiness.PointRate{
				CarbonCredit: decimal.NewFromInt(1),
				Point:        1,
			}.PointRateSetting(),
			MaxCount: 10,
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
	db.Create([]ebusiness.CarbonCreditsLog{
		{
			ID:            1,
			TransactionId: "mock-transaction-id1",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromInt(100),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
		},
		{
			ID:            2,
			TransactionId: "mock-transaction-id2",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromInt(200),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
		},
		{
			ID:            3,
			TransactionId: "mock-transaction-id3",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromInt(100),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
		},
		{
			ID:            4,
			TransactionId: "mock-transaction-id4",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromInt(200),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
		},
		{
			ID:            5,
			TransactionId: "mock-transaction-id5",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromInt(200),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
		},
		{
			ID:            6,
			TransactionId: "mock-transaction-id6",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromInt(200),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
		},

		{
			ID:            7,
			TransactionId: "mock-transaction-id7",
			BUserId:       2,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromFloat(50.12),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
		},
		{
			ID:            8,
			TransactionId: "mock-transaction-id8",
			BUserId:       2,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromFloat(38.99),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -1)},
		},
		{
			ID:            9,
			TransactionId: "mock-transaction-id9",
			BUserId:       2,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromFloat(99.66),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
		},
		{
			ID:            10,
			TransactionId: "mock-transaction-id10",
			BUserId:       2,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromFloat(88.88),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, 0, -7)},
		},
		{
			ID:            11,
			TransactionId: "mock-transaction-id11",
			BUserId:       2,
			Type:          ebusiness.CarbonTypeOnlineMeeting,
			Value:         decimal.NewFromFloat(33.33).Round(2),
			Info:          "{\"OneCityDuration\":3600000000000,\"manyCityDuration\":3600000000000}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
		},
		{
			ID:            12,
			TransactionId: "mock-transaction-id6",
			BUserId:       3,
			Type:          ebusiness.CarbonTypeEvCar,
			Value:         decimal.NewFromFloat(199.99).Round(2),
			Info:          "{\"Electricity\":34}",
			CreatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
			UpdatedAt:     model.Time{Time: time.Now().AddDate(0, -1, 0)},
		},
	})

	sbusiness.DefaultCarbonRankService.InitCompanyUserRank(1, ebusiness.RankDateTypeDay)
	sbusiness.DefaultCarbonRankService.InitCompanyUserRank(1, ebusiness.RankDateTypeWeek)
	sbusiness.DefaultCarbonRankService.InitCompanyUserRank(1, ebusiness.RankDateTypeMonth)

	sbusiness.DefaultCarbonRankService.InitCompanyDepartmentRank(1, ebusiness.RankDateTypeDay)
	sbusiness.DefaultCarbonRankService.InitCompanyDepartmentRank(1, ebusiness.RankDateTypeWeek)
	sbusiness.DefaultCarbonRankService.InitCompanyDepartmentRank(1, ebusiness.RankDateTypeMonth)
}
