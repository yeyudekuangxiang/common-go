package email

import (
	"bytes"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"os"
)

type mail struct {
	FromUser string
	Passwd   string
}

//初始化用户名和密码

func NewEmailClient(fromUser string, passwd string) mail {
	temp := mail{
		FromUser: fromUser,
		Passwd:   passwd,
	}
	return temp
}

type SendInvoiceParam struct {
	ToUser      string
	Subject     string
	ApplyDate   string
	InvoiceDate string
	Title       string
	Price       string
	Annex       []string
}

func (m mail) SendInvoice(param SendInvoiceParam) error {
	e := email.NewEmail()
	e.From = m.FromUser
	e.To = []string{param.ToUser}
	e.Subject = param.Subject
	/*// 解析html模板
	t1, err := template.ParseFiles("email-template.html")
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = t1.Execute(body, struct {
		ApplyDate   string
		InvoiceDate string
		Title       string
		Price       string
	}{
		ApplyDate:   param.ApplyDate,
		InvoiceDate: param.InvoiceDate,
		Title:       param.Title,
		Price:       param.Price,
	})
	if err != nil {
		return err
	}*/
	body1 := `<html><body><div><p>尊敬的客户您好：</p><p>您于` + param.InvoiceDate + `申请开具电子发票，现我们将电子发票发送给您，以便作为您的报销凭证。</br>发票信息如下：</br>开票日期：` + param.InvoiceDate + `</br>购方名称：` + param.Title + `</br>价税合计：` + param.Price + `</br>详情请见附件</br></p></div></body></html>`
	e.HTML = []byte(body1)
	for _, annex := range param.Annex {
		permissions := os.FileMode(0644) // 设置新的权限模式
		err := os.Chmod(annex, permissions)
		_, err = e.AttachFile(annex)
		if err != nil {
			continue
		}
	}
	err := e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", m.FromUser, m.Passwd, "smtp.gmail.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	if err != nil {
		return err
	}
	return nil
}

func (m mail) SendInvoiceV2(param SendInvoiceParam) error {
	e := email.NewEmail()
	e.From = m.FromUser
	e.To = []string{param.ToUser}
	e.Subject = param.Subject
	// 解析html模板
	t1, err := template.ParseFiles("https://miotech-mio.oss-cn-shanghai.aliyuncs.com/static/mp2c/test/email-template.html")
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = t1.Execute(body, struct {
		ApplyDate   string
		InvoiceDate string
		Title       string
		Price       string
	}{
		ApplyDate:   param.ApplyDate,
		InvoiceDate: param.InvoiceDate,
		Title:       param.Title,
		Price:       param.Price,
	})
	if err != nil {
		return err
	}

	e.HTML = body.Bytes()
	for _, annex := range param.Annex {
		_, err = e.AttachFile(annex)
		if err != nil {
			continue
		}
	}
	err = e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", m.FromUser, m.Passwd, "smtp.gmail.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	if err != nil {
		return err
	}
	return nil
}
