package thirdPlatformPdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/queue/types/message/thirdPlatform"
	"mio/internal/pkg/queue/types/routerkey"
	"mio/internal/pkg/service/platform/tianjin_metro"
	"mio/internal/pkg/util"
)

func SendRobotMessage(msg thirdPlatform.TjMetroMessage) error {
	//判断是否可以发放天津地铁优惠券
	serviceTianjin := tianjin_metro.NewTianjinMetroService(context.NewMioContext())
	userInfo, ticketErr := serviceTianjin.GetTjMetroTicketStatus(msg.ThirdCouponTypes, msg.OpenId)
	//发放优惠券
	if ticketErr != nil {
		app.Logger.Infof("天津地铁答题发电子票,发放前失败 %+v", ticketErr)
		return ticketErr
	} else {
		sendMsg := thirdPlatform.TjMetroSendMessage{
			Uid:                 userInfo.ID,
			OpenId:              userInfo.OpenId,
			Phone:               userInfo.PhoneNumber,
			BizId:               util.UUID(),
			CouponCardTypeId:    msg.ThirdCouponTypes,
			DistributionChannel: "天津地铁答题发电子票",
		}
		data, err := json.Marshal(sendMsg)
		if err != nil {
			return err
		}
		err = app.QueueProduct.Publish(data, []string{routerkey.TjMetroSend}, rabbitmq.WithPublishOptionsExchange("quizExchange"))
		if err != nil {
			app.Logger.Infof("答题发天津地铁优惠券失败,发放后失败 %+v", ticketErr)
		}
	}
	return nil
}
