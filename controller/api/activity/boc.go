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
