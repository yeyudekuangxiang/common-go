/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package db

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/event"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// v2Cmd represents the v2 command
var v2Cmd = &cobra.Command{
	Use:   "v2",
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
			v2(tx)
			return nil
		})
	},
}

func init() {
	seedCmd.AddCommand(v2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// v2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// v2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	v2Cmd.Flags().StringP("config", "c", "./config.ini", "config file")
}
func v2(db *gorm.DB) {
	eventCategoryList := []event.EventCategory{
		{
			EventCategoryId: "cbddf0af60ecf9f11676bcbd6482736f",
			Title:           "公益善心",
			Active:          true,
			ImageUrl:        "https://resources.miotech.com/static/mp2c/images/event/shouye/sy_rw.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/hc.png",
			Sort:            1,
		},
		{
			EventCategoryId: "79550af260ecf9df2635751c3273b269",
			Title:           "生态环保",
			Active:          true,
			ImageUrl:        "https://resources.miotech.com/static/mp2c/images/event/shouye/sy_sthb.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/eep.png",
			Sort:            2,
		},
		{
			EventCategoryId: "79550af260ecfcd4263627ff7c516d0b",
			Title:           "碳减排证书",
			Active:          true,
			ImageUrl:        "https://resources.miotech.com/static/mp2c/images/event/shouye/sy_dtjp.png",
			Icon:            "https://resources.miotech.com/static/mp2c/images/event/category/icon/lcaer.png",
			Sort:            3,
		},
	}
	eventList := []event.Event{
		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "b00064a760ecfeac26f5a3f512c493aa",
			EventTemplateType:     "EEP",
			Title:                 "守护栖息地任鸟飞",
			Subtitle:              "守护中国东海濒危水鸟及其栖息地，建立民间候鸟及栖息地保护行动网络、开展候鸟保护行动、推进鸟类保护政策倡导。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shqxdrnf_20220722%20%283%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "b00064a760f400a42850b68e1f783c22",
			ParticipationCount:    113,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "守护栖息地任鸟飞",
			Sort:                  10,
			Tag:                   []string{"候鸟保护", "生态保护", "政策倡导"},
			TemplateSetting:       "{\"recipient\":\"「守护栖息地任鸟飞」\",\"money\":\"1元\",\"desc\":\"守护鸟儿栖息地，共享美好自然！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_shqxdrnf_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "79550af260ecff352636b0b4531bae00",
			EventTemplateType:     "EEP",
			Title:                 "保护海洋你我同行",
			Subtitle:              "修复被侵蚀的海岸线和红树林，防止海洋“荒漠化”，守护300万平方公里蓝色家园。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/bhhynwtx_20220722%20%283%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "cbddf0af60f4017417b047590e2601cf",
			ParticipationCount:    110,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "保护海洋你我同行",
			Sort:                  20,
			Tag:                   []string{"海洋保护", "生态修复", "红树林保育"},
			TemplateSetting:       "{\"recipient\":\"「保护海洋你我同行」\",\"money\":\"1元\",\"desc\":\"修复滨海湿地，守护蓝色家园！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_%20bhhynwtx_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "79550af260ecfde52636644875e06696",
			EventTemplateType:     "EEP",
			Title:                 "一亿棵梭梭",
			Subtitle:              "恢复阿拉善关键生态区域200万亩以梭梭为代表的荒漠植被，防止荒漠蔓延。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yykss_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "b00064a760f3ff8e28507ac777424d10",
			ParticipationCount:    0,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "一亿棵梭梭",
			Sort:                  30,
			Tag:                   []string{"荒漠防治", "生态保护", "经济增收"},
			TemplateSetting:       "{\"recipient\":\"「一亿棵梭梭」\",\"money\":\"1元\",\"desc\":\"每棵拯救10m²，让荒漠不再蔓延！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_yykss_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "cbddf0af60ecfe6c1677b1f5468ee936",
			EventTemplateType:     "EEP",
			Title:                 "留住长江的微笑——江豚",
			Subtitle:              "保护极度濒危的哺乳动物长江江豚，留住长江的微笑。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/lzcjdwx_20220722%20%282%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "79550af260f4003f278f20ea1b87024e",
			ParticipationCount:    103,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "留住长江的微笑——江豚",
			Sort:                  40,
			Tag:                   []string{"江豚保护", "生物多样性", "繁育研究"},
			TemplateSetting:       "{\"recipient\":\"「留住江豚的微笑」\",\"money\":\"1元\",\"desc\":\"每棵拯救10m²，让荒漠不再蔓延！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_lzcjdwx_20220722.png\",\"organization\":\"保护濒危物种，关注生物多样性、共筑永续家园！\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "b00064a760ed001b26f5f61a47ed7c2b",
			EventTemplateType:     "HC",
			Title:                 "远方书声图书角——云南施甸县",
			Subtitle:              "为云南施甸县老麦乡杨柳小学提供“远方书声”图书角，通过配备优质适龄的书籍让孩子们共享阅读、收获知识。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfssynsmx_20220722%20%284%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "79550af260f40281278f9c3e2372ac00",
			ParticipationCount:    102,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "远方书声图书角——云南施甸县",
			Sort:                  50,
			Tag:                   []string{"人文关怀", "教育支持", ""},
			TemplateSetting:       "{\"project\":\"云南施甸县老麦乡杨柳小学\",\"money\":\"1元\",\"desc\":\"让图书启迪边远学校孩子们,共同助力儿童快乐成长\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_yfssynsmx_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "cbddf0af60ed266416807a0b6838a821",
			EventTemplateType:     "HC",
			Title:                 "远方书声图书馆——新疆喀什",
			Subtitle:              "在新疆喀什泽普县阿依库勒乡中心小学开展“远方书声”项目，通过配备优质适龄的书籍让孩子们共享阅读、收获知识。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfshxjks_20220722%20%284%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "28ee4e3e60f403512b79968b11d86c15",
			ParticipationCount:    108,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "远方书声图书馆——新疆喀什",
			Sort:                  60,
			Tag:                   []string{"人文关怀", "教育支持", ""},
			TemplateSetting:       "{\"project\":\"新疆阿依库勒中心小学\",\"money\":\"1元\",\"desc\":\"让图书启迪边远学校孩子们,共同助力儿童快乐成长\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_yfshxjks_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "79550af260ed25d12640112b4b2e41a8",
			EventTemplateType:     "HC",
			Title:                 "远方书声图书角——云南老麦乡幼儿园",
			Subtitle:              "为云南施甸县老麦乡幼儿园提供“远方书声”图书角，通过配备优质适龄的书籍让孩子们共享阅读、收获知识。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfssynlmx_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "cbddf0af60f402f717b0987b79709209",
			ParticipationCount:    109,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "远方书声图书角——云南老麦乡幼儿园",
			Sort:                  70,
			Tag:                   []string{"人文关怀", "教育支持", ""},
			TemplateSetting:       "{\"project\":\"云南施甸县老麦乡幼儿园\",\"money\":\"1元\",\"desc\":\"让图书启迪边远学校孩子们,共同助力儿童快乐成长\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_yfssynlmx_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecfcd4263627ff7c516d0b",
			EventId:               "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
			EventTemplateType:     "LCAER",
			Title:                 "郑州快速公交项目——联合国认证碳减排",
			Subtitle:              "郑州快速公交（BRT）项目利用油电混合动力减少碳排放，并作为认证减排项目可用于抵消温室气体排放量。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/zzksgj_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "5fbaf3bc-b5d8-4710-8fc4-d04245fe1a24",
			ParticipationCount:    168,
			ParticipationTitle:    "已减少{X}kgCO₂",
			ParticipationSubtitle: "郑州快速公交项目——联合国认证碳减排",
			Sort:                  80,
			Tag:                   []string{"环境提升", "公共福祉", "经济发展"},
			TemplateSetting:       "{\"project\":\"郑州电动公交车项目\",\"money\":\"50千克\",\"desc\":\"守护蓝色地球，全球气候中和行动有您的一份力量！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_zzksgj_20220722.png\"}",
		},

		{
			EventCategoryId:       "79550af260ecfcd4263627ff7c516d0b",
			EventId:               "d7c133f20c782a75480e96493aa7de3e",
			EventTemplateType:     "LCAER",
			Title:                 "青海曲格二级水电站项目——联合国认证碳减排",
			Subtitle:              "青海玛沁格曲二级水电站通过零碳排放形式提供清洁能源，并作为认证减排项目可用于抵消温室气体排放量。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/qhqgejsdz_20220722%20%283%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "708ce50c87cd689a28ad524dae37f8e8",
			ParticipationCount:    105,
			ParticipationTitle:    "已减少{X}kgCO₂",
			ParticipationSubtitle: "青海曲格二级水电站项目——联合国认证碳减排",
			Sort:                  90,
			Tag:                   []string{"清洁能源", "岗位提供", "经济发展"},
			TemplateSetting:       "{\"project\":\"青海曲格二级水电站项目\",\"money\":\"30千克\",\"desc\":\"守护蓝色地球，全球气候中和行动有您的一份力量！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_qhqgejsdz_20220722.png\"}",
		},

		{
			EventCategoryId:       "79550af260ecfcd4263627ff7c516d0b",
			EventId:               "3f37194e271cbd177b3830120fa631f9",
			EventTemplateType:     "LCAER",
			Title:                 "四川九节滩水电项目——联合国认证碳减排",
			Subtitle:              "四川九节滩水电项目为当地提供可再生能源电力，并作为认证减排项目可用于抵消温室气体排放量。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/scjjtsd_20220722%20%283%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "54c2561e6a6f2ea090f082dca36135d2",
			ParticipationCount:    108,
			ParticipationTitle:    "已减少{X}kgCO₂",
			ParticipationSubtitle: "四川九节滩水电项目——联合国认证碳减排",
			Sort:                  100,
			Tag:                   []string{"清洁能源", "岗位提供", "技术发展"},
			TemplateSetting:       "{\"project\":\"四川九节滩水电项目\",\"money\":\"50千克\",\"desc\":\"守护蓝色地球，全球气候中和行动有您的一份力量！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_scjjtsd_20220722.png\"}",
		},

		{
			EventCategoryId:       "79550af260ecfcd4263627ff7c516d0b",
			EventId:               "903acf9cc687e85dc9de0bb0f054ad2f",
			EventTemplateType:     "LCAER",
			Title:                 "新疆阿勒泰华宁水电项目——联合国认证碳减排",
			Subtitle:              "新疆阿勒泰华宁水电项目通过水力发电实现净零排放，并作为认证减排项目可用于抵消温室气体排放量。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/xjalthnsd_20220722%20%282%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "6104c497e29c1bc0628cffa67c97e171",
			ParticipationCount:    88,
			ParticipationTitle:    "已减少{X}kgCO₂",
			ParticipationSubtitle: "新疆阿勒泰华宁水电项目——联合国认证碳减排",
			Sort:                  110,
			Tag:                   []string{"清洁能源", "岗位提供", "经济发展"},
			TemplateSetting:       "{\"project\":\"新疆阿勒泰华宁水电项目\",\"money\":\"50千克\",\"desc\":\"守护蓝色地球，全球气候中和行动有您的一份力量！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_xjalthnsd_20220722.png\"}",
		},

		{
			EventCategoryId:       "79550af260ecfcd4263627ff7c516d0b",
			EventId:               "5bd4792be982c6cd64c7b0de67a6ce34",
			EventTemplateType:     "LCAER",
			Title:                 "四川沐川火谷水电站项目——联合国认证碳减排",
			Subtitle:              "四川沐川县火谷水电项目通过水力发电实现净零排放，并作为认证减排项目可用于抵消温室气体排放量。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/scmchgsdz_20220722%20%282%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "ace42fb85a02640ea95aecc1cbc4e078",
			ParticipationCount:    103,
			ParticipationTitle:    "已减少{X}kgCO₂",
			ParticipationSubtitle: "四川沐川火谷水电站项目——联合国认证碳减排",
			Sort:                  120,
			Tag:                   []string{"清洁能源", "岗位提供", "经济发展"},
			TemplateSetting:       "{\"project\":\"四川沐川火谷水电站项目\",\"money\":\"50千克\",\"desc\":\"守护蓝色地球，全球气候中和行动有您的一份力量！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_scmchgsdz_20220722.png\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "f706818efa576bd56629765e840f5e62",
			EventTemplateType:     "EEP",
			Title:                 "守三江水护万物源",
			Subtitle:              "守护三江源，让七亿人水源重回清澈。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/ssjshwwy_20220722%20%282%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "7c5af0e0035248ee3f50bfd3880fc377",
			ParticipationCount:    66,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "守三江水护万物源",
			Sort:                  130,
			Tag:                   []string{"水源保护", "生物多样性", "生态保护"},
			TemplateSetting:       "{\"recipient\":\"「守三江水护万物源」\",\"money\":\"1元\",\"desc\":\"守护三江源，让七亿人水源重回清澈！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_ssjshwwy_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "02d3b3452aa82f364a487a77f152775e",
			EventTemplateType:     "EEP",
			Title:                 "守护勺嘴鹬",
			Subtitle:              "守护湿地，为勺嘴鹬漫漫迁飞路保驾护航！",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shszh_20220722%20%283%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "e05ed9cbfee9e8dddb176d7e33aebdb3",
			ParticipationCount:    1034,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "守护勺嘴鹬",
			Sort:                  140,
			Tag:                   []string{"物种保护", "生物多样性", "湿地保护"},
			TemplateSetting:       "{\"recipient\":\"「守护勺嘴鹬」\",\"money\":\"1元\",\"desc\":\"拯救极危勺嘴鹬，为它们保驾护航！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_shszh_20220722.png\",\"organization\":\"深圳市红树林湿地保护基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "73dd769e56993d84bfdb95f39b22d7de",
			EventTemplateType:     "EEP",
			Title:                 "守护地球之肾",
			Subtitle:              "守护“地球之肾”湿地，让地球远离“肾衰竭”！",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shdqzs_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "ae4096de38e12b8b6847dd0bf40b7be7",
			ParticipationCount:    55,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "守护地球之肾",
			Sort:                  150,
			Tag:                   []string{"生态保护", "自然教育", "湿地保护"},
			TemplateSetting:       "{\"recipient\":\"「守护地球之肾」\",\"money\":\"1元\",\"desc\":\"保护滨海湿地,共筑永续家园！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_shdqzs_20220722.png\",\"organization\":\"深圳市红树林湿地保护基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "6bc0ffb20c2f1032314ee2dc9aa75e64",
			EventTemplateType:     "EEP",
			Title:                 "守护藏地每一个生命",
			Subtitle:              "支持当地野生动物救助，培训三江源地区本地乡镇兽医",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shzdmygsm_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "a236a10054c1efb2e6b489506694d959",
			ParticipationCount:    556,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "守护藏地每一个生命",
			Sort:                  160,
			Tag:                   []string{"野生动物救助", "生物多样性", "兽医培训"},
			TemplateSetting:       "{\"recipient\":\"「守护藏地每一个生命」\",\"money\":\"1元\",\"desc\":\"助力野生动物救助，守护藏地生物多样性！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_shzdmygsm_20220722.png\",\"organization\":\"浙江省微笑明天慈善基金会\"}",
		},

		{
			EventCategoryId:       "79550af260ecf9df2635751c3273b269",
			EventId:               "02af181705e54d1b0bbc392973b9a029",
			EventTemplateType:     "EEP",
			Title:                 "蒙新河狸方舟计划",
			Subtitle:              "建立中国首个专业的河狸救助中心",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/mxhlfzjh_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "ee5e9cc109201da8ffb626a0279490ac",
			ParticipationCount:    22,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "蒙新河狸方舟计划",
			Sort:                  170,
			Tag:                   []string{"生物保护", "生物多样性", "河狸救助"},
			TemplateSetting:       "{\"recipient\":\"「蒙新河狸方舟计划」\",\"money\":\"1元\",\"desc\":\"为濒危河狸筑起生存的方舟！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_%20mxhlfzjh_20220722.png\",\"organization\":\"爱德基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "0f2cacf4f89046681957f1cad276f64b",
			EventTemplateType:     "HC",
			Title:                 "乡村孩子的创新课",
			Subtitle:              "让乡村孩子们在家乡享有属于他们的好教育。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/xchzdcxk_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "5e3e607260a8ee5817451c2c6f5c4262",
			ParticipationCount:    45,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "乡村孩子的创新课",
			Sort:                  180,
			Tag:                   []string{"人文关怀", "教育支持", "乡村教师赋能"},
			TemplateSetting:       "{\"project\":\"乡村孩子的创新课\",\"money\":\"1元\",\"desc\":\"让更多的孩子在家乡享有属于他们的好教育！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_xchzdcxk_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "181aa44452c1d4869d138977502d573b",
			EventTemplateType:     "HC",
			Title:                 "流浪动物温饱计划",
			Subtitle:              "支持流浪动物救助站，让流浪动物得到妥善安置",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/lldwwbjh_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "3f7f1ff6670b2144b192809af3e84c76",
			ParticipationCount:    1033,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "流浪动物温饱计划",
			Sort:                  190,
			Tag:                   []string{"流浪动物保护", "城市文明建设", "可持续发展"},
			TemplateSetting:       "{\"project\":\"流浪动物温饱计划\",\"money\":\"1元\",\"desc\":\"让流浪小动物不再挨饿！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_lldwwbjh_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "1153583a3af4843cb82f8a67dc293a4e",
			EventTemplateType:     "HC",
			Title:                 "长棘海星大作战",
			Subtitle:              "支持潜水员清理长棘海星，保护珊瑚礁生态。",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/cjhxdzz_20220722%20%282%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "81aea0d228f0e7673815b00055948f25",
			ParticipationCount:    2213,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "长棘海星大作战",
			Sort:                  200,
			Tag:                   []string{"人文关怀", "潜水员支持", "生态保护"},
			TemplateSetting:       "{\"project\":\"长棘海星大作战\",\"money\":\"1元\",\"desc\":\"清理长棘海星，拯救珊瑚生态！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_%20cjhxdzz_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},

		{
			EventCategoryId:       "cbddf0af60ecf9f11676bcbd6482736f",
			EventId:               "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
			EventTemplateType:     "HC",
			Title:                 "让动物不再受到伤害",
			Subtitle:              "促进《反虐待动物法》立法工作的推动，捍卫动物福利",
			Active:                true,
			CoverImageUrl:         "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/rdwbzsdsh_20220722%20%281%29.png",
			StartTime:             panicTime("2022-03-01 00:00:00+08"),
			EndTime:               panicTime("2022-12-31 00:00:00+08"),
			ProductItemId:         "8539c7f021f1bee18d7c6c73da6cc8a4",
			ParticipationCount:    44,
			ParticipationTitle:    "已支持{X}次",
			ParticipationSubtitle: "让动物不再受到伤害",
			Sort:                  210,
			Tag:                   []string{"动物福祉", "政策推进", "行业建设"},
			TemplateSetting:       "{\"project\":\"让动物不再受到伤害\",\"money\":\"1元\",\"desc\":\"共同捍卫所有动物的尊严和福利！\",\"image\":\"https://resources.miotech.com/static/mp2c/images/event/gyzs/gyzs_rdwbzsdsh_20220722.png\",\"organization\":\"北京市企业家环保基金会\"}",
		},
	}
	eventDetailList := []event.EventDetail{{
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "中国有鸟类1445种，其中候鸟600余种，约占世界候鸟的20%。全球候鸟迁徙线路主要有九条，其中四条从我国经过，几乎覆盖了整个中国版图。中国作为重要的中转站（繁殖地、停歇地与越冬地），中国候鸟及其栖息地保护对全球候乌保护具有非常重要的地位。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shqxdrnf_20220722%20%282%29.png",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "由于围垦、开发建设、污染物排放、外来物种入侵等原因，在中国适宜鸟类生存的栖息地面积逐年缩减，生态系统受到破坏，生物多样性缩减，这也严重威胁到了在此栖息的鸟儿。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shqxdrnf_20220722%20%281%29.png",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shqxdrnf_20220722%20%285%29.png",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "任鸟飞是守护中国候鸟及其栖息地的一个综合性生态保护项目。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "项目广泛发动社会力量，覆盖填补保护空缺，形成与官方自然保护体系互补的民间保护网络；开展科学研究、同步调查和公民科学等多种形式的保护行动，提高公众对鸟类及其栖息地保护的关注；推动政策完善，共同守护鸟类及其栖息地。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shqxdrnf_20220722%20%284%29.png",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "守护中国候鸟及其栖息地，让鸟儿与人类共享美好自然。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "红树林、珊瑚礁、盐沼等是连接海洋和陆地的重要过渡地带，是海洋生物繁殖育幼的重要场所。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/bhhynwtx_20220722%20%281%29.png",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "由于不合理及过度的围垦填海、不科学的海岸带工程、过度捕捞、污染排放入海等原因，20世纪50年代以来，我国已经损失近60%的滨海湿地，包括73%的红树林，超过66%的海岸带受到严重侵蚀。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/bhhynwtx_20220722%20%282%29.png",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "开展保护行动让海洋重现生机，不仅可以持续提供丰富的渔业资源还可以通过滨海海岸带生态系统为沿海人口提供天然的保护屏障。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/bhhynwtx_20220722%20%284%29.png",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "为了应对海洋“荒漠化“、恢复海洋生境，保护海洋你我同行项目通过联合政府部门及地方中心，开展红树林、盐沼、海草床等不同滨海生态系统进行科学的修复工作；同时向公众宣传滨海湿地保护知识、提升其保护意识。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "保护海洋，你我同行。邀您一起守护我们的蓝色家园！",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "每年新增荒漠化土地为2460平方千米——即，每一秒钟有78平方米的土地被荒漠化侵蚀，相当于一个羽毛球场的大小。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yykss_20220722%20%283%29.png",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "一亿棵梭梭项目于2014年正式发起，SEE基金会联合阿拉善盟政府相关部门、当地牧民、合作社，以及民间环保组织、企业家、公众，搭建多方参与平台，共同致力于用十年的时间在阿拉善关键生态区域种植一亿棵以梭梭为代表的沙生植物，恢复200万亩荒漠植被，从而改善当地生态环境，遏制荒漠化蔓延趋势，借助梭梭的衍生经济价值提升牧民的生活水平。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yykss_20220722%20%284%29.png",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "2014-2020年，一亿棵梭梭项目已在阿拉善关键生态区域累计推广种植以梭梭为代表的沙生植物6558.85万棵（约137.75万亩）。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yykss_20220722%20%282%29.png",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "和我们一起种下梭梭树，让荒漠不再蔓延！",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "长江江豚的数量为1012头左右，而野外大熊猫的数量有1800只左右。也就是说，长江江豚的数量比国宝大熊猫还要少！",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lzcjdwx_20220722%20%284%29.png",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "长江江豚极可能已成为长江里的唯一旗舰物种，其生存和健康状况是长江生命力最重要的指示。",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lzcjdwx_20220722%20%281%29.png",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "江豚保护项目借助目视考察、声学仪器、无人机影像等科学手段监测分析江豚的分布和种群状况，研究人工环境下江豚的生活习性、支持长江江豚繁育研究工作，建设“守护江豚示范学校”等方式，保护极度濒危的长江江豚。",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lzcjdwx_20220722%20%283%29.png",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "共同参与江豚保护，守护江豚的微笑！",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "云南施甸县老麦乡杨柳小学位于县城东北方向35千米处的杨柳村，是一所高寒山区小学。该小学教学服务覆盖杨柳村14个村民小组。杨柳小学办学历史悠久，始建于1947年。",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynsmx_20220722%20%283%29.png",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "学校现有办学教学班6个，在校学生221人。该校现有图书馆藏书有3千余册，但十几年来未更新，目前图书缺口在5000册。",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynsmx_20220722%20%285%29.png",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynsmx_20220722%20%281%29.png",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynsmx_20220722%20%282%29.png",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "新疆喀什泽普县阿依库勒乡中心小学始建于2009年，位于阿依库勒乡阿依丁库勒村5组010号，学校现下设两个教学点，巴什吾斯塘教学点和上翁热特教学点。",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfshxjks%E2%80%94%E2%80%9420220722%20%283%29.png",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "泽普县阿依库勒乡中心幼儿园始建于2011年9月，幼儿园现下设5个村级园，库塔依村幼儿园、巴什吾斯塘村幼儿园、上翁热特村幼儿园、塔勒克其村幼儿园和其纳尔勒克村幼儿园。",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "学校现有教学班42个，在校学生1821人，现有图书10705册，平均每人5本书，图书缺口在3万册左右。",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfshxjks%E2%80%94%E2%80%9420220722%20%281%29.png",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfshxjks%E2%80%94%E2%80%9420220722%20%285%29.png",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfshxjks%E2%80%94%E2%80%9420220722%20%282%29.png",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "云南施甸县老麦乡幼儿园首创建于2001年9月1日，2003年8月15日首次取得《云南省举办幼儿园登记证》。地处太和社区太和街，位于县城东北角，属三县之交，海拔2000米左右，学校服务全乡7个村。",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynlmx%E2%80%94%E2%80%9420220722%20%284%29.png",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "在园7个班级，共262位学生。目前尚无适合学龄前孩子阅读的绘本或儿童期刊，此类图书缺口在1000册左右。",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynlmx%E2%80%94%E2%80%9420220722%20%282%29.png",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/yfssynlmx%E2%80%94%E2%80%9420220722%20%283%29.png",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "该项目为联合国气候变化框架公约（UNFCCC）认证碳减排（CER）项目，项目编号4744。",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "郑州BRT是一个现代化的大众交通系统，由超过3100辆油电混合动力巴士组成。该项目涵盖超过100公里的公交专用车道、综合智能票务系统及实时乘务信息系统。",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/zzksgj20220722%20%283%29.png",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "混合动力公交车与传统公交相比，可减少40-60%的温室气体排放，同时在市中心以零排放方式运行。该项目对气候产生积极影响、减少交通堵塞、并降低事故率。",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/zzksgj20220722%20%282%29.png",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "该项目为联合国气候变化框架公约（UNFCCC）认证碳减排（CER）项目，项目编号7507。",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "青海玛沁格曲二级水电站通过水力发电的方式，以可持续的方式满足青海省日益增长的电力需求，并取代西北电网传统的通过燃煤发电从而达到净零排放。",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/qhqgejsdz_20220722%20%282%29.png",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "该项目每年可减少17.7万吨的二氧化碳，并通过提高供电能力，改善供电质量，减少输电线路损耗，为当地社区带来积极的社会和环境效益，为当地的可持续发展做出贡献。",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/qhqgejsdz_20220722%20%281%29.png",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "该项目为联合国气候变化框架公约（UNFCCC）认证碳减排（CER）项目，项目编号7152。",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "四川达州市达县九节滩水电项目位于石梯镇上游，是一个径流式水电站并利用巴河现有的水力潜能发电。该项目加快中国可再生能源并网技术和市场的商业化进程。",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/scjjtsd_20220722%20%282%29.png",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "水电站装机容量为39兆瓦，年发电量约187.51亿瓦时，每年向电网提供净电量180.95亿瓦时。该项目通过四川省电网与华中电网连接，代替化石燃料发电，从而达到净零排放。",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/scjjtsd_20220722%20%281%29.png",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "该项目为联合国气候变化框架公约（UNFCCC）认证碳减排（CER）项目，项目编号7497。",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "华宁水电项目位于阿勒泰苏木达依尔日克河河流中游，包含河上大坝、厂房和变电站以及2套涡轮发电机组。总装机容量为100兆瓦，每年平均运行时间为3322小时。",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/xjalthnsd_20220722%20%281%29.png",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "项目通过新疆电网与中国西北电网相连，缓解了新疆和中国西北电网的电力短缺问题。本项目通过提高电力供应能力、改善电力质量、减少输电线路损耗以及给当地社区带来积极的社会和环境效益。",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/xjalthnsd_20220722%20%283%29.png",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "该项目为联合国气候变化框架公约（UNFCCC）认证碳减排（CER）项目，项目编号3309。",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "沐川火谷项目位于四川省乐山市沐川县利店镇，包含2台水轮发电机组，每台机组容量为20兆瓦，水电站总装机容量为40兆瓦。项目预计每年运行4982小时。",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/scmchgsdz_20220722%20%281%29.png",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "该项目通过提高电力供应能力、改善电力质量、减少该地区的输电线路损耗，以及给当地社区带来积极的社会和环境效益，为当地可持续发展作出贡献。",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/scmchgsdz_20220722%20%283%29.png",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "三江源地区保留着地球上残存的完整自然食物链，有着许多珍贵且濒危的野生动物，也是世界上的最神秘的大猫——雪豹最完整的栖息地。",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/ssjshwwy_20220722%20%284%29.png",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "随着全球气候变化加剧、人口的迅猛增长以及现代文明的强势冲击，这片水源地受到种种威胁：",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/ssjshwwy_20220722%20%281%29.png",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/ssjshwwy_20220722%20%283%29.png",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "北京市企业家环保基金会生态保护与自然教育议题下的“三江源保护〞项目，自2012年起至今，己带动超过104家当地环保组织参与到保护网络，累计覆盖的保护面积达117,446平方千米。",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "守护中国最独特高原生态系统及七亿人的水源地！",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "“极危”物种勺嘴鹬现存约600只，比大熊猫还要稀少。每年春秋两季，勺嘴鹬会短暂停留在中国东部沿海湿地地区进行休息和觅食。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shszh20220722%20%284%29.png",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "大面积的围海造地、不可持续的水产养殖、外来物种入侵、水环境污染等，已经严重破坏了勺嘴鹬的栖息地——滨海湿地。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shszh20220722%20%281%29.png",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "勺嘴鹬也是滨海湿地的旗舰物种，具有重大保护和象征意义，保护勺嘴鹬就是保护跟勺嘴鹬共同迁徙的5000万只迁徙候鸟及它们的重要栖息地——滨海湿地。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shszh20220722%20%282%29.png",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "拯救极危勺嘴鹬，我们需要你！",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "目前，全球湿地消失的速度是森林的3倍。50年间，在中国，滨海湿地消失了53%，红树林消失了73%，珊瑚礁减少了80%。“肾衰竭”为地球“健康”拉响警报。",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shdqzs20220722%20%284%29.png",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "湿地保护需要全民参与，才能从量变促成质变。",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shdqzs20220722%20%283%29.png",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "自2012年7月成立以来，红树林基金会聚焦滨海湿地，致力于推动以红树林为代表的滨海湿地保护和公众自然教育，组建了涵盖保育、教育、科研、国际交流等方面的专业团队，通过连结政府、企业、个人等多方力量，打造“社会化参与的自然保育和教育模式”。 ",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shdqzs20220722%20%282%29.png",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "为人与湿地储蓄一个更美好的未来！",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "三江源地区是中国生物多样性最丰富的区域之一。三江源地区由于地处偏远，保护基础较为薄弱，对于受伤的野生动物无法实施有效的救助。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shzdmygsm_20220722%20%284%29.png",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "十几年前，藏獒产业骤然兴起，藏獒被过多人为繁殖。几年后，藏獒市场毫无征兆地崩塌，无力经营的养狗场留下无数无人看管的流浪藏獒，而被弃养的藏獒由于体型大，生性又凶猛，是当地人与野生动物的极大威胁。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shzdmygsm_20220722%20%283%29.png",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "亟需通过呼吁当地人领养流浪藏獒，再加上科学绝育的手段能改变野生动物们所面临的严峻处境。本项目将充分对接外部志愿者队伍和本地兽医，帮助管理部门建立兽医救助体系。并通过在当地开展流浪狗节育培训，最终实现帮助当地管理部门建立完善的野生动物救助体系的目标。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/shzdmygsm_20220722%20%282%29.png",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "共同守护藏地的每一个生命！",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "蒙新河狸是国家一级重点保护动物，也被称作动物界的建筑师。在国内它们仅分布于新疆阿勒泰地区的乌伦古河流域，目前仅有 190个家族，600只左右。这190个家族中，只有40个家族在保护区之内，剩余的150个家族都在保护区之外，",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "蒙新河狸不仅是濒危动物，更是整个乌伦古河流域的关键物种，它们为许多其他物种创造了小生境，它们的存活状态也指示着乌伦古河流域的健康状态。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/mxhlfzjh_20220722%20%284%29.png",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "蒙新河狸190个家族，每年都有三百余只小河狸会自立门户，它们需要寻找新的领地，在这个过程中，它们会遇到各种问题。最严重和频发的有两种：",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/mxhlfzjh_20220722%20%283%29.png",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "该项目将在中国建立首个专业的蒙新河狸救助中心，这也将是中国唯一一个河狸救助中心。河狸救助中心可以为不同季节遇到生存问题的河狸提供康复的环境，帮助它们度过困难时期后放归自然。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/mxhlfzjh_20220722%20%282%29.png",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "一同携手，为中国仅有的河狸种群打造载着它们生存希望的方舟。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "城乡教育资源的不平等愈发严峻，这也造成了“寒门再难出贵子”。从全国范围来看，农村中学生占54.63%，城市中学生占45.37%，占中学生总人数54.63%的农村学生，考入北大的仅仅只有16.3%，而占中学生总人数45.37%的城市学生，考入北大的却占83.7%。",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/xchzdcxk_20220722%20%284%29.png",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "017年初，田字格公益与正安县教育局签订委托协议，在兴隆村开启了一场乡村教育创新的探索，创立田字格兴隆实验小学，并最终在实践中根据中国乡村特点及乡村孩子需求创立了一种乡村教育新模式——乡土人本教育。",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/xchzdcxk_20220722%20%283%29.png",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/xchzdcxk_20220722%20%282%29.png",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "推动城乡教育资源差异互助，协同共生，形成“一校带多校”、“多点成面”的县域教育模式，旨在培育一批小而美、小而优的未来学校，培养“走出大山能生存，留在大山能生活”的乡村子弟，为乡村教师赋能，为推动乡村振兴助力。",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "让更多的孩子在家乡享有属于他们的好教育！",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "因社会的物质生活水平日渐提高，家庭养宠的数量也随之上升，但由于养宠素质的参差不齐，直接或间接地导致流浪动物的数量增加，这对于社区居民和谐共处，城市卫生文明都造成了直接的影响。",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lldwwbjh20220722%20%284%29.png",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "因地处偏远，加上自身造血机能不足，很多救助站都出现动物日渐增多，但运营资金越来越少的恶性循环，甚至会出现站内动物的温饱都无法解决的问题。",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lldwwbjh20220722%20%283%29.png",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/lldwwbjh20220722%20%282%29.png",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "项目旨在为全中国各地区的正规流浪动物救助站及流浪动物群护点，筹集用于日常喂养流浪动物的粮，帮助它们熬过饥饿，健康地活下去，坚持至找到家的那一天。",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "你的支持让流浪小动物不再挨饿！",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "从2003年开始，我国珊瑚礁就频现长棘海星暴发。三亚最漂亮的珊瑚礁——亚龙湾，到大小东海附近的珊瑚，再到西沙群岛的珊瑚，几乎被吃光。",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/cjhxdzz_20220722%20%281%29.png",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "2019年，我国长棘海星大暴发，政府斥巨资动员渔民和相关机构，清除掉4.5万余只长棘海星，才让珊瑚开始慢慢恢复，但近期发现，部分区域长棘海星又卷土重来了！",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/cjhxdzz_20220722%20%284%29.png",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/cjhxdzz_20220722%20%283%29.png",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "长棘海星以珊瑚为食，啃食珊瑚虫，在珊瑚礁表面留下一块块代表死亡的白色骨骼印记。一只长棘海星一年平均要吃掉5-13平方米的活体珊瑚。长棘海星身上还布满了毒刺，潜水员需要经过专业的培训去找到它们，再使用专业的工具才能将它们一一打捞上岸晒干。",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "项目通过长棘海星清理联合行动、清理长棘海星挑战赛、长棘海星调查及清理专长培训班等活动，让更多潜水员参与对长棘海星科学和系统化的调查与清理，将海星的种群数量减少到暴发数量以下，以减少珊瑚礁退化的威胁，并提高已经受损的珊瑚礁区域恢复的机会。",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "支持潜水员清理长棘海星，拯救珊瑚生态！",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "由于法律的空白、缺失，而导致虐待动物的行为无法受到惩罚和禁止因此，很多动物保护志愿者一直在呼吁，中国需要加紧出台《反虐待动物法》。",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/rdwbzsdsh_20220722%20%284%29.png",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "在出台这部法律之前，需要有其他的已经出台的法律能够完善一些条款，以保障在过渡时期，能够形成有效的震慑和提醒作用。研究动物的专家指出，任何“使动物遭受无法忍受的痛苦”的行为，都是虐待动物。",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/rdwbzsdsh_20220722%20%283%29.png",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "https://resources.miotech.com/static/mp2c/images/event/gyxq/rdwbzsdsh_20220722%20%282%29.png",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "深圳、珠海已相继通过禁食猫狗肉的立法表决，都于2020年5月1日开始实施。继广州动物园取消动物表演后，青岛森林野生动物世界也宣布永久取消动物表演。这些都是我国为动物争取到的福利，但这些仍十分局限和片面。",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "项目通过邀请相关的专家和媒体，组建了“反虐待动物法”立法促进工作组，争取能够促进这部法律的立法通过，并施行。",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "让我们一起捍卫所有动物的尊严和福利！",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	},
	}
	eventRuleList := []event.EventRule{{
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "1、您的500积分将用于支持任鸟飞项目用于守护东海濒危水鸟栖息地。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "b00064a760ecfeac26f5a3f512c493aa",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "1、您的500积分将用于支持北京市企业家环保基金会(SEE基金会）发起的保护海洋你我同行项目用于修复滨海湿地。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "79550af260ecff352636b0b4531bae00",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "1、您的500积分将用于支持北京市企业家环保基金会( SEE基金会）发起的一亿棵梭梭项目用于荒漠化防治。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "79550af260ecfde52636644875e06696",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "1、您的500积分将用于支持武汉白鱀豚保护基金会发起的江豚保护项目用于保护极度濒危物种长江江豚。",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "cbddf0af60ecfe6c1677b1f5468ee936",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "1、您的500积分将用于支持远方书声图书角项目，并定向用于杨柳小学的图书角建设;",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "2、成功兑换积分后，您将获得绿喵mio电子证书;",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "b00064a760ed001b26f5f61a47ed7c2b",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "1、您的500积分将用于支持上海享物公益基金会发起的远方书声图书馆项目，并定向用于阿依库勒乡中心小学的图书馆建设。",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "2、成功兑换积分后，您将获得绿喵mio电子证书;",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "cbddf0af60ed266416807a0b6838a821",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "1、您的500积分将用于支持远方书声图书角项目，并定向用于老麦乡幼儿园的图书角建设;",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "2、成功兑换积分后，您将获得绿喵mio电子证书;",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "79550af260ed25d12640112b4b2e41a8",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "1.您的500积分将用于支持郑州BRT项目，相当于为地球减少二氧化碳排放30千克;",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "⒉您将获得绿喵mio电子证书;",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "3.该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "8d8f78fc-0b76-4d1d-a88d-c86cf9191449",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "1、您的500积分将用于支持格曲水电站项目，相当于为地球减少二氧化碳排放50千克;",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "2、您将获得绿喵mio电子证书;",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "d7c133f20c782a75480e96493aa7de3e",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "1、您的500积分将用于支持九节滩水电站项目，相当于为地球减少二氧化碳排放50千克;",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "2、您将获得绿喵mio电子证书;",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "3f37194e271cbd177b3830120fa631f9",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "1、您的500积分将用于支持阿勒泰华宁水电项目，相当于为地球减少二氧化碳排放50千克;",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "2、您将获得绿喵mio电子证书;",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "903acf9cc687e85dc9de0bb0f054ad2f",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "1、您的500积分将用于支持火谷水电站项目，相当于为地球减少二氧化碳排放50千克;",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "2、您将获得绿喵mio电子证书;",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "3、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "5bd4792be982c6cd64c7b0de67a6ce34",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "1、您的500积分将用于支持北京市企业家环保基金会(SEE基金会）发起的三江源保护项目用于守护三江源地区。",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "f706818efa576bd56629765e840f5e62",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "1、您的500积分将用于支持深圳市红树林湿地保护基金会发起的守护勺嘴鹬项目用于勺嘴鹬保育行动。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "02d3b3452aa82f364a487a77f152775e",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "1、您的500积分将用于支持深圳市红树林湿地保护基金会发起的守护地球之肾项目用于保护滨海湿地。",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "73dd769e56993d84bfdb95f39b22d7de",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "1、您的500积分将用于支持浙江省微笑明天慈善基金会发起的守护藏地每一个生命项目用于当地野生动物救治。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "6bc0ffb20c2f1032314ee2dc9aa75e64",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "1、您的500积分将用于支持爱德基金会发起的蒙新河狸方舟计划项目用于监测河狸生存情况和提供所需的救助。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "02af181705e54d1b0bbc392973b9a029",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "1、您的500积分将用于支持上海联劝公益基金会发起的乡村孩子的创新课项目，支持贵州正安县乡村教育;",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项;",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "0f2cacf4f89046681957f1cad276f64b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "1、您的500积分将用于支持深圳市社会公益基金会发起的流浪动物温饱计划项目，用于流浪动物保护;",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项;",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "181aa44452c1d4869d138977502d573b",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "1、您的500积分将用于支持海南成美慈善基金会发起的长棘海星大作战项目，支持潜水员进行长棘海星清理;",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项;",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "1153583a3af4843cb82f8a67dc293a4e",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "1、您的500积分将用于支持中国生物多样性保护与绿色发展基金会发起的让动物不再受到伤害项目，用于推进《反虐待动物法》立法工作;",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "2、用户每次成功兑换500积分，绿喵mio向该项目投入1元的支持款项;",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "3、用户成功兑换后，将获得绿喵mio电子证书。(证书将在兑换后的5个工作日内发放，可在我的-我的证书查看)",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "4、该项目属于虚拟权益，兑换后不会寄送实物。",
	}, {
		EventId: "d8a3e5d8fd0dbe72b5918df2bc87dc8c",
		Content: "",
	},
	}
	productList := []entity.ProductItem{{
		ProductItemId:          "b00064a760f400a42850b68e1f783c22",
		Virtual:                true,
		Title:                  "守护栖息地任鸟飞",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shqxdrnf_20220722%20%283%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   10,
	}, {
		ProductItemId:          "cbddf0af60f4017417b047590e2601cf",
		Virtual:                true,
		Title:                  "保护海洋你我同行",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/bhhynwtx_20220722%20%283%29.png",
		RemainingCount:         300,
		SalesCount:             290,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   20,
	}, {
		ProductItemId:          "b00064a760f3ff8e28507ac777424d10",
		Virtual:                true,
		Title:                  "一亿棵梭梭",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yykss_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   30,
	}, {
		ProductItemId:          "79550af260f4003f278f20ea1b87024e",
		Virtual:                true,
		Title:                  "留住长江的微笑——江豚",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/lzcjdwx_20220722%20%282%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   40,
	}, {
		ProductItemId:          "79550af260f40281278f9c3e2372ac00",
		Virtual:                true,
		Title:                  "远方书声图书角——云南施甸县",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfssynsmx_20220722%20%284%29.png",
		RemainingCount:         300,
		SalesCount:             290,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   50,
	}, {
		ProductItemId:          "28ee4e3e60f403512b79968b11d86c15",
		Virtual:                true,
		Title:                  "远方书声图书馆——新疆喀什",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfshxjks_20220722%20%284%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   60,
	}, {
		ProductItemId:          "cbddf0af60f402f717b0987b79709209",
		Virtual:                true,
		Title:                  "远方书声图书角——云南老麦乡幼儿园",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/yfssynlmx_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   70,
	}, {
		ProductItemId:          "5fbaf3bc-b5d8-4710-8fc4-d04245fe1a24",
		Virtual:                true,
		Title:                  "郑州快速公交项目——联合国认证碳减排",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/zzksgj_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             290,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   80,
	}, {
		ProductItemId:          "708ce50c87cd689a28ad524dae37f8e8",
		Virtual:                true,
		Title:                  "青海曲格二级水电站项目——联合国认证碳减排",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/qhqgejsdz_20220722%20%283%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   90,
	}, {
		ProductItemId:          "54c2561e6a6f2ea090f082dca36135d2",
		Virtual:                true,
		Title:                  "四川九节滩水电项目——联合国认证碳减排",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/scjjtsd_20220722%20%283%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   100,
	}, {
		ProductItemId:          "6104c497e29c1bc0628cffa67c97e171",
		Virtual:                true,
		Title:                  "新疆阿勒泰华宁水电项目——联合国认证碳减排",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/xjalthnsd_20220722%20%282%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   110,
	}, {
		ProductItemId:          "ace42fb85a02640ea95aecc1cbc4e078",
		Virtual:                true,
		Title:                  "四川沐川火谷水电站项目——联合国认证碳减排",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/scmchgsdz_20220722%20%282%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   120,
	}, {
		ProductItemId:          "7c5af0e0035248ee3f50bfd3880fc377",
		Virtual:                true,
		Title:                  "守三江水护万物源",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/ssjshwwy_20220722%20%282%29.png",
		RemainingCount:         300,
		SalesCount:             290,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   130,
	}, {
		ProductItemId:          "e05ed9cbfee9e8dddb176d7e33aebdb3",
		Virtual:                true,
		Title:                  "守护勺嘴鹬",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shszh_20220722%20%283%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   140,
	}, {
		ProductItemId:          "ae4096de38e12b8b6847dd0bf40b7be7",
		Virtual:                true,
		Title:                  "守护地球之肾",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shdqzs_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   150,
	}, {
		ProductItemId:          "a236a10054c1efb2e6b489506694d959",
		Virtual:                true,
		Title:                  "守护藏地每一个生命",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/shzdmygsm_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   160,
	}, {
		ProductItemId:          "ee5e9cc109201da8ffb626a0279490ac",
		Virtual:                true,
		Title:                  "蒙新河狸方舟计划",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/mxhlfzjh_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   170,
	}, {
		ProductItemId:          "5e3e607260a8ee5817451c2c6f5c4262",
		Virtual:                true,
		Title:                  "乡村孩子的创新课",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/xchzdcxk_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   180,
	}, {
		ProductItemId:          "3f7f1ff6670b2144b192809af3e84c76",
		Virtual:                true,
		Title:                  "流浪动物温饱计划",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/lldwwbjh_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   190,
	}, {
		ProductItemId:          "81aea0d228f0e7673815b00055948f25",
		Virtual:                true,
		Title:                  "长棘海星大作战",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/cjhxdzz_20220722%20%282%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   200,
	}, {
		ProductItemId:          "8539c7f021f1bee18d7c6c73da6cc8a4",
		Virtual:                true,
		Title:                  "让动物不再受到伤害",
		Cost:                   500,
		ImageUrl:               "https://resources.miotech.com/static/mp2c/images/event/gyxmzt/rdwbzsdsh_20220722%20%281%29.png",
		RemainingCount:         300,
		SalesCount:             0,
		Active:                 true,
		ProductItemReferenceId: nil,
		Sort:                   210,
	},
	}

	for _, ec := range eventCategoryList {
		old := event.EventCategory{}
		err := db.Where("event_category_id = ?", ec.EventCategoryId).Delete(&old).Error
		panicOnErr(err)
		panicOnErr(db.Create(&ec).Error)
	}
	for _, ev := range eventList {
		old := event.Event{}
		err := db.Where("event_id = ?", ev.EventId).Delete(&old).Error
		panicOnErr(err)
		panicOnErr(db.Create(&ev).Error)
	}

	edMap := make(map[string]string)
	for _, eventDetail := range eventDetailList {
		old := event.EventDetail{}
		if edMap[eventDetail.EventId] == "" {
			err := db.Where("event_id = ?", eventDetail.EventId).Delete(&old).Error
			panicOnErr(err)
		}
		edMap[eventDetail.EventId] = eventDetail.EventId
		panicOnErr(db.Create(&eventDetail).Error)
	}

	erMap := make(map[string]string)
	for _, er := range eventRuleList {
		old := event.EventRule{}
		if erMap[er.EventId] == "" {
			err := db.Where("event_id = ?", er.EventId).Delete(&old).Error
			panicOnErr(err)
		}
		erMap[er.EventId] = er.EventId
		panicOnErr(db.Create(&er).Error)
	}

	for _, p := range productList {
		old := entity.ProductItem{}
		err := db.Where("product_item_id = ?", p.ProductItemId).Delete(&old).Error
		panicOnErr(err)
		panicOnErr(db.Create(&p).Error)
	}

}

func panicOnErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
func v2File() {
	file, err := excelize.OpenFile("event.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if file.SheetCount == 0 {
		log.Fatal("没有数据")
	}

	rows, err := file.GetRows(file.GetSheetList()[2])
	if err != nil {
		log.Fatal(err)
	}

	ev := strings.Builder{}
	// a b c d e f g h i j k  l  m  n  o  p  q  r  s  t
	// 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19

	ev.WriteString("package cmd \n" +
		"func data(){\n eventList:=[]event.Event{")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		ev.WriteString(fmt.Sprintf(`
{
			EventCategoryId:       "%s",
			EventId:               "%s",
			EventTemplateType:     "%s",
			Title:                 "%s",
			Subtitle:              "%s",
			Active:                %v,
			CoverImageUrl:         "%s",
			StartTime:             panicTime("%s"),
			EndTime:               panicTime("%s"),
			ProductItemId:         "%s",
			ParticipationCount:    %s,
			ParticipationTitle:    "%s",
			ParticipationSubtitle: "%s",
			Sort:                  %s,
			Tag:                   []string{"%s"},
			TemplateSetting:       "%s",
		},
`, row[1], row[2], row[14], row[3], row[4], row[5] == "t", row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], strings.ReplaceAll(row[15], ",", `","`), se(row[16])))
	}
	ev.WriteString("\n}\n")

	rows, err = file.GetRows(file.GetSheetList()[3])
	if err != nil {
		log.Fatal(err)
	}

	ev.WriteString("eventDetailList:=[]event.EventDetail{")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		content := ""
		if len(row) > 2 {
			content = row[2]
		}
		ev.WriteString(fmt.Sprintf(`{
			EventId: "%s",
			Content: "%s",
		},`, row[1], strings.Trim(content, "\n")))
	}
	ev.WriteString("\n}\n")

	rows, err = file.GetRows(file.GetSheetList()[4])
	if err != nil {
		log.Fatal(err)
	}

	ev.WriteString("eventRuleList:=[]event.EventRule{")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		content := ""
		if len(row) > 2 {
			content = row[2]
		}
		ev.WriteString(fmt.Sprintf(`{
			EventId: "%s",
			Content: "%s",
		},`, row[1], content))
	}
	ev.WriteString("\n}\n")

	rows, err = file.GetRows(file.GetSheetList()[5])
	if err != nil {
		log.Fatal(err)
	}
	ev.WriteString("productList:=[]entity.ProductItem{")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		ev.WriteString(fmt.Sprintf(`{
			ProductItemId:          "%s",
			Virtual:                true,
			Title:                  "%s",
			Cost:                   %s,
			ImageUrl:               "%s",
			RemainingCount:         %s,
			SalesCount:             %s,
			Active:                 true,
			ProductItemReferenceId: nil,
			Sort:                   %s,
		},`, row[1], row[2], row[3], row[8], row[5], row[6], row[9]))
	}
	ev.WriteString("\n}}")

	ioutil.WriteFile("internal/app/cmd/event.go", []byte(ev.String()), 0777)
}

func panicTime(str string) model.Time {
	t, err := time.Parse("2006-01-02 15:04:05", str[:19])
	if err != nil {
		panic(err)
	}
	return model.Time{Time: t}
}
func se(str string) string {
	str = strings.ReplaceAll(str, "\"", "\\\"")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}
