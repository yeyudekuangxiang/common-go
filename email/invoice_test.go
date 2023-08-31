package email

import (
	"testing"
)

func Test555(t *testing.T) {
	a := NewEmailClient("mei.liu@miotech.com", "dadynuhualdsrmfi").SendInvoice(SendInvoiceParam{
		ToUser:      "2661440161@qq.com",
		Subject:     "测试",
		ApplyDate:   "2023年1月1日",
		InvoiceDate: "2023年1月2日",
		Title:       "绿喵",
		Price:       "测试，测试",
		Annex:       []string{"1", "12"},
	})
	println(a)
}

/**
邮箱：lvmiao@miotech.com
密码：lvmiao123
"mei.liu@miotech.com", "JN<dA6TF1"
*/
