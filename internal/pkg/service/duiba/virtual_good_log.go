package duiba

import (
	duibaApi "gitlab.miotech.com/miotech-application/backend/common-go/duiba/api/model"
	rduiba "mio/internal/pkg/repository/duiba"
	"mio/internal/pkg/util"
)
import eduiba "mio/internal/pkg/model/entity/duiba"

var DefaultVirtualGoodLogService = VirtualGoodLogService{repo: rduiba.DefaultVirtualGoodLogRepository}

type VirtualGoodLogService struct {
	repo rduiba.VirtualGoodLogRepository
}

func (srv VirtualGoodLogService) FindVirtualGoodLog(param FindVirtualGoodLogParam) (*eduiba.VirtualGoodLog, error) {
	log := srv.repo.FindBy(rduiba.FindVirtualGoodLogBy{
		OrderNum: param.OrderNum,
		Params:   param.Params,
	})
	return &log, nil
}
func (srv VirtualGoodLogService) CreateVirtualGoodLog(good duibaApi.VirtualGood) (*eduiba.VirtualGoodLog, error) {
	log := eduiba.VirtualGoodLog{
		AppKey:        good.AppKey,
		Timestamp:     good.Timestamp.ToInt(),
		Uid:           good.Uid,
		Sign:          good.Sign,
		OrderNum:      good.OrderNum,
		DevelopBizId:  good.DevelopBizId,
		Params:        good.Params,
		Description:   good.Description,
		Account:       good.Account,
		SupplierBizId: util.UUID(),
	}
	return &log, srv.repo.Create(&log)
}
