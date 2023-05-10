package initialize

import (
	sdk "github.com/sensorsdata/sa-sdk-go"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
)

//var SensorsClient = &sdk.SensorsAnalytics{}

func InitSensors() {
	if config.Config.Sensors.SaServerUrl != "" {
		log.Println("开始初始化神策全局客户端...")
		if config.Config.App.Env != "prod" {
			// 初始化 Debug Consumer
			consumer, err := sdk.InitDebugConsumer(config.Config.Sensors.SaServerUrl, config.Config.Sensors.Debug, config.Config.Sensors.SaRequestTimeout)
			if err != nil {
				return
			}
			// 使用 Consumer 来构造 SensorsAnalytics 对象
			sa := sdk.InitSensorsAnalytics(consumer, "default", false)
			defer sa.Close()

			*app.SensorsClient = sa
		} else {
			//生产用
			consumer, err := sdk.InitBatchConsumer(config.Config.Sensors.SaServerUrl, config.Config.Sensors.BatchMax, config.Config.Sensors.SaRequestTimeout)
			if err != nil {
				return
			}
			// 使用 Consumer 来构造 SensorsAnalytics 对象
			sa := sdk.InitSensorsAnalytics(consumer, "production", false)
			defer sa.Close()
			*app.SensorsClient = sa
		}

		log.Println("初始化全局神策客户端成功...")
		/*
			guid := "test"
			properties := map[string]interface{}{
				"price": 12,
				"name":  "apple",
			}
			a := TrackInfo{
				EventName: "testlm",
				Guid:      &guid,
				Attr:      properties,
			}
			data, err := json.Marshal(a)
			err = app.QueueProduct.Publish(data, []string{"sensors.track"}, rabbitmq.WithPublishOptionsExchange("lvmio"))
			if err != nil {
				app.Logger.Errorf("答题发天津地铁优惠券失败,发放后失败 %+v %v", "1212", err.Error())
			}
		*/

		/*guid := "testd"
		properties := map[string]interface{}{
			"price": 12,
			"name":  "apple",
		}
		a := TrackInfo{
			EventName: "testlm",
			Guid:      &guid,
			Attr:      properties,
		}
		data, err := json.Marshal(a)
		err = app.QueueProduct.Publish(data, []string{"sensors.track"}, rabbitmq.WithPublishOptionsExchange("lvmio"))
		if err != nil {
			app.Logger.Errorf("答题发天津地铁优惠券失败,发放后失败 %+v %v", "1212", err.Error())
		}
		*/

	}
}
