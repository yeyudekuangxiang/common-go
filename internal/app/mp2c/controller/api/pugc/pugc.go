package pugc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/rwxrob/uniq"
	"mio/internal/pkg/model/entity/pugc"
	"mio/internal/pkg/service"
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

	f := excelize.NewFile()
	index := f.NewSheet("code")
	f.SetCellValue("code", "A1", "工商银行")

	for i := 0; i <= 6002; i++ {
		println(uniq.Hex(6))
		f.SetCellValue("code", "A"+strconv.Itoa(i+2), uniq.Hex(6))
	}
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("/Users/leo/Downloads/工商银行601.xlsx"); err != nil {
		fmt.Println(err)
	}
	return gin.H{
		"Pugc": "",
	}, nil
}
