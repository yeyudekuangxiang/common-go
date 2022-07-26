package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCompanyController = CompanyController{}

type CompanyController struct{}

func (CompanyController) GetCompanyInfo(ctx *gin.Context) (gin.H, error) {
	//先拿token,然后拿到公司id
	user := apiutil.GetAuthBusinessUser(ctx)
	Company := business.DefaultCompanyService.GetCompanyById(user.BCompanyId)
	//企业碳减排信息
	totalCarbonReduce := business.DefaultCarbonCreditsLogService.GetUserTotalCarbonCreditsByCid(Company.ID)
	return gin.H{
		"info":   Company,
		"carbon": totalCarbonReduce.Total,
	}, nil
}
