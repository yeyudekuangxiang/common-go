package charge

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"strings"
	"time"
)

func TimeToDuration(startTimeStr string, endTimeStr string) string {
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return ""
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return ""
	}
	duration := endTime.Sub(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60
	seconds := int(duration.Seconds()) - hours*3600 - minutes*60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

//获取时长，分钟单位，保留2位小数

func TimeToDurationMinutes(startTimeStr string, endTimeStr string) float64 {
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return 0
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return 0
	}
	duration := endTime.Sub(startTime)
	durationMinutesTime, _ := decimal.NewFromFloat(duration.Minutes()).Round(2).Float64()
	return durationMinutesTime
}

//获取流水号订单号

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSerialNumber() string {
	operatorID := "MA1G55M8X"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 18)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	result := string(b)
	return operatorID + result
}

func GetOutInvoiceId() string {
	operatorID := "MA1G55M8X_"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	result := string(b)
	return operatorID + result
}

/**
第一步：在扫码完成之后会获取到“二维码格式1”或者“二维码格式2”，在获取二维码格式2后需进行截取终端编号，如https://qrcode.starcharge.com/#/1012081501，截取得到1012081501
第二步：将获取到的文本或者截取到文本转换成互联互通编号，大部分的星星充电桩上无互联互通编码，是使用的8位与10位的编号。所以需要接入方根据扫码后字符串的长度来判断，如果字符串长度是8位或10位，则说明是星星充电。
例如：
终端编号：10004451
转换成互联互通编号（ConnectorID）后为：11000000000000010004451000
终端编号：1012081501
互联互通编号（ConnectorID）：12000000000000010120815001
桩编号为10位的转换时是12开头，转换为枪号的时候我们取前8位，枪编号的后三位中第一个0目前都是默认的，是为了以后做扩展，枪编号的最后两位设置的是枪号。
*/

func GetConnectorID(qrCode string) string {
	// 截取终端编号
	if strings.HasPrefix(qrCode, "https://qrcode.starcharge.com/#/") {
		qrCode = strings.TrimPrefix(qrCode, "https://qrcode.starcharge.com/#/")
	}
	// 转换为互联互通编号
	var connectorID string
	if len(qrCode) == 8 {
		connectorID = "110000000000000" + qrCode + "000"
	} else if len(qrCode) == 10 {
		connectorID = "120000000000000" + qrCode[:8] + "0" + qrCode[8:]
	} else {
		return ""
	}
	return connectorID
}

func RestoreConnectorID(connectorID string) string {
	if len(connectorID) != 26 {
		return ""
	}
	connectorIdStr := connectorID[0:2]
	if connectorIdStr == "11" {
		return connectorID[15:23]
	} else if connectorIdStr == "12" {
		return connectorID[15:23] + connectorID[24:26]
	} else {
		return ""
	}
}
