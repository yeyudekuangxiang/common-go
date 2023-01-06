package activity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/timeutils"
	"time"
)

var DefaultReportController = ReportController{}

type ReportController struct {
}

var carbonSql = `select type,sum(value) from carbon_transaction where   created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' and  openid = ?  group by "type"`

var badgeSql = `SELECT event.event_category_id  FROM badge LEFT JOIN event on  event.product_item_id = badge.product_item_id  where  create_time >= '2022-01-01 00:00:01' and create_time <= '2023-01-01 00:00:01' and  openid = ?  
`
var badgeSqlV2 = `
select dd.*,event_category.title from (
select event_category_id,count(*) as cateTotal  from (
SELECT badge.product_item_id,event.event_category_id  FROM badge LEFT JOIN event on  event.product_item_id = badge.product_item_id  where  create_time >= '2022-01-01 00:00:01' and create_time <= '2023-01-01 00:00:01' and  openid = 'oy_BA5EsE0mPQvll8eAqPCkBvI8Q'  

) ll  GROUP BY event_category_id order by cateTotal asc 
) dd  LEFT JOIN event_category on event_category.event_category_id = dd.event_category_id 
`

type carbonList struct {
	Type string
	Sum  float64
}

type badgeList struct {
	EventCategoryId string
}

func (ctr ReportController) Index(ctx *gin.Context) (gin.H, error) {
	//当前登陆用户信息
	user := apiutil.GetAuthUser(ctx)
	println(user.ID)

	//用户基本信息
	userPage := make(map[string]string)
	userPage["register_time"] = user.Time.Format("2006-01-02")
	userPage["register_days"] = fmt.Sprint(timeutils.Now().GetDiffDays(time.Now(), user.Time.Time))

	//select type,sum(value) from carbon_transaction where   created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' and  openid = 'oy_BA5IGl1JgkJKbD14wq_-Yorqw' group by "type"

	openId := "oy_BA5IGl1JgkJKbD14wq_-Yorqw"

	//碳量
	carbon := make([]carbonList, 0)
	err := app.DB.Raw(carbonSql, openId).Scan(&carbon).Error
	if err != nil {
		fmt.Println("异常了", err)
	}

	var carbonTotal float64
	var carbonFavourite float64
	var carbonFavouriteType entity.PointTransactionType
	for _, list := range carbon {
		carbonTotal += list.Sum
		if carbonFavourite < list.Sum {
			carbonFavourite = list.Sum
			carbonFavouriteType = entity.PointTransactionType(list.Type)
		}
	}
	carbonTotalDec := decimal.NewFromFloat(carbonTotal)
	carbonFavouriteDec := decimal.NewFromFloat(carbonFavourite)
	overPer := carbonFavouriteDec.Div(carbonTotalDec).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"

	carbonPage := make(map[string]string)
	carbonPage["favourite_type"] = carbonFavouriteType.Text()
	carbonPage["favourite_carbon"] = util.CarbonToRate(carbonFavourite)
	carbonPage["favourite_ratio"] = overPer

	//公益

	badge := make([]badgeList, 0)
	err = app.DB.Raw(badgeSql, openId).Scan(&badge).Error
	if err != nil {
		fmt.Println("异常了", err)
	}
	badgeTotal := 0

	for _, list := range badge {
		switch list.EventCategoryId {
		case "cbddf0af60ecf9f11676bcbd6482736f":

		}
		badgeTotal++
	}

	badgePage := make(map[string]string)
	badgePage["badge_total"] = "6"
	carbonPage["favourite_type"] = "公益善心"

	answerPage := make(map[string]string)
	answerPage["answer_total"] = "100"
	answerPage["correct_rate"] = "50%" //正确率

	orderPage := make(map[string]string)
	orderPage["order_total"] = "100"
	orderPage["most_point_goods"] = "骑行券"

	communityPage := make(map[string]string)
	communityPage["total_topic"] = "3"
	communityPage["favourite_topic"] = "3"
	communityPage["comment_topic"] = "10"

	newModuleExperiencePage := make(map[string]string)
	newModuleExperiencePage["desc"] = "热爱探索绿喵新功能的你\n是“我的碳排”的种子用户\n你的碳排放等级是几级？"

	lastPage := make(map[string]string)
	lastPage["carbon"] = "111g"
	lastPage["badge"] = "5"
	lastPage["topic"] = "3"
	lastPage["identity"] = "最爱骑行的乐活家"

	return gin.H{
		"user_page":                  userPage,
		"carbon_page":                carbonPage,
		"badge_page":                 badgePage,
		"answer_page":                answerPage,
		"community_page":             communityPage,
		"new_module_experience_page": newModuleExperiencePage,
		"last_page":                  lastPage,
	}, nil

}
