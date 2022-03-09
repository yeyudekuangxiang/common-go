package activity

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service/activity"
)

var DefaultBocController = BocController{}

type BocController struct {
}

func (b BocController) GetRecordList(c *gin.Context) (gin.H, error) {
	form := GetBocRecordListForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	if user.ID == 0 {
		return gin.H{
			"list":  make([]interface{}, 0),
			"total": 0,
			"page":  form.Page,
			"size":  form.PageSize,
		}, nil
	}
	list, total, err := activity.BocService{}.GetApplyRecordPageList(activity.GetRecordPageListParam{
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
	form := AddBocRecordFrom{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)
	record, err := activity.DefaultBocService.FindOrCreateApplyRecord(activity.AddApplyRecordParam{
		UserId:      user.ID,
		ShareUserId: form.ShareUserId,
	})

	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": record,
	}, nil
}
func (b BocController) Answer(c *gin.Context) (gin.H, error) {
	form := AnswerBocQuestionFrom{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := util.GetAuthUser(c)

	err := activity.DefaultBocService.AnswerQuestion(user.ID, int(form.Right))
	return nil, err
}
func (b BocController) FindRecordOfMini(c *gin.Context) (gin.H, error) {

	user := util.GetAuthUser(c)
	record, err := activity.DefaultBocService.FindApplyRecord(user.ID)

	if err != nil {
		return nil, err
	}
	return gin.H{
		"record": record,
	}, nil
}
