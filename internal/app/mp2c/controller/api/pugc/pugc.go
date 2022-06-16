package pugc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"mio/config"
	"mio/internal/pkg/model/entity/pugc"
	"mio/internal/pkg/service"
	"mio/pkg/wxapp"
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
