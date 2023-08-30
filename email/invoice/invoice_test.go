package invoice

import (
	"bytes"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
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
func TestName4(t *testing.T) {
	e := email.NewEmail()

	mailUserName := "mei.liu@miotech.com" //邮箱账号
	mailPassword := "dadynuhualdsrmfi"    //邮箱授权码
	code := "12345678"                    //发送的验证码
	Subject := "验证码发送测试"                  //发送的主题

	e.From = "mei.liu@miotech.com"
	e.To = []string{"2661440161@qq.com"}
	e.Subject = Subject
	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")

	// 解析html模板
	t1, err := template.ParseFiles("email-template.html")
	if err != nil {

	}
	body := new(bytes.Buffer)
	t1.Execute(body, struct {
		TimeDateApply string
		TimeDate      string
		Title         string
		Price         string
		Annex         string
		Annex1        string
		Annex2        string
		Annex3        string
		Annex4        string
	}{
		TimeDateApply: "2022年1月3日",
		TimeDate:      "2023年2月3日",
		Title:         "绿喵生活",
		Price:         "绿喵：5元 星星 10元",
		Annex:         "附件1 /Users/liumei/code/go/src/gitlab.miotech.com/miotech-application/backend/common-go/email/email-template.html </br>附件2 /Users/liumei/code/go/src/gitlab.miotech.com/miotech-application/backend/common-go/email/email-template.html",
	})

	e.HTML = body.Bytes()
	//e.Attach(body, "email-template.html", "text/html")
	e.AttachFile("/Users/liumei/code/go/src/gitlab.miotech.com/miotech-application/backend/common-go/email/email-template.html")
	e.AttachFile("/Users/liumei/code/go/src/gitlab.miotech.com/miotech-application/backend/common-go/email/email-template.html")

	err = e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", mailUserName, mailPassword, "smtp.gmail.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	if err != nil {
		t.Fatal(err)
	}

}
