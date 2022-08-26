package pugc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"mio/config"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/pugc"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/wxapp"
	"os"
	"strconv"
	"time"
)

var DefaultPugcController = PugcController{}

type PugcController struct {
}

func (PugcController) AddPugc(c *gin.Context) (gin.H, error) {
	f, err := excelize.OpenFile("/Users/leo/Downloads/test1.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("sheet1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var p Pugc.PugcAddModel
	for _, row := range rows {
		p.UserId, _ = strconv.Atoi(row[0])
		p.Title = row[1]
		p.CreatedTime = time.Now()
		service.DefaultPugcService.InsertPugc(&p)
		fmt.Println()
	}

	if err != nil {
		return nil, err
	}
	return gin.H{
		"Pugc": "Pugc",
	}, nil
}

func (PugcController) ExportExcel(c *gin.Context) (gin.H, error) {

	resp, _ := service.DefaultUserService.CheckUserRisk(wxapp.UserRiskRankParam{
		AppId:    config.Config.Weapp.AppId,
		OpenId:   "oy_BA5Is69X-X2hlNn9HxRH3lOhI",
		Scene:    0,
		ClientIp: c.ClientIP(),
	})
	fmt.Println(resp)

	//f, err := excelize.OpenFile("/Users/leo/Desktop/10元话费充值名单test.xlsx")
	//
	//// Get all the rows in the Sheet1.
	//rows, err := f.GetRows("Sheet1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(rows)
	//var openidArr []string
	//for _, row := range rows {
	//	openidArr = append(openidArr, row[0])
	//}
	//cas := wxamp.BatchGetUserRiskCase(openidArr)
	//fmt.Println(cas)

	//list, _ := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
	//	Limit:   10,
	//	Offset:  0,
	//	OrderBy: "id desc",
	//})
	//
	//var ids []string
	//for _, v := range list {
	//	ids = append(ids, v.OpenId)
	//}
	//
	////openid 一次最多传十个
	//cas := wxamp.BatchGetUserRiskCase(ids)
	////保存risk
	//for _, v := range list {
	//	for _, c := range cas.List {
	//		if v.OpenId == c.Openid {
	//
	//		}
	//	}
	//}
	//os.Exit(0)
	return nil, nil
}

//周年庆双倍积分奖励明细

func (PugcController) SendPoint(c *gin.Context) (gin.H, error) {
	f, err := excelize.OpenFile("/Users/apple/Desktop/0808.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	// Get all the rows in the Sheet2.
	rows, err := f.GetRows("Sheet2")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//后续进行实际话费充值操作                                      //需要手机号码
	for i, row := range rows {
		pointVal, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			continue
		}
		if i > 1 {
			return nil, nil
		}
		pointService := service.NewPointService(context.NewMioContext())
		_, sendErr := pointService.ChangeUserPointByOffline(srv_types.ChangeUserPointDTO{
			OpenId:      row[0],
			ChangePoint: pointVal,
			Type:        entity.POINT_PLATFORM,
			BizId:       util.UUID(),
		})
		println(sendErr)
		println("给" + row[0] + "发积分" + row[1])
	}
	fmt.Println("发完了")
	return nil, nil
}

func (c PugcController) carbonInit() {
	f, err := excelize.OpenFile("/Users/leo/Desktop/10元话费充值名单test.xlsx")
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows)

	var option []qnrEntity.Option
	var subject []qnrEntity.Subject
	for _, row := range rows {
		var typeSub int8

		if row[3] == "正常" {
			typeSub = 1
		}
		if row[3] == "正常" {
			typeSub = 1
		}

		subject = append(subject, qnrEntity.Subject{
			Title:      row[1],
			Remind:     row[2],
			Type:       typeSub,
			IsHide:     1,
			QnrId:      1,
			CategoryId: 1,
		})

		option = append(option, qnrEntity.Option{
			Title:  row[1],
			Remind: row[2],
		})
	}

	//qnrService.NewSubjectService(context.NewMioContext()).Create()

	os.Exit(0)
}
