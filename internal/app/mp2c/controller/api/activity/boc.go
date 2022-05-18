package activity

import (
	"github.com/gin-gonic/gin"
	"log"
	activityM "mio/internal/pkg/model/entity/activity"
	activity2 "mio/internal/pkg/service/activity"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultBocController = BocController{}

type BocController struct {
}

var bocEndTime, _ = time.Parse("2006-01-02 15:04:05", "2022-06-30 23:59:59")

func (b BocController) GetRecordList(c *gin.Context) (gin.H, error) {
	form := GetBocRecordListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	if user.ID == 0 {
		return gin.H{
			"list":  make([]interface{}, 0),
			"total": 0,
			"page":  form.Page,
			"size":  form.PageSize,
		}, nil
	}
	list, total, err := activity2.BocService{}.GetApplyRecordPageList(activity2.GetRecordPageListParam{
		UserId:      user.ID,
		ApplyStatus: form.ApplyStatus,
		Offset:      form.Offset(),
		Limit:       form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":  list,
		"total": total,
		"page":  form.Page,
		"size":  form.PageSize,
	}, nil
}
func (b BocController) FindOrCreateRecord(c *gin.Context) (gin.H, error) {
	log.Println("活动结束时间", bocEndTime, time.Now(), bocEndTime.Before(time.Now()))

	form := AddBocRecordFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	source := form.Source
	switch form.ShareUserId {
	case 297489:
		source = "mio-oa"
	case 297108:
		source = "boc-sms"
	case 297490:
		source = "boc-oa"
	case 297492:
		source = "boc-app"
	case 297549:
		source = "mio-poster"
	}

	user := apiutil.GetAuthUser(c)
	record, err := activity2.DefaultBocService.FindOrCreateApplyRecord(activity2.AddApplyRecordParam{
		UserId:      user.ID,
		ShareUserId: form.ShareUserId,
		Source:      source,
	})

	if err != nil {
		return nil, err
	}
	isOldUser, err := activity2.DefaultBocService.IsOldUserById(user.ID)
	if err != nil {
		return nil, err
	}

	type recordDetail struct {
		*activityM.BocRecord
		IsOldUser     bool   `json:"isOldUser"`
		ActivityIsEnd bool   `json:"activityIsEnd"`
		Now           string `json:"now"`
	}

	return gin.H{
		"record": recordDetail{
			BocRecord:     record,
			IsOldUser:     isOldUser,
			ActivityIsEnd: bocEndTime.Before(time.Now()),
			Now:           time.Now().String(),
		},
	}, nil
}
func (b BocController) Answer(c *gin.Context) (gin.H, error) {
	form := AnswerBocQuestionFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	err := activity2.DefaultBocService.AnswerQuestion(user.ID, int(form.Right))
	return nil, err
}
func (b BocController) FindRecordOfMini(c *gin.Context) (gin.H, error) {
	log.Println("活动结束时间", bocEndTime, time.Now(), bocEndTime.Before(time.Now()))
	user := apiutil.GetAuthUser(c)
	record, err := activity2.DefaultBocService.FindApplyRecordMini(user.ID)
	if err != nil {
		return nil, err
	}
	isOldUser, err := activity2.DefaultBocService.IsOldUserById(user.ID)
	if err != nil {
		return nil, err
	}

	type recordDetail struct {
		*activityM.BocRecord
		IsOldUser     bool   `json:"isOldUser"`
		ActivityIsEnd bool   `json:"activityIsEnd"`
		Now           string `json:"now"`
	}

	return gin.H{
		"record": recordDetail{
			BocRecord:     record,
			IsOldUser:     isOldUser,
			ActivityIsEnd: bocEndTime.Before(time.Now()),
			Now:           time.Now().String(),
		},
	}, nil
}
func (b BocController) ApplySendBonus(c *gin.Context) (gin.H, error) {
	form := ApplySendBonusForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	switch form.Type {
	case "apply":
		return nil, activity2.DefaultBocService.ApplySendApplyBonus(user.ID)
	case "bind":
		return nil, activity2.DefaultBocService.ApplySendBindWechatBonus(user.ID)
	case "boc":
		return nil, activity2.DefaultBocService.ApplySendBocBonus(user.ID)
	}
	return nil, nil
}
