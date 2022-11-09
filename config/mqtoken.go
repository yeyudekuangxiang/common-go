package config

type mqToken struct {
	ZhuGeSendPost      string
	SmsSendSendPost    string
	YzmSmsSendSendPost string
}

//模版id
var mqTokens = mqToken{
	ZhuGeSendPost:      "RD5KLoOe3axKtOzzkFaiM", //诸葛
	SmsSendSendPost:    "OSiM9W3dkaSsPDrd1Dllp", //sms
	YzmSmsSendSendPost: "OSiM9W3dkaSsPDrd1Dllm", //sms
}

//根据路径查找token
func FindMqToken(fullPath string) string {
	switch fullPath {
	case "/api/mp2c/mq/send_zhuge":
		return mqTokens.ZhuGeSendPost
	case "/api/mp2c/mq/send_sms":
		return mqTokens.SmsSendSendPost
	case "/api/mp2c/mq/send_yzm_sms":
		return mqTokens.YzmSmsSendSendPost
	default:
		return ""
	}
	return ""
}
