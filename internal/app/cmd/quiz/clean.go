/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package quiz

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/service/quiz"
	"mio/internal/pkg/util/timeutils"
	"strings"
	"time"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
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

		type OpenID struct {
			Openid string
			Tdate  string
		}
		list := make([]OpenID, 0)
		err = app.DB.Raw("SELECT record.openid,record.tdate FROM ( SELECT openid, to_char( answer_time, 'YYYY-MM-DD' ) tdate, COUNT ( * ) COUNT  FROM quiz_single_record  WHERE answer_time >= '2022-10-10' and answer_time < '2022-10-14' GROUP BY openid, to_char( answer_time, 'YYYY-MM-DD' )  HAVING COUNT ( * ) >= 4  ) record LEFT JOIN ( SELECT openid, to_char( answer_time, 'YYYY-MM-DD' ) tdate FROM quiz_daily_result WHERE answer_time >= '2022-10-10' and answer_time < '2022-10-14' ) RESULT ON record.openid = RESULT.openid  AND record.tdate = RESULT.tdate WHERE RESULT.openid IS NULL").Find(&list).Error
		dealList := make([]string, 0)
		today := time.Now().Format("2006-01-02")
		for _, item := range list {
			if item.Tdate == today {
				dealList = append(dealList, fmt.Sprintf("跳过今天 %s %s", item.Openid, item.Tdate))
				continue
			}

			day, err := time.ParseInLocation("2006-01-02", item.Tdate, time.Local)
			if err != nil {
				dealList = append(dealList, fmt.Sprintf("日期异常 %s %s", item.Openid, item.Tdate))
				continue
			}

			todayResult, err := quiz.DefaultQuizDailyResultService.CompleteTodayQuiz(item.Openid, timeutils.ToTime(day))
			if err != nil {
				dealList = append(dealList, fmt.Sprintf("提交失败 %s %s", item.Openid, item.Tdate))
				continue
			}
			err = quiz.DefaultQuizSummaryService.UpdateTodaySummary(quiz.UpdateSummaryParam{
				OpenId:           item.Openid,
				TodayCorrectNum:  todayResult.CorrectNum,
				TodayAnsweredNum: todayResult.IncorrectNum + todayResult.CorrectNum,
			})
			if err != nil {
				dealList = append(dealList, fmt.Sprintf("统计失败 %s %s", item.Openid, item.Tdate))
				continue
			}
			p, err := quiz.DefaultQuizService.SendAnswerPoint(item.Openid, todayResult.CorrectNum)
			if err != nil {
				dealList = append(dealList, fmt.Sprintf("积分失败 %s %s", item.Openid, item.Tdate))
				continue
			} else {
				dealList = append(dealList, fmt.Sprintf("成功 %s %s %d", item.Openid, item.Tdate, p))
			}
		}
		fmt.Println(strings.Join(dealList, ","))
		ioutil.WriteFile("./result22.txt", []byte(strings.Join(dealList, "\n")), 777)
	},
}

func init() {
	cleanCmd.Flags().StringP("config", "c", "./config.ini", "config file")

	QuizCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
