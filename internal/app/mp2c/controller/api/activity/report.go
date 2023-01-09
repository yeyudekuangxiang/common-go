package activity

import (
	"encoding/json"
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
SELECT badge.product_item_id,event.event_category_id  FROM badge LEFT JOIN event on  event.product_item_id = badge.product_item_id  where  create_time >= '2022-01-01 00:00:01' and create_time <= '2023-01-01 00:00:01' and  openid = ?  

) ll  GROUP BY event_category_id order by cateTotal asc 
) dd  LEFT JOIN event_category on event_category.event_category_id = dd.event_category_id  order by catetotal desc`

var answerSql = `select sum(correct_num) as correct_num,sum(incorrect_num) as incorrect_num  from quiz_daily_result where openid = ? and  answer_time >= '2022-01-01 00:00:01' and answer_time <= '2023-01-01 00:00:01' 
`

var orderSql = `select count(*) as total    from duiba_order where   user_id = ? and   source = '普通兑换' and  total_credits != 0 and  create_time >= 1640966401000  and create_time <= 1672502401000 
`
var orderV2Sql = ` select order_item_list   from duiba_order where   user_id = ? and   source = '普通兑换' and total_credits != 0  and  create_time >= 1640966401000 and   create_time <= 1672502401000    order by total_credits desc  limit 1 
`
var communityCommentSql = `select  count(*)    from comment_index where   member_id  = ? and   created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`
var communityLikeSql = `select  count(*)    from topic_like where   user_id  = 828463 and   created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

var communityTopicSql = `select  count(*)    from topic where   user_id  = 828463 and   created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

//我的碳排
var questionSql = `select count(*) from question_user where third_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' `

//二手市集-发商品

var secondHandGoodsSql = `select count(*) from carbon_second_hand_commodity where user_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

//二手市集-商品点赞
var secondHandGoodsLikeSql = `select count(*) from carbon_commodity_like where user_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

//二手市集-评论点赞
var secondHandCommentLikeSql = `select count(*) from carbon_comment_like where user_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

//二手市集-评论

var secondHandCommentSql = `select count(*) from carbon_comment_index where member_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

//二手市集-订单
var secondHandOrderSql = `select count(*) from second_order where consume_user_open_id = ?  and  created_at >= '2022-01-01 00:00:01' and created_at <= '2023-01-01 00:00:01' 
`

type carbonList struct {
	Type string
	Sum  float64
}

type badgeList struct {
	EventCategoryId string
	Catetotal       int64
	Title           string
}

type answerList struct {
	CorrectNum   int64
	IncorrectNum int64
}

type orderItem struct {
	Title string
}

