package activity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/timeutils"
	"time"
)

var DefaultReportController = ReportController{}

type ReportController struct {
}

func (ctr ReportController) Index(ctx *gin.Context) (gin.H, error) {
	//当前登陆用户信息
	user := apiutil.GetAuthUser(ctx)
	println(user.ID)

	userPage := make(map[string]string)
	userPage["register_time"] = user.Time.Format("2006-01-02")
	userPage["register_days"] = fmt.Sprint(timeutils.Now().GetDiffDays(time.Now(), user.Time.Time))

	carbonPage := make(map[string]string)
	carbonPage["favourite_type"] = "步行"
	carbonPage["favourite_carbon"] = "100kg"
	carbonPage["favourite_ratio"] = "70%"

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
