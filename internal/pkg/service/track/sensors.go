package track

import (
	"github.com/square/go-jose/v3/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/routerkey"
)

type SensorsService struct {
	//是否开启打点
	Open bool
}

func DefaultSensorsService() *SensorsService {
	return NewSensorsService(config.Config.App.Env == "prod")
}
func NewSensorsService(open bool) *SensorsService {
	return &SensorsService{
		Open: open,
	}
}

type TrackInfo struct {
	EventName string  `json:"eventName"`
	Guid      *string `json:"guid"`
	Attr      map[string]interface{}
}

func (srv SensorsService) Track(sync bool, eventName, guid string, attr map[string]interface{}) {
	if !srv.Open {
		return
	}
	if eventName == "" || guid == "" {
		return
	}
	//同步埋点
	if sync {
		err := app.SensorsClient.Track(guid, eventName, attr, false)
		if err != nil {
			app.Logger.Errorf("神策埋点失败 %v %v %+v %v", eventName, guid, attr, err)
		}
		return
	}
	//异步埋点
	trackInfo := TrackInfo{
		EventName: eventName,
		Guid:      &guid,
		Attr:      attr,
	}
	data, err := json.Marshal(trackInfo)
	if err != nil {
		app.Logger.Errorf("神策埋点失败 %v %v %+v %v", eventName, guid, attr, err)
		return
	}
	err = app.QueueProduct.Publish(data, []string{routerkey.SensorsSend}, rabbitmq.WithPublishOptionsExchange("lvmio"))
	if err != nil {
		app.Logger.Errorf("神策埋点失败 %v %v %+v %v", eventName, guid, attr, err)
	}
	return
}
