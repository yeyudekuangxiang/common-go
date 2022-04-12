package service

import (
	"encoding/json"
	"errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/duiba"
	duibaApi "mio/pkg/duiba/api"
)

var DefaultDuiBaService DuiBaService

func NewDuiBaService(client *duiba.Client) DuiBaService {
	return DuiBaService{
		client: client,
	}
}

// InitDefaultDuibaService 配置文件加载后调用此方法初始化默认兑吧服务
func InitDefaultDuibaService() {
	client := duiba.NewClient(config.Config.DuiBa.AppKey, config.Config.DuiBa.AppSecret)
	DefaultDuiBaService = NewDuiBaService(client)
}

type DuiBaService struct {
	client *duiba.Client
}

// AutoLogin 获取免登陆地址
func (srv DuiBaService) AutoLogin(param AutoLoginParam) (string, error) {
	userInfo, err := DefaultUserService.GetUserById(param.UserId)
	if err != nil {
		return "", err
	}
	b, err := DefaultPointService.FindByUserId(param.UserId)
	if err != nil {
		return "", err
	}
	return srv.client.AutoLogin(duiba.AutoLoginParam{
		Uid:      userInfo.OpenId,
		Credits:  int64(b.Balance),
		Redirect: param.Path,
		DCustom:  param.DCustom,
		Transfer: param.Transfer,
		SignKeys: param.SignKeys,
	})
}
func (srv DuiBaService) AutoLoginOpenId(param AutoLoginOpenIdParam) (string, error) {
	/*userInfo, err := DefaultUserService.GetUserById(param.UserId)
	if err != nil {
		return "", err
	}*/
	b, err := DefaultPointService.FindByUserId(param.UserId)
	if err != nil {
		return "", err
	}
	return srv.client.AutoLogin(duiba.AutoLoginParam{
		Uid:      param.OpenId,
		Credits:  int64(b.Balance),
		Redirect: param.Path,
		DCustom:  param.DCustom,
		Transfer: param.Transfer,
		SignKeys: param.SignKeys,
	})
}

// ExchangeCallback 扣积分回调
func (srv DuiBaService) ExchangeCallback(form duibaApi.ExchangeForm) (*ExchangeCallbackResult, error) {
	err := srv.client.CheckSign(form)
	if err != nil {
		return nil, err
	}
	userInfo, err := DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId: form.Uid,
	})
	if err != nil {
		return nil, err
	}
	if userInfo.ID == 0 {
		return nil, errors.New("用户信息不存在")
	}
	data, err := json.Marshal(form)
	if err != nil {
		app.Logger.Errorf("%+v %v", form, err)
		return nil, errors.New("系统异常,请联系管理员")
	}
	pointTran, err := DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       form.Uid,
		Type:         entity.POINT_DUIBA,
		Value:        int(-form.Credits),
		AdditionInfo: string(data),
	})
	if err != nil {
		app.Logger.Errorf("%+v %v", form, err)
		return nil, errors.New("系统异常,请联系管理员")
	}
	point, err := DefaultPointService.FindByUserId(userInfo.ID)
	if err != nil {
		return nil, err
	}

	return &ExchangeCallbackResult{
		BizId:   pointTran.TransactionId,
		Credits: point.Balance,
	}, nil
}

// ExchangeResultNoticeCallback 积分兑换结果回调
func (srv DuiBaService) ExchangeResultNoticeCallback(form duibaApi.ExchangeResultForm) error {
	err := srv.client.CheckSign(form)
	if err != nil {
		return err
	}
	if form.Success {
		return nil
	}

	userInfo, err := DefaultUserService.GetUserBy(repository.GetUserBy{
		OpenId: form.Uid,
	})
	if err != nil {
		return err
	}
	if userInfo.ID == 0 {
		return errors.New("用户信息不存在")
	}

	pt, err := DefaultPointTransactionService.FindBy(repository.FindPointTransactionBy{
		TransactionId: form.BizId,
	})
	if err != nil {
		return err
	}
	if pt.Id == 0 {
		return nil
	}

	data, err := json.Marshal(form)
	if err != nil {
		return err
	}

	_, err = DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       form.Uid,
		Type:         entity.POINT_DUIBA,
		Value:        -pt.Value,
		AdditionInfo: string(data),
	})
	return err
}