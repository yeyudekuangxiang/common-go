package email

import (
	"testing"
)

func Test555(t *testing.T) {
	//mxk ceoeneoq w i y sr _ _ _ _ _ _ _ _ _ _
	a := NewEmailClient("lvmiao@miotech.com", "mxkceoeneoqwiysr").SendInvoice(SendInvoiceParam{

		//a := NewEmailClient("", "").SendInvoice(SendInvoiceParam{
		ToUser:      "18840853003@163.com",
		Subject:     "测试",
		ApplyDate:   "2023年1月1日",
		InvoiceDate: "2023年1月2日",
		Title:       "绿喵",
		Price:       "测试，测试",
		Annex:       []string{"https://resources-cn.miotech.com/static/mp2c/invoice/MA1G55M8X_2CWgDhmk1l.pdf", "https://invoice-cdn.starcharge.com/89b95b3f-b1e8-4615-9640-33f1197564ca_1690817079425.pdf"},
	})
	println(a)
}