func (ctr ReportController) Index(ctx *gin.Context) (gin.H, error) {
	//当前登陆用户信息
	user := apiutil.GetAuthUser(ctx)
	userId := user.ID
	openId := user.OpenId

	userId = 921216
	openId = "oy_BA5IGl1JgkJKbD14wq_-Yorqw"

	//用户基本信息
	userPage := make(map[string]string)
	userPage["register_time"] = user.Time.Format("2006-01-02")
	userPage["register_days"] = fmt.Sprint(timeutils.Now().GetDiffDays(time.Now(), user.Time.Time))

	//碳量
	carbon := make([]carbonList, 0)
	err := app.DB.Raw(carbonSql, openId).Scan(&carbon).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取碳异常", err)
	}

	var carbonTotal float64
	var carbonFavourite float64
	var carbonFavouriteType entity.PointTransactionType
	var overPer string
	for _, list := range carbon {
		carbonTotal += list.Sum
		if carbonFavourite < list.Sum {
			carbonFavourite = list.Sum
			carbonFavouriteType = entity.PointTransactionType(list.Type)
		}
		carbonTotalDec := decimal.NewFromFloat(carbonTotal)
		carbonFavouriteDec := decimal.NewFromFloat(carbonFavourite)
		overPer = carbonFavouriteDec.Div(carbonTotalDec).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"
	}
	carbonPage := make(map[string]string)
	carbonPage["favourite_type"] = carbonFavouriteType.Text()
	carbonPage["favourite_carbon"] = util.CarbonToRate(carbonFavourite)
	carbonPage["favourite_ratio"] = overPer

	/*公益开始*/
	badge := make([]badgeList, 0)
	err = app.DB.Raw(badgeSqlV2, openId).Scan(&badge).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取证书异常", err)
	}
	badgePage := make(map[string]interface{})
	var badgeTotal int64
	if len(badge) != 0 {
		firstBadge := badge[0]
		for _, badgeInfo := range badge {
			badgeTotal = badgeTotal + badgeInfo.Catetotal
		}
		badgePage["badge_total"] = badgeTotal
		badgePage["favourite_type"] = firstBadge.Title
	} else {
		badgePage["no_date"] = true
	}
	/*公益结束*/

	/*答题开始*/
	answerPage := make(map[string]interface{})
	var answerTotal decimal.Decimal
	var correctRate string
	answer := answerList{}
	err = app.DB.Raw(answerSql, openId).Scan(&answer).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取答题数据异常", err)
	}
	correctNum := decimal.NewFromInt(answer.CorrectNum)
	incorrectNum := decimal.NewFromInt(answer.IncorrectNum)
	answerTotal = correctNum.Add(incorrectNum)
	if answerTotal.IsZero() {
		answerPage["no_date"] = true
	} else {
		correctRate = correctNum.Div(answerTotal).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"
		answerPage["answer_total"] = answerTotal
		answerPage["correct_rate"] = correctRate //正确率
	}

	/*答题结束*/

	/*订单数据开始*/
	orderPage := make(map[string]interface{})
	var orderTotal int64
	err = app.DB.Raw(orderSql, userId).Scan(&orderTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取订单数据异常", err)
	}
	orderPage["order_total"] = orderTotal

	var orderItemStr string
	err = app.DB.Raw(orderV2Sql, userId).Scan(&orderItemStr).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:订单商品转换数据异常", err)
	}
	var maxPointGoods string
	if orderItemStr != "" {
		orderItemArr := make([]orderItem, 0)
		err = json.Unmarshal([]byte(orderItemStr), &orderItemArr)
		if err != nil {
			return nil, err
		}
		for _, item := range orderItemArr {
			maxPointGoods = maxPointGoods + item.Title + " "
		}
	}
	orderPage["max_point_goods"] = maxPointGoods

	if orderTotal == 0 {
		orderPage["no_date"] = true
	}
	communityPage := make(map[string]interface{})

	//评论
	var communityCommentTotal int64
	err = app.DB.Raw(communityCommentSql, userId).Scan(&communityCommentTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取酷喵圈评论数据异常", err)
	}
	communityPage["comment_topic"] = communityCommentTotal

	//点赞
	var communityFavouriteTotal int64
	err = app.DB.Raw(communityLikeSql, userId).Scan(&communityFavouriteTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取酷喵圈帖子点赞数据异常", err)
	}
	communityPage["favourite_topic"] = communityFavouriteTotal

	//发帖
	var communityTopicTotal int64
	err = app.DB.Raw(communityTopicSql, userId).Scan(&communityTopicTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取酷喵圈发帖数据异常", err)
	}
	communityPage["total_topic"] = communityTopicTotal

	if communityCommentTotal == 0 && communityFavouriteTotal == 0 && communityTopicTotal == 0 {
		communityPage["no_date"] = true
	}

	//我的碳排
	var questionTotal int64
	err = app.DB.Raw(questionSql, userId).Scan(&questionTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我的碳排数据异常", err)
	}

	var secondHandGoodsTotal int64
	err = app.DB.Raw(secondHandGoodsSql, userId).Scan(&secondHandGoodsTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我二手交易商品数据异常", err)
	}

	var secondHandGoodsLikeTotal int64
	err = app.DB.Raw(secondHandGoodsLikeSql, userId).Scan(&secondHandGoodsLikeTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我二手交易点赞商品数据异常", err)
	}

	var secondHandCommentLikeTotal int64
	err = app.DB.Raw(secondHandCommentLikeSql, userId).Scan(&secondHandCommentLikeTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我二手交易点赞评论数据异常", err)
	}

	var secondHandCommentTotal int64
	err = app.DB.Raw(secondHandCommentSql, userId).Scan(&secondHandCommentTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我二手交易评论数据异常", err)
	}

	var secondHandOrderTotal int64
	err = app.DB.Raw(secondHandOrderSql, openId).Scan(&secondHandOrderTotal).Error
	if err != nil {
		app.Logger.Error("个人减碳成就报告:获取我二手交易买家购买订单数据异常", err)
	}

	experienceHand := false

	if secondHandGoodsTotal != 0 || secondHandGoodsLikeTotal != 0 || secondHandCommentLikeTotal != 0 || secondHandCommentTotal != 0 || secondHandOrderTotal != 0 {
		experienceHand = true
	}

	newModuleExperiencePage := make(map[string]interface{})
	if questionTotal != 0 && !experienceHand {
		newModuleExperiencePage["desc"] = "热爱探索绿喵新功能的你\n是“我的碳排”的种子用户\n你的碳排放等级是几级？"
	} else if questionTotal == 0 && experienceHand {
		newModuleExperiencePage["desc"] = "热爱探索绿喵新功能的你\n是“二手市集”的种子用户\n你成功交换了闲置物品吗？"
	} else if questionTotal != 0 && experienceHand {
		newModuleExperiencePage["desc"] = "热爱探索绿喵新功能的你\n重大功能更新可都没错过\n你最喜欢的新功能是哪些？"
	} else {
		newModuleExperiencePage["no_date"] = true
	}

	/*帖子数据结束*/
	identity := "保留实力的环保潜力股"
	if carbonTotal > 0 {
		if communityTopicTotal > 0 {
			identity = "最爱" + carbonFavouriteType.Text() + "的乐活家"
		} else if badgeTotal > 0 {
			identity = "最爱" + carbonFavouriteType.Text() + "的地球守护者"
		} else if orderTotal > 0 {
			identity = "最爱" + carbonFavouriteType.Text() + "的喵生活实践家"
		}
	}

	lastPage := make(map[string]interface{})
	lastPage["carbon"] = util.CarbonToRate(carbonTotal)
	lastPage["badge"] = badgeTotal
	lastPage["topic"] = communityTopicTotal
	lastPage["identity"] = identity

	return gin.H{
		"user_page":                  userPage,
		"carbon_page":                carbonPage,
		"badge_page":                 badgePage,
		"answer_page":                answerPage,
		"order_page":                 orderPage,
		"community_page":             communityPage,
		"new_module_experience_page": newModuleExperiencePage,
		"last_page":                  lastPage,
	}, nil

}
