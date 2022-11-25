/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package point

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"sync"
	"time"
)

var pointSql = `SELECT ID,
	point.openid,
	point.balance,
	COALESCE ( log.aa, 0 ) new_point 
FROM
	point
	LEFT JOIN ( SELECT openid, SUM ( VALUE ) aa FROM point_transaction WHERE create_time >= '2022-01-01 00:00:00' AND VALUE > 0  GROUP BY openid ) log ON point.openid = log.openid`

var ant, _ = ants.NewPool(100)

// DealExpireCmd represents the DealExpireCmd command
var DealExpireCmd = &cobra.Command{
	Use:   "dealexpire",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Initialize("./config-dev.ini")
		list := make([]PointExpire, 0)
		wg := sync.WaitGroup{}
		fmt.Println("haha")
		app.DB.Table(fmt.Sprintf("(%s) point_expire", pointSql)).FindInBatches(&list, 100, func(tx *gorm.DB, batch int) error {
			fmt.Println("处理批次", batch, "处理条数", len(list))
			for i, item := range list {
				i := i
				item := item
				wg.Add(1)
				err := ant.Submit(func() {
					defer wg.Done()
					dealOne(i, item)
				})
				if err != nil {
					log.Printf("提交任务失败 %v %+v\n", i, item)
				}
			}
			return nil
		})
		wg.Wait()
	},
}

func dealOne(i int, item PointExpire) {
	log.Println(i, item.Id, item.Balance, item.NewPoint)
	//全部为2022年之前获取的积分 全部于2022年12月月底过期
	if item.NewPoint == 0 {
		err := dealOld(item.Openid, item.Balance)
		if err != nil {
			log.Printf("创建可用积分异常1 %+v %v\n", item, err)
		}
		return
	}

	//全部为2022年内获取的积分 按时间倒序 group by month 插入 pa 总额为item.Balance
	if item.Balance <= item.NewPoint {
		err := dealNew(item.Openid, item.Balance)
		if err != nil {
			log.Printf("创建可用积分异常2 %+v %v\n", item, err)
		}
		return
	}

	//既有22年前的积分 也有22年的积分
	oldPoint := item.Balance - item.NewPoint
	err := dealOld(item.Openid, oldPoint)
	if err != nil {
		log.Printf("创建21年可用积分异常1 %+v %v\n", item, err)
		return
	}

	//全部为2022年内获取的积分 按时间倒序 group by month 插入 pa 总额为item.NewPoint
	err = dealNew(item.Openid, item.NewPoint)
	if err != nil {
		log.Printf("创建22年可用积分异常2 %+v %v\n", item, err)
	}
}

type UserMonthPoint struct {
	Openid string
	Tm     string
	Point  int64
}

var userMonthSql = `SELECT
	* 
FROM
	(
	SELECT
		openid,
		to_char( create_time, 'yyyymm' ) tm,
		SUM ( VALUE ) point 
	FROM
		point_transaction 
	WHERE
		
	VALUE
		> 0 
		AND create_time >= '2022-01-01 00:00:00' 
		AND openid = ? 
	GROUP BY
		openid,
		to_char( create_time, 'yyyymm' ) 
	) A 
ORDER BY
	tm DESC`

func dealNew(openId string, newPoint int64) error {
	if newPoint == 0 {
		return nil
	}
	list := make([]UserMonthPoint, 0)
	err := app.DB.Raw(userMonthSql, openId).Scan(&list).Error
	if err != nil {
		fmt.Println("异常了", err)
		return err
	}

	return app.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			if newPoint == 0 {
				return nil
			}
			if newPoint < 0 {
				return errors.New("处理积分异常")
			}
			point := item.Point
			if newPoint < item.Point {
				point = newPoint
			}
			newPoint -= point

			tm, err := time.Parse("200601", item.Tm)
			if err != nil {
				return err
			}
			pa := PointAvailable{
				Openid:         openId,
				AvailablePoint: point,
				IsExpired:      false,
				TimePoint:      item.Tm,
				ExpireMonth:    tm.AddDate(1, 0, 0).Format("200601"),
			}
			err = tx.Create(&pa).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func dealOld(openId string, oldPoint int64) error {
	if oldPoint == 0 {
		return nil
	}
	pa := PointAvailable{
		Openid:         openId,
		AvailablePoint: oldPoint,
		IsExpired:      false,
		TimePoint:      "202112",
		ExpireMonth:    "202212",
	}
	err := app.DB.Create(&pa).Error
	if err != nil {
		log.Printf("创建可用积分异常 %+v %v\n", pa, err)
	}
	return err
}

type PointAvailable struct {
	PointAvailableId int64        `gorm:"primaryKey;autoIncrement;column:point_available_id"` // 可用积分记录id 数据库自增的
	Openid           string       `gorm:"column:openid"`                                      // 用户openid
	AvailablePoint   int64        `gorm:"column:available_point"`                             // 可用积分余额
	IsExpired        bool         `gorm:"column:is_expired"`                                  // 是否已经失效
	ActualExpireTime sql.NullTime `gorm:"column:actual_expire_time"`                          // 实际失效时间
	CreatedAt        time.Time    `gorm:"column:created_at"`                                  // 创建时间
	UpdatedAt        time.Time    `gorm:"column:updated_at"`                                  // 更新时间
	TimePoint        string       `gorm:"column:time_point"`                                  // 获取积分当月 202206
	ExpireMonth      string       `gorm:"column:expire_month"`                                // 积分失效月份 202306 此月最后一天失效
}
type PointExpire struct {
	Id       int64
	Balance  int64
	NewPoint int64
	Openid   string
}

func init() {
	PointCmd.AddCommand(DealExpireCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
