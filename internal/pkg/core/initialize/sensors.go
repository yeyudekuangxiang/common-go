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

		// 初始化 ConcurrentLoggingConsumer

		// 初始化一个 Consumer，用于数据发送
		consumer, err := sdk.InitConcurrentLoggingConsumer(config.Config.Sensors.SaServerUrl, true)
		//...
		if err != nil {
			return
		}
		// ...
		// 使用 Consumer 来构造 SensorsAnalytics 对象
		sa := sdk.InitSensorsAnalytics(consumer, "default", false)

		// 程序结束前调用 Close() ，让 Consumer 刷新所有缓存数据到文件中
		defer sa.Close()

		/*// 初始化 Debug Consumer
		consumer, err := sdk.InitDebugConsumer(config.Config.Sensors.SaServerUrl, config.Config.Sensors.Debug, config.Config.Sensors.SaRequestTimeout)
		if err != nil {
			return
		}
		// 使用 Consumer 来构造 SensorsAnalytics 对象
		sa := sdk.InitSensorsAnalytics(consumer, "default", false)
		defer sa.Close()
		*/
		*app.SensorsClient = sa
		log.Println("初始化全局神策客户端成功...")

		zhuGeAttr := make(map[string]interface{}, 0) //诸葛打点
		zhuGeAttr["shop"] = "1"
		err = app.SensorsClient.Track("37473t374t3", "sgsdgdggsd", zhuGeAttr, false)
		if err != nil {
			println("11212")
		}
		println("44444")

		/*
			guid := "qwqw121212"
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

		/*guid := "qwqw121212"
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
