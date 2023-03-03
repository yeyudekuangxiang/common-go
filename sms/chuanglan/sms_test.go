package chuanglan

import (
	"testing"
)

func TestSign(t *testing.T) {

	//2c 验证码
	/*	content := "验证码123456，30分钟有效。参与低碳任务，体验格调生活。如非本人操作请忽略。"
		var c = SmsClient{
			Account:  "YZM7795025",
			Password: "P4tDNsDCXI5380",
			BaseUrl:  "https://smssh1.253.com/msg/v1/send/json",
		}
		_, err := c.Send("18840853003", content)
		if err != nil {
			return
		}*/

	//2b验证码
	/*	contentV2 := "【企业员工碳减排管理平台】验证码：654321，10分钟有效。如非本人操作，请忽略"
		var cV2 = SmsClient{
			Account:  "YZM7795025",
			Password: "P4tDNsDCXI5380",
			BaseUrl:  "https://smssh1.253.com/msg/v1/send/json",
		}
		_, errV2 := cV2.Send("18840853003", contentV2)
		if errV2 != nil {
			return
		}*/

	//营销短信，不到参数
	var cV3 = MarketSmsClient{
		Account:  "M4232956",
		Password: "8Xx53be5pXc568",
		BaseUrl:  "https://smssh1.253.com/msg/v1/send/json",
	}

	contentV3 := "恭喜你通过了绿喵社区的乐活家认证申请，请添加活动运营人员绿大可wx： 19117399953 进入乐活家社群获取相关乐活家身份权益，期待看到你的更多创作与分享。践行可持续生活方式，绿喵与你同行~退订回T "
	_, err := cV3.Send("18840853003", contentV3)
	if err != nil {
		return
	}

	var cV4 = MarketTemplateSmsClient{
		Account:  "M4232956",
		Password: "8Xx53be5pXc568",
		BaseUrl:  "https://smssh1.253.com/msg/variable/json",
	}
	contentV4 := "很遗憾，由于【{$var}】，你并未通过绿喵社区的乐活家认证申请。践行可持续生活方式，绿喵与你同行~退订回T"
	_, err = cV4.Send("18840853003", contentV4, "6767")
	if err != nil {
		return
	}

}
