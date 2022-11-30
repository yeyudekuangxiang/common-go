package quizpdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/quizmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func SendMessage(msg quizmsg.QuizMessage) error {
	//发放优惠券
	sendMsg := quizmsg.QuizSendMessage{
		Uid:              msg.Uid,
		OpenId:           msg.OpenId,
		BizId:            msg.BizId,
		QuizTime:         msg.QuizTime,
		TodayCorrectNum:  msg.TodayCorrectNum,
		TodayAnsweredNum: msg.TodayAnsweredNum,
	}
	data, err := json.Marshal(sendMsg)
	if err != nil {
		app.Logger.Errorf("答题发天津地铁优惠券失败,发放后失败 %+v %v", msg.OpenId, err.Error())
		return err
	}
	err = app.QueueProduct.Publish(data, []string{routerkey.TjMetroSend}, rabbitmq.WithPublishOptionsExchange("quiz"))
	if err != nil {
		app.Logger.Errorf("答题发天津地铁优惠券失败,发放后失败 %+v %v", msg.OpenId, err.Error())
		return err
	}
	return nil
}
