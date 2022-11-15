package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service"
)

var DefaultIndexIconController = NewIndexIconController()

func NewIndexIconController() IndexIconController {
	return IndexIconController{}
}

type IndexIconController struct {
}

func (ctl IndexIconController) Page(c *gin.Context) (gin.H, error) {
	list, _, err := service.NewIndexIconService(context.NewMioContext()).Page(repotypes.GetIndexIconPageDO{
		Offset: 0,
		Limit:  20,
	})
	if err != nil {
		return gin.H{}, err
	}
	voRow1 := make([]api_types.IndexIconApiVO, 0)
	voRow2 := make([]api_types.IndexIconApiVO, 0)
	for _, activity := range list {
		info := api_types.IndexIconApiVO{
			ID:     activity.ID,
			Title:  activity.Title,
			Type:   activity.Type,
			RowNum: activity.RowNum,
			Sort:   activity.Sort,
			Pic:    activity.Pic,
			Custom: activity.Custom,
			Abbr:   activity.Abbr,
		}
		switch activity.RowNum {
		case "1":
			voRow1 = append(voRow1, info)
		case "2":
			voRow2 = append(voRow2, info)
		}
	}
	if err != nil {
		return nil, err
	}
	return gin.H{
		"row1": voRow1,
		"row2": voRow2,
	}, nil
}
