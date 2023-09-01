package email

import (
	"testing"
)

func Test555(t *testing.T) {
	a := NewEmailClient("mei.liu@miotech.com", "dadynuhualdsrmfi").SendInvoice(SendInvoiceParam{
		ToUser:      "18840853003@163.com",
		Subject:     "测试",
		ApplyDate:   "2023年1月1日",
		InvoiceDate: "2023年1月2日",
		Title:       "绿喵",
		Price:       "测试，测试",
		Annex:       []string{"http://open-cdn.starcharge.com/08d17cd5-3b6f-4912-80c0-5b23d7bc7890_1690099421421.pdf"},
	})
	println(a)
}
