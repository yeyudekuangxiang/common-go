package duiba

import "strconv"

type Param interface {
	ToMap() map[string]string
}

// AutoLoginParam https://www.duiba.com.cn/tech_doc_book/server/servers/general_autologin_url_2.html#%E5%85%8D%E7%99%BBurl%E5%9F%9F%E5%90%8D
type AutoLoginParam struct {
	Uid      string `json:"uid" binding:"required"` //用户唯一性标识，对应唯一一个用户且不可变 （用not_login作为uid标记游客用户，详见
	Credits  int64  `json:"credits"`                //用户积分余额（无积分体系独立活动，积分可以传0）
	Redirect string `json:"redirect"`               //登录成功后的重定向地址（需要进行urlencode编码），可以直达积分商城内的任意页面,如果不带redirect参数，默认跳转到积分商城首页
	Alipay   string `json:"alipay"`                 //支付宝账号（如启用支付宝账号锁定则必填） 使用场景详见
	RealName string `json:"realname"`               //支付宝实名（如启用支付宝账号锁定则必填，URL中需进行utf-8编码）
	QQ       string `json:"qq"`                     //QQ号（如启用Q币账号锁定则必填）
	Phone    string `json:"phone"`                  //手机号码（如启用话费账号锁定则必填
	DCustom  string `json:"dcustom"`                //自定义参数
	Transfer string `json:"transfer"`               //自定义参数
	SignKeys string `json:"signKeys"`               //自定义参数
}

func (param AutoLoginParam) ToMap() map[string]string {
	return map[string]string{
		"uid":      param.Uid,
		"credits":  strconv.FormatInt(param.Credits, 10),
		"redirect": param.Redirect,
		"alipay":   param.Alipay,
		"realname": param.RealName,
		"qq":       param.QQ,
		"phone":    param.Phone,
		"dcustom":  param.DCustom,
		"transfer": param.Transfer,
		"signKeys": param.SignKeys,
	}
}
